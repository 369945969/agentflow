import { createRequire } from "node:module";
import fs from "node:fs";

const require = createRequire(import.meta.url);

const DB_PATH = "/opt/duckdb/agentflow.duckdb";
const BACKEND_MODELS_URL =
  process.env.AGENTFLOW_MODELS_URL || "http://127.0.0.1:3000/api/models/";
const MODELS_CACHE_TTL_MS = 5000;
let db = null;
const modelProviders = new Map();
let backendModelsCache = [];
let backendModelsCacheAt = 0;

async function ensureInitialized() {
  if (db) return;
  if (!fs.existsSync(DB_PATH)) {
    console.warn(`[duckdb-model-router] missing db: ${DB_PATH}`);
    return;
  }

  const duckdb = require("duckdb");
  db = new duckdb.Database(DB_PATH);
  await new Promise((resolve) => {
    db.all(
      "SELECT id, provider, base_url, api_key, model_name, name FROM models",
      (err, rows) => {
        if (err) {
          console.error("[duckdb-model-router] load models error", err);
          resolve();
          return;
        }
        for (const row of rows || []) {
          if (!row.base_url) continue;
          try {
            const { createOpenAICompatible } = require("@ai-sdk/openai-compatible");
            const providerInstance = createOpenAICompatible({
              name: row.id,
              baseURL: row.base_url,
              apiKey: row.api_key || "dummy",
            });
            modelProviders.set(row.id, {
              provider: providerInstance,
              modelName: row.model_name,
            });
          } catch (e) {
            console.error("[duckdb-model-router] provider init error", row.id, e);
          }
        }
        console.info(`[duckdb-model-router] loaded models=${modelProviders.size}`);
        resolve();
      },
    );
  });
}

async function loadBackendModels(force = false) {
  const now = Date.now();
  if (!force && now - backendModelsCacheAt < MODELS_CACHE_TTL_MS && backendModelsCache.length > 0) {
    return backendModelsCache;
  }

  try {
    const resp = await fetch(BACKEND_MODELS_URL, {
      headers: { Accept: "application/json" },
    });
    if (!resp.ok) {
      throw new Error(`backend models http ${resp.status}`);
    }
    const rows = await resp.json();
    backendModelsCache = Array.isArray(rows) ? rows : [];
    backendModelsCacheAt = now;
    for (const row of backendModelsCache) {
      ensureProviderLoaded(row);
    }
    console.info(`[duckdb-model-router] backend models refreshed count=${backendModelsCache.length}`);
  } catch (err) {
    console.warn("[duckdb-model-router] backend models refresh failed", err);
  }
  return backendModelsCache;
}

function ensureProviderLoaded(row) {
  if (!row?.id || !row?.base_url || modelProviders.has(row.id)) return;
  try {
    const { createOpenAICompatible } = require("@ai-sdk/openai-compatible");
    const providerInstance = createOpenAICompatible({
      name: row.id,
      baseURL: row.base_url,
      apiKey: row.api_key || "dummy",
    });
    modelProviders.set(row.id, {
      provider: providerInstance,
      modelName: row.model_name,
    });
  } catch (e) {
    console.error("[duckdb-model-router] provider init error", row.id, e);
  }
}

function getHeader(headers, key) {
  if (!headers) return "";
  return (
    headers[key] ||
    headers[key.toLowerCase()] ||
    headers[key.toUpperCase()] ||
    ""
  );
}

function getSessionHeaders(input, output) {
  return (
    input?.request?.headers ||
    input?.headers ||
    output?.request?.headers ||
    output?.headers ||
    {}
  );
}

async function getTargetModel(userId, modelIdOverride) {
  const backendModels = await loadBackendModels();
  if (modelIdOverride) {
    const backendMatch = backendModels.find((row) => row?.id === modelIdOverride);
    if (backendMatch) {
      ensureProviderLoaded(backendMatch);
      return backendMatch;
    }
  }

  if (!db) return null;
  return new Promise((resolve) => {
    if (modelIdOverride) {
      db.get("SELECT * FROM models WHERE id = ?", [modelIdOverride], (err, row) => {
        if (err) console.error("[duckdb-model-router] query by id error", err);
        resolve(row || null);
      });
      return;
    }

    db.get(
      `
        SELECT m.* FROM user_model_bindings b
        JOIN models m ON b.model_id = m.id
        WHERE b.user_id = ?
      `,
      [userId],
      (err, row) => {
        if (err) {
          console.error("[duckdb-model-router] query by user error", err);
          resolve(null);
          return;
        }
        if (row) {
          resolve(row);
          return;
        }
        db.get(
          "SELECT * FROM models WHERE is_default = TRUE LIMIT 1",
          (err2, def) => {
            if (err2) console.error("[duckdb-model-router] query default error", err2);
            resolve(def || null);
          },
        );
      },
    );
  });
}

async function routeModel(input, output) {
  try {
    await ensureInitialized();
    if (!db) return;

    const headers = getSessionHeaders(input, output);
    const userId = String(getHeader(headers, "x-user-id")).trim();
    const modelIdOverride = String(getHeader(headers, "x-model-id")).trim();
    const requestedModelID = String(
      input?.model?.modelID || output?.model?.modelID || "",
    ).trim();
    const requestedProviderID = String(
      input?.model?.providerID || output?.model?.providerID || "",
    ).trim();

    console.info(
      `[duckdb-model-router] hook user=${userId || "-"} model_override=${modelIdOverride || "-"} requested=${requestedProviderID || "-"}:${requestedModelID || "-"} input_keys=${Object.keys(input || {}).join(",")} output_keys=${Object.keys(output || {}).join(",")}`,
    );

    let resolvedModelId = modelIdOverride;
    if (!resolvedModelId && requestedModelID) {
      const backendModels = await loadBackendModels();
      const requestedMatch =
        backendModels.find((row) => row?.model_name === requestedModelID) ||
        backendModels.find((row) => row?.name === requestedModelID) ||
        (backendModels.length === 1 ? backendModels[0] : null);
      if (requestedMatch?.id) {
        resolvedModelId = requestedMatch.id;
      }
    }

    if (!userId && !resolvedModelId) return;

    const targetModel = await getTargetModel(userId, resolvedModelId);
    if (!targetModel) {
      console.warn("[duckdb-model-router] no target model resolved");
      return;
    }

    const inst = modelProviders.get(targetModel.id);
    if (!inst) {
      console.warn(`[duckdb-model-router] model not preloaded id=${targetModel.id}`);
      return;
    }

    if (input) {
      input.model = {
        provider: inst.provider,
        modelId: inst.modelName,
      };
    }
    if (output) {
      output.model = {
        provider: inst.provider,
        modelId: inst.modelName,
      };
    }

    console.info(
      `[duckdb-model-router] routed model_id=${targetModel.id} model_name=${inst.modelName}`,
    );
  } catch (err) {
    console.error("[duckdb-model-router] route error", err);
  }
}

export default async function DuckdbModelRouterPlugin() {
  await ensureInitialized();
  return {
    "sse.message": async (input, output) => {
      await routeModel(input, output);
    },
    "chat.message": async (input, output) => {
      await routeModel(input, output);
    },
  };
}
