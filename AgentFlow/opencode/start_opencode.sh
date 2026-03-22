#!/bin/bash

# OpenCode + Daytona 后台启动脚本
# 逻辑：杀掉旧进程 -> 启动 Daytona Server -> 启动 OpenCode

# --- 配置区 ---
PORT=3001
LOG_LEVEL="INFO"
LOG_FILE="opencode.log"
export DEEPSEEK_API_KEY=sk-6696ebd5e41f4fd9a1a218b57a85ad6b

echo "🛑 正在停止现有的进程..."

# 1. 杀掉 OpenCode 进程
PIDS=$(pgrep -f "opencode")
if [ -n "$PIDS" ]; then
    kill $PIDS 2>/dev/null || true
    sleep 1
    kill -9 $PIDS 2>/dev/null || true
    echo "✅ 已清理 OpenCode 旧进程。"
fi

# 2. 检查并启动 Daytona Server
echo "🚀 检查 Daytona Server..."
if ! daytona server status &> /dev/null; then
    echo "📦 正在启动 Daytona Server..."
    nohup daytona server > /dev/null 2>&1 &
    sleep 3
fi
echo "✅ Daytona Server 已就绪。"

# 3. 启动 OpenCode 服务
echo "🚀 正在后台启动 OpenCode 服务 (端口: $PORT)..."
nohup opencode serve --port $PORT --log-level $LOG_LEVEL --print-logs > "$LOG_FILE" 2>&1 &

NEW_PID=$!
sleep 1

if ps -p $NEW_PID > /dev/null; then
    echo "✅ 服务已成功在后台启动！"
    echo "📌 PID: $NEW_PID"
    echo "📂 日志文件: $(pwd)/$LOG_FILE"
    echo "📝 您可以运行 'tail -f $LOG_FILE' 实时查看日志。"
else
    echo "❌ 启动失败，请检查 $LOG_FILE 中的错误信息。"
    exit 1
fi
