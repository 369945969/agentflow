#!/bin/bash

# OpenCode 多模型路由动态验证脚本 (从数据库动态查询 ID)
# 运行前请确保 opencode server 已启动

BASE_URL="http://localhost:3001"
SSE_PATH="/api/sse/chat"
DB_FILE="/opt/duckdb/agentflow.duckdb"

echo "🔍 开始功能动态验证..."

# 1. 检查数据库文件
if [ ! -f "$DB_FILE" ]; then
    echo "❌ 错误: 未找到数据库文件 $DB_FILE"
    exit 1
fi

# 2. 检查 DuckDB CLI
if ! command -v duckdb &> /dev/null; then
    echo "❌ 错误: 未安装 duckdb CLI，无法进行动态验证。"
    exit 1
fi

# 3. 动态获取模型 ID
echo "📡 正在从数据库读取模型列表..."
# 获取默认模型 (is_default = TRUE)
DEFAULT_MODEL_ID=$(duckdb "$DB_FILE" -c "SELECT id FROM models WHERE is_default = TRUE LIMIT 1;" | sed -n '3p' | xargs)
DEFAULT_MODEL_NAME=$(duckdb "$DB_FILE" -c "SELECT name FROM models WHERE is_default = TRUE LIMIT 1;" | sed -n '3p' | xargs)

# 获取一个非默认模型 (is_default = FALSE)
OTHER_MODEL_ID=$(duckdb "$DB_FILE" -c "SELECT id FROM models WHERE is_default = FALSE LIMIT 1;" | sed -n '3p' | xargs)
OTHER_MODEL_NAME=$(duckdb "$DB_FILE" -c "SELECT name FROM models WHERE is_default = FALSE LIMIT 1;" | sed -n '3p' | xargs)

if [ -z "$DEFAULT_MODEL_ID" ]; then
    echo "❌ 错误: 数据库中未找到默认模型。请先在 AgentFlow 界面设置默认模型。"
    exit 1
fi

echo "   ✅ 默认模型: $DEFAULT_MODEL_NAME ($DEFAULT_MODEL_ID)"
[ -n "$OTHER_MODEL_ID" ] && echo "   ✅ 备选模型: $OTHER_MODEL_NAME ($OTHER_MODEL_ID)" || echo "   ⚠️  提示: 未找到备选模型，将仅测试默认路由。"

# 4. 检查服务状态
echo -n "🚀 检查 OpenCode 服务 (Port 3001): "
HEALTH=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health" || echo "000")
if [ "$HEALTH" == "200" ]; then
    echo "✅ 正常"
else
    echo "❌ 失败 (状态码: $HEALTH)"
    echo "请运行 'opencode server start' (端口已改为 3001) 后再测试。"
    exit 1
fi

# 5. 测试 1: 匿名用户自动路由到默认模型
echo "5. 测试: 匿名用户自动路由 (预期路由到: $DEFAULT_MODEL_NAME)"
RESULT1=$(curl -s -X POST "$BASE_URL$SSE_PATH" \
  -H "X-User-ID: anonymous_$(date +%s)" \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}' | grep "ROUTED_TO" | head -n 1)

echo "   结果: $RESULT1"
if [[ "$RESULT1" == *"$DEFAULT_MODEL_ID"* ]]; then
    echo "   ✅ 验证成功: 已正确路由到默认模型"
else
    echo "   ❌ 验证失败"
fi

# 6. 测试 2: 手动指定模型覆盖 (如果存在备选模型)
if [ -n "$OTHER_MODEL_ID" ]; then
    echo "6. 测试: 手动指定模型 (强制切换到: $OTHER_MODEL_NAME)"
    RESULT2=$(curl -s -X POST "$BASE_URL$SSE_PATH" \
      -H "X-User-ID: tester" \
      -H "X-Model-ID: $OTHER_MODEL_ID" \
      -H "Content-Type: application/json" \
      -d '{"message": "你好"}' | grep "ROUTED_TO" | head -n 1)

    echo "   结果: $RESULT2"
    if [[ "$RESULT2" == *"$OTHER_MODEL_ID"* ]]; then
        echo "   ✅ 验证成功: 手动指定生效"
    else
        echo "   ❌ 验证失败"
    fi
fi

echo ""
echo "✨ 动态验证结束！"
