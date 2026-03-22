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

# 2. 创建 Session
echo "📦 创建 Session ..."
SESSION_ID=$(curl -sS -X POST "$BASE_URL/session" -H "Content-Type: application/json" -d '{}' | python3 -c 'import sys, json; print(json.load(sys.stdin)["id"])')
TEST_SESSION="$SESSION_ID"
echo "   ✅ session_id: $TEST_SESSION"

# 3. 选择 provider/model（尽量用 deepseek 的默认模型）
read -r PROVIDER_ID MODEL_ID < <(curl -sS "$BASE_URL/provider" | python3 - <<'PY'
import json,sys
p=json.load(sys.stdin)
default=p.get("default") or {}
connected=p.get("connected") or []
provider_id="deepseek" if "deepseek" in connected else (connected[0] if connected else "deepseek")
model_id=default.get(provider_id) or "deepseek-chat"
print(provider_id, model_id)
PY)

echo "🧪 发送消息到 Session（触发沙箱创建）"
echo "   - provider/model: $PROVIDER_ID/$MODEL_ID"
curl -sS -X POST "$BASE_URL/session/$TEST_SESSION/message" \
  -H "Content-Type: application/json" \
  -d "$(python3 - <<PY
import json
print(json.dumps({
  "model": {"providerID": "$PROVIDER_ID", "modelID": "$MODEL_ID"},
  "variant": "verify_sandbox",
  "parts": [{"type": "text", "text": "你好，请确认你的沙箱环境，并输出 [SANDBOX] 标记。"}]
}))
PY)" >/dev/null

echo "⏳ 等待响应并验证 [SANDBOX] 标记..."
SANDBOX_LINE=$(python3 - <<PY
import json, time, urllib.request
base="$BASE_URL"
sid="$TEST_SESSION"
deadline=time.time()+20
while time.time()<deadline:
    try:
        data=urllib.request.urlopen(f"{base}/session/{sid}/message", timeout=3).read().decode("utf-8")
        msgs=json.loads(data) if data.strip().startswith("[") else []
        for m in msgs:
            parts=m.get("parts") or []
            for p in parts:
                if p.get("type")=="text" and "[SANDBOX]" in (p.get("text") or ""):
                    print(p.get("text"))
                    raise SystemExit(0)
    except SystemExit:
        raise
    except Exception:
        pass
    time.sleep(0.5)
raise SystemExit(1)
PY || true)

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
