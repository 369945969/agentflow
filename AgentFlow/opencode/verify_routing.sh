#!/bin/bash

# OpenCode 多模型路由动态验证脚本 (从 API 动态查询 ID)
# 运行前请确保：1. Go 后端已启动(3000) 2. opencode server 已启动(3001)

BASE_URL="http://localhost:3001"
BACKEND_API="http://localhost:3000/api/models/"
SSE_PATH="/api/sse/chat"

echo "🔍 开始功能动态验证..."

# 1. 检查 jq 是否安装 (用于解析 JSON)
if ! command -v jq &> /dev/null; then
    echo "❌ 错误: 未安装 jq，无法解析 API 返回的 JSON。"
    echo "请运行 'brew install jq' (macOS) 或 'sudo apt install jq' (Linux)"
    exit 1
fi

# 2. 动态获取模型 ID (通过 Go 后端 API)
echo "📡 正在从后端 API 读取模型列表..."
MODELS_JSON=$(curl -s "$BACKEND_API")

if [ -z "$MODELS_JSON" ] || [ "$MODELS_JSON" == "[]" ]; then
    echo "❌ 错误: 后端返回的模型列表为空。请确保 Go 后端正在运行且已配置模型。"
    exit 1
fi

# 获取默认模型 (is_default = true)
DEFAULT_MODEL_ID=$(echo "$MODELS_JSON" | jq -r '.[] | select(.is_default == true) | .id' | head -n 1)
DEFAULT_MODEL_NAME=$(echo "$MODELS_JSON" | jq -r '.[] | select(.is_default == true) | .name' | head -n 1)

# 获取一个非默认模型 (is_default = false)
OTHER_MODEL_ID=$(echo "$MODELS_JSON" | jq -r '.[] | select(.is_default == false) | .id' | head -n 1)
OTHER_MODEL_NAME=$(echo "$MODELS_JSON" | jq -r '.[] | select(.is_default == false) | .name' | head -n 1)

if [ "$DEFAULT_MODEL_ID" == "null" ] || [ -z "$DEFAULT_MODEL_ID" ]; then
    echo "⚠️  未找到标记为 '默认' 的模型，将尝试使用第一个模型。"
    DEFAULT_MODEL_ID=$(echo "$MODELS_JSON" | jq -r '.[0].id')
    DEFAULT_MODEL_NAME=$(echo "$MODELS_JSON" | jq -r '.[0].name')
fi

echo "   ✅ 默认模型: $DEFAULT_MODEL_NAME ($DEFAULT_MODEL_ID)"
[ "$OTHER_MODEL_ID" != "null" ] && [ -n "$OTHER_MODEL_ID" ] && echo "   ✅ 备选模型: $OTHER_MODEL_NAME ($OTHER_MODEL_ID)"

# 3. 检查 OpenCode 服务状态
echo -n "🚀 检查 OpenCode 服务 (Port 3001): "
HEALTH=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health" || echo "000")
if [ "$HEALTH" == "200" ]; then
    echo "✅ 正常"
else
    echo "❌ 失败 (状态码: $HEALTH)"
    echo "请运行 'bash start_opencode.sh' 启动服务。"
    exit 1
fi

# 4. 测试 1: 匿名用户自动路由到默认模型
echo "4. 测试: 匿名用户自动路由 (预期路由到: $DEFAULT_MODEL_NAME)"
# 发送消息并捕获日志中注入的调试信息
RESULT1=$(curl -s -X POST "$BASE_URL$SSE_PATH" \
  -H "X-User-ID: anonymous_$(date +%s)" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}' | grep "ROUTED_TO" | head -n 1)

echo "   结果: $RESULT1"
if [[ "$RESULT1" == *"$DEFAULT_MODEL_ID"* ]]; then
    echo "   ✅ 验证成功: 已正确路由到默认模型"
else
    echo "   ❌ 验证失败: 路由结果不匹配"
fi

# 5. 测试 2: 手动指定模型覆盖
if [ "$OTHER_MODEL_ID" != "null" ] && [ -n "$OTHER_MODEL_ID" ]; then
    echo "5. 测试: 手动指定模型 (强制切换到: $OTHER_MODEL_NAME)"
    RESULT2=$(curl -s -X POST "$BASE_URL$SSE_PATH" \
      -H "X-User-ID: tester" \
      -H "X-Model-ID: $OTHER_MODEL_ID" \
      -H "Content-Type: application/json" \
      -d '{"message": "你好"}' | grep "ROUTED_TO" | head -n 1)

    echo "   结果: $RESULT2"
    if [[ "$RESULT2" == *"$OTHER_MODEL_ID"* ]]; then
        echo "   ✅ 验证成功: 手动指定生效"
    else
        echo "   ❌ 验证失败: 路由结果不匹配"
    fi
fi

echo ""
echo "✨ 动态验证结束！"
