#!/bin/bash

# Daytona 沙箱路由功能验证脚本
# 运行前请确保 opencode server 已启动

BASE_URL="http://localhost:3001"
SSE_PATH="/api/sse/chat"
TEST_SESSION="session_$(date +%s)"
export DEEPSEEK_API_KEY=sk-6696ebd5e41f4fd9a1a218b57a85ad6b

echo "🔍 开始 Daytona 沙箱功能验证..."

# 1. 检查服务状态
echo -n "🚀 检查 OpenCode 服务 (Port 3001): "
HEALTH=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health" || echo "000")
if [ "$HEALTH" == "200" ]; then
    echo "✅ 正常"
else
    echo "❌ 失败 (状态码: $HEALTH)"
    exit 1
fi

# 2. 模拟请求并触发沙箱创建
echo "📦 模拟会话 $TEST_SESSION 的第一个请求 (触发沙箱创建)..."
echo "⏳ 这可能需要一点时间，因为 Daytona 正在初始化工作区..."

RESULT=$(curl -s -X POST "$BASE_URL$SSE_PATH" \
  -H "X-Session-ID: $TEST_SESSION" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好，请确认你的沙箱环境"}' | grep "SANDBOX" | head -n 1)

echo "   结果: $RESULT"

if [[ "$RESULT" == *"$TEST_SESSION"* ]]; then
    echo "   ✅ 验证成功: 沙箱已为会话 $TEST_SESSION 正确分配。"
    
    # 3. 检查 Daytona 命令行
    echo "📊 检查 Daytona 内部状态..."
    WORKSPACES=$(daytona list)
    if [[ "$WORKSPACES" == *"$TEST_SESSION"* ]]; then
        echo "   ✅ 验证成功: Daytona 工作区 'opencode-$TEST_SESSION' 已在系统中创建。"
    else
        echo "   ❌ 验证失败: 命令行未找到对应工作区。"
    fi
else
    echo "   ❌ 验证失败: 未能从响应中捕获沙箱就绪信息。"
fi

echo ""
echo "✨ Daytona 沙箱集成验证结束。"
