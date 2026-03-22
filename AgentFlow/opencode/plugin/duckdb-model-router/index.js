import { createRequire } from 'node:module';
import fs from 'node:fs';

const require = createRequire(import.meta.url);

const DB_PATH = '/opt/duckdb/agentflow.duckdb';
let db = null;
const modelProviders = new Map();

async function ensureInitialized() {
    if (db) return;
    if (!fs.existsSync(DB_PATH)) return;

    const duckdb = require('duckdb');
    db = new duckdb.Database(DB_PATH);
    return new Promise((resolve) => {
        db.all("SELECT id, provider, base_url, api_key, model_name, name FROM models", (err, rows) => {
            if (!err && rows) {
                rows.forEach(row => {
                    if (row.base_url) {
                        try {
                            const { createOpenAICompatible } = require('@ai-sdk/openai-compatible');
                            const providerInstance = createOpenAICompatible({
                                name: row.id,
                                baseURL: row.base_url,
                                apiKey: row.api_key || 'dummy',
                            });
                            modelProviders.set(row.id, {
                                provider: providerInstance,
                                modelName: row.model_name
                            });
                        } catch (e) {}
                    }
                });
            }
            resolve();
        });
    });
}

export const DuckdbModelRouter = async () => {
    return {
        'sse.message': async (context) => {
            try {
                await ensureInitialized();
                if (!db) return context;

                const userId = context.request?.headers?.['x-user-id'];
                const modelIdOverride = context.request?.headers?.['x-model-id'];

                if (!userId && !modelIdOverride) return context;

                const targetModel = await new Promise(resolve => {
                    if (modelIdOverride) {
                        db.get("SELECT * FROM models WHERE id = ?", [modelIdOverride], (err, row) => resolve(row));
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
                            if (row) resolve(row);
                            else db.get("SELECT * FROM models WHERE is_default = TRUE LIMIT 1", (err2, def) => resolve(def));
                        },
                    );
                });

                if (targetModel && modelProviders.has(targetModel.id)) {
                    const inst = modelProviders.get(targetModel.id);
                    context.model = {
                        provider: inst.provider,
                        modelId: inst.modelName
                    };
                    if (context.response?.messages) {
                        context.response.messages.push({
                            role: 'system',
                            content: `[DEBUG] ROUTED_TO: ${targetModel.id}`
                        });
                    }
                }
            } catch (err) {
                console.error("[Router Plugin Error]", err);
            }
            return context;
        }
    };
};