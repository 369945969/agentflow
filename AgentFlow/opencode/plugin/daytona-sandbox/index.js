import { execFileSync } from "node:child_process";
import { tool } from "@opencode-ai/plugin";

const activeWorkspaces = new Map();
const daytonaCandidates = [
  process.env.DAYTONA_BIN,
  "daytona",
  "/usr/local/bin/daytona",
  "/opt/homebrew/bin/daytona",
  `${process.env.HOME || ""}/.daytona/bin/daytona`,
  `${process.env.HOME || ""}/bin/daytona`,
].filter(Boolean);

function sandboxError(error) {
  if (error instanceof Error && error.message) return error.message;
  return String(error);
}

function resolveDaytonaBin() {
  for (const candidate of daytonaCandidates) {
    try {
      execFileSync(candidate, ["--version"], { stdio: ["ignore", "pipe", "pipe"] });
      return candidate;
    } catch {
      continue;
    }
  }
  throw new Error("daytona command not found; please install Daytona and ensure it is on PATH");
}

function isOptionalDaytonaError(error) {
  const message = sandboxError(error);
  return (
    message.includes("daytona command not found") ||
    message.includes("no profiles found") ||
    message.includes("daytona login")
  );
}

function isDaytonaAvailable() {
  try {
    const daytonaBin = resolveDaytonaBin();
    listWorkspaces(daytonaBin);
    return { available: true, daytonaBin };
  } catch (error) {
    if (isOptionalDaytonaError(error)) {
      return { available: false, reason: sandboxError(error) };
    }
    throw error;
  }
}

function listWorkspaces(daytonaBin) {
  return execFileSync(daytonaBin, ["list"], {
    stdio: ["ignore", "pipe", "pipe"],
    encoding: "utf8",
  });
}

function createWorkspace(daytonaBin, workspaceId) {
  execFileSync(daytonaBin, ["create", workspaceId, "--non-interactive"], {
    stdio: ["ignore", "pipe", "pipe"],
    encoding: "utf8",
  });
}

function workspaceExists(listOutput, workspaceId) {
  return listOutput.includes(workspaceId);
}

async function ensureWorkspace(sessionId) {
  const daytona = isDaytonaAvailable();
  if (!daytona.available) return null;
  const daytonaBin = daytona.daytonaBin;

  const workspaceId = `opencode-${sessionId}`;
  const cached = activeWorkspaces.get(sessionId);
  if (cached === workspaceId) return workspaceId;

  const listOutput = listWorkspaces(daytonaBin);
  if (!workspaceExists(listOutput, workspaceId)) {
    createWorkspace(daytonaBin, workspaceId);
  }

  activeWorkspaces.set(sessionId, workspaceId);
  return workspaceId;
}

export default async function DaytonaSandboxPlugin(ctx) {
  const { directory } = ctx;
  return {
    "chat.message": async (input, output) => {
      const sessionId = input.sessionID;
      if (!sessionId) return;

      try {
        const workspaceId = await ensureWorkspace(sessionId);
        if (!workspaceId) return;
        output?.parts?.unshift({
          id: `prt-sandbox-${Date.now()}`,
          sessionID: sessionId,
          messageID: output.message?.id,
          type: "text",
          text: `[SANDBOX] 隔离环境已就绪: ${workspaceId}`,
          synthetic: true,
        });
      } catch (error) {
        if (isOptionalDaytonaError(error)) return;
        output?.parts?.unshift({
          id: `prt-sandbox-${Date.now()}`,
          sessionID: sessionId,
          messageID: output.message?.id,
          type: "text",
          text: `[SANDBOX_ERROR] ${sandboxError(error)}`,
          synthetic: true,
        });
      }
    },
    tool: {
      sandbox: tool({
        description: "Create or get a Daytona workspace bound to the current session.",
        args: {
          action: tool.schema.enum(["ensure", "status"]).optional(),
        },
        async execute(args, toolCtx) {
          const sessionId = toolCtx?.sessionID;
          if (!sessionId) {
            return { ok: false, error: "missing sessionID" };
          }

          try {
            const workspaceId = await ensureWorkspace(sessionId);
            if (!workspaceId) {
              return {
                ok: false,
                skipped: true,
                reason: "daytona unavailable or not logged in",
              };
            }
            if (args?.action === "status") {
              const listOutput = listWorkspaces(resolveDaytonaBin());
              return {
                ok: true,
                workspaceId,
                listed: workspaceExists(listOutput, workspaceId),
              };
            }
            return { ok: true, workspaceId };
          } catch (error) {
            if (isOptionalDaytonaError(error)) {
              return {
                ok: false,
                skipped: true,
                reason: sandboxError(error),
              };
            }
            return { ok: false, error: sandboxError(error) };
          }
        },
      }),
    },
    meta: {
      directory,
    },
  };
}
