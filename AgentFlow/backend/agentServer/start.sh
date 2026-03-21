#!/bin/bash

# Go WebSocket 服务启动脚本
PORT=3002
LOG_FILE="agent_server.log"
BACKGROUND="${BACKGROUND:-0}"
FOLLOW_LOGS="${FOLLOW_LOGS:-1}"

echo "🛑 停止现有的 WebSocket 服务..."
# 杀死占用 3002 端口的进程
PIDS=$(lsof -t -i:$PORT 2>/dev/null || true)
if [ -n "$PIDS" ]; then
    kill -9 $PIDS
    echo "✅ 已清理旧进程。"
fi

echo "🚀 正在编译并启动 Go WebSocket 服务 (端口: $PORT)..."
# 编译并运行
go build -o agent_server main.go

if [ "$BACKGROUND" = "1" ]; then
    nohup ./agent_server > "$LOG_FILE" 2>&1 &
    NEW_PID=$!

    sleep 1
    if ps -p $NEW_PID > /dev/null; then
        echo "✅ WebSocket 服务已在后台启动！"
        echo "📌 PID: $NEW_PID"
        echo "📂 日志文件: $(pwd)/$LOG_FILE"

        if [ "$FOLLOW_LOGS" = "1" ]; then
            tail -f "$LOG_FILE"
        fi
        exit 0
    fi
    echo "❌ 启动失败，请检查 $LOG_FILE"
    exit 1
fi

./agent_server 2>&1 | tee "$LOG_FILE"
