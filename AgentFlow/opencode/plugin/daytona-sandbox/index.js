import { execSync } from "node:child_process";
import { tool } from "@opencode-ai/plugin";

const activeWorkspaces = new Set();

async function ensureWorkspace(sessionId) {
  const workspaceId = `opencode-${sessionId}`;
  if (activeWorkspaces.has(workspaceId)) return workspaceId;

  const list = execSync("daytona list", { stdio: ["ignore", "pipe", "pipe"] }).toString();
  if (!list.includes(workspaceId)) {
    execSync(`daytona create ${workspaceId} --non-interactive`, {
      stdio: ["ignore", "pipe", "pipe"],
    });
  }
  activeWorkspaces.add(workspaceId);
  return workspaceId;
}

export default async function DaytonaSandboxPlugin(ctx) {
  const { directory } = ctx;
  return {
    "chat.message": async (input, output) => {
      const sessionId = input.sessionID;
      if (!sessionId) return;

      let workspaceId = null;
      try {
        workspaceId = await ensureWorkspace(sessionId);
      } catch (error) {
        const msg = `[SANDBOX_ERROR] ${String(error)}`;
        if (output?.parts) {
          output.parts.unshift({
            id: `prt-sandbox-${Date.now()}`,
            sessionID: sessionId,
            messageID: output.message?.id,
            type: "text",
            text: msg,
            synthetic: true,
          });
        }
        return;
      }

      if (!workspaceId) return;

      const text = `[SANDBOX] 隔离环境已就绪: ${workspaceId}`;
      if (output?.parts) {
        output.parts.unshift({
          id: `prt-sandbox-${Date.now()}`,
          sessionID: sessionId,
          messageID: output.message?.id,
          type: "text",
          text,
          synthetic: true,
        });
      }
    },
    tool: {
      sandbox: tool({
        description: "Create or get a Daytona sandbox workspace bound to current session.",
        args: {
          action: tool.schema.enum(["ensure", "status"]).optional(),
        },
        async execute(args, toolCtx) {
          const sessionId = toolCtx?.sessionID;
          if (!sessionId) {
            return { ok: false, error: "missing sessionID" };
          }
          const workspaceId = await ensureWorkspace(sessionId);
          if (args?.action === "status") {
            const list = execSync("daytona list", { stdio: ["ignore", "pipe", "pipe"] }).toString();
            return { ok: true, workspaceId, listed: list.includes(workspaceId) };
          }
          return { ok: true, workspaceId };
        },
      }),
    },
    meta: {
      directory,
    },
  };
}
