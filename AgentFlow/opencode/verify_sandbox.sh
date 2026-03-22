#!/bin/bash

# Daytona 沙箱路由功能验证脚本
# 运行前请确保 opencode server 已启动

BASE_URL="http://localhost:3001"
TEST_SESSION=""

echo "🔍 开始 Daytona 沙箱功能验证..."

# 1. 检查服务状态
echo -n "🚀 检查 OpenCode 服务 (Port 3001): "
HEALTH=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/global/health" || echo "000")
if [ "$HEALTH" == "200" ]; then
    echo "✅ 正常"
else
    echo "❌ 失败 (状态码: $HEALTH)"
    exit 1
fi

# 2/3. 使用 Node.js 创建 Session、选择 provider/model、发送消息、轮询验证 [SANDBOX]
echo "📦 创建 Session 并触发沙箱..."
NODE_OUT="$(BASE_URL="$BASE_URL" node - <<'JS'
const base = process.env.BASE_URL || "http://localhost:3001";
const deadlineMs = 20000;
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

async function jsonOrNull(res) {
  const text = await res.text();
  try { return JSON.parse(text); } catch { return null; }
}

async function main() {
  const createRes = await fetch(`${base}/session`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: "{}",
  });
  const created = await jsonOrNull(createRes);
  const sessionId = created && typeof created === "object" ? created.id : "";
  if (!sessionId) {
    process.stderr.write(`[verify_sandbox] create session failed: ${JSON.stringify(created)}\n`);
    process.exit(2);
  }

  const provRes = await fetch(`${base}/provider`);
  const prov = await jsonOrNull(provRes) || {};
  const connected = Array.isArray(prov.connected) ? prov.connected : [];
  const defaults = (prov.default && typeof prov.default === "object") ? prov.default : {};
  const providerId = connected.includes("deepseek") ? "deepseek" : (connected[0] || "deepseek");
  const modelId = defaults[providerId] || "deepseek-chat";

  await fetch(`${base}/session/${sessionId}/message`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      model: { providerID: providerId, modelID: modelId },
      variant: "verify_sandbox",
      parts: [{ type: "text", text: "你好，请确认你的沙箱环境，并输出 [SANDBOX] 标记。" }],
    }),
  });

  const endAt = Date.now() + deadlineMs;
  let sandboxLine = "";
  while (Date.now() < endAt) {
    try {
      const msgRes = await fetch(`${base}/session/${sessionId}/message`);
      const payload = await jsonOrNull(msgRes);
      const msgs = Array.isArray(payload) ? payload : (payload && Array.isArray(payload.data) ? payload.data : []);
      for (const m of msgs) {
        const parts = Array.isArray(m.parts) ? m.parts : [];
        for (const p of parts) {
          if (p && p.type === "text" && typeof p.text === "string" && p.text.includes("[SANDBOX]")) {
            sandboxLine = p.text;
            break;
          }
        }
        if (sandboxLine) break;
      }
      if (sandboxLine) break;
    } catch {
    }
    await sleep(500);
  }

  const b64 = Buffer.from(sandboxLine, "utf8").toString("base64");
  process.stdout.write(`${sessionId} ${providerId} ${modelId} ${b64}\n`);
}

main().catch((err) => {
  process.stderr.write(`[verify_sandbox] fatal: ${String(err)}\n`);
  process.exit(2);
});
JS
)"
NODE_STATUS=$?
if [ $NODE_STATUS -ne 0 ]; then
  echo "   ❌ Node 验证脚本执行失败"
  exit 1
fi

TEST_SESSION="$(printf '%s' "$NODE_OUT" | awk '{print $1}')"
PROVIDER_ID="$(printf '%s' "$NODE_OUT" | awk '{print $2}')"
MODEL_ID="$(printf '%s' "$NODE_OUT" | awk '{print $3}')"
SANDBOX_B64="$(printf '%s' "$NODE_OUT" | awk '{print $4}')"

if [ -z "$TEST_SESSION" ]; then
  echo "   ❌ 创建 session 失败"
  exit 1
fi
echo "   ✅ session_id: $TEST_SESSION"
echo "🧪 发送消息到 Session（触发沙箱创建）"
echo "   - provider/model: $PROVIDER_ID/$MODEL_ID"

echo "⏳ 等待响应并验证 [SANDBOX] 标记..."
SANDBOX_LINE="$(echo "$SANDBOX_B64" | base64 -D 2>/dev/null || true)"

if [[ "$SANDBOX_LINE" == *"[SANDBOX]"* ]]; then
  echo "   ✅ 找到标记: $SANDBOX_LINE"
else
  echo "   ❌ 未找到 [SANDBOX] 标记（可能插件未生效或没有注入输出）"
  exit 1
fi

# 4. 检查 Daytona 命令行
echo "📊 检查 Daytona 内部状态..."
WORKSPACE_ID="opencode-$TEST_SESSION"
WORKSPACES=$(daytona list 2>/dev/null || true)
if [[ "$WORKSPACES" == *"$WORKSPACE_ID"* ]]; then
  echo "   ✅ 验证成功: Daytona 工作区 '$WORKSPACE_ID' 已创建。"
else
  echo "   ❌ 验证失败: 未找到工作区 '$WORKSPACE_ID'。"
  exit 1
fi

echo ""
echo "✨ Daytona 沙箱集成验证结束。"
