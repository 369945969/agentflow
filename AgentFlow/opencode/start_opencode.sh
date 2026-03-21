#!/bin/bash

export DEEPSEEK_API_KEY=sk-6696ebd5e41f4fd9a1a218b57a85ad6b
# OpenCode 后台启动脚本
# 逻辑：杀掉旧进程 -> 并在后台启动新进程 -> 输出日志到文件

# --- 配置区 ---
PORT=3001
LOG_LEVEL="INFO"
LOG_FILE="opencode.log"
HOST="${HOST:-}"
BACKGROUND="${BACKGROUND:-0}"
FOLLOW_LOGS="${FOLLOW_LOGS:-1}"

if [ -z "${DEEPSEEK_API_KEY:-}" ]; then
    echo "🧠 提示：未设置 DEEPSEEK_API_KEY，deepseek provider 与 opencode-mem 可能不可用。"
fi

echo "🛑 正在停止现有的 OpenCode 进程..."

PIDS="$(pgrep -f "opencode serve" 2>/dev/null || true)"
if [ -n "$PIDS" ]; then
    echo "发现正在运行的进程 PIDs: $PIDS"
    for pid in $PIDS; do
        if [ "$pid" = "$$" ] || [ "$pid" = "$PPID" ]; then
            continue
        fi
        kill "$pid" 2>/dev/null || true
    done
    sleep 2
    for pid in $PIDS; do
        if [ "$pid" = "$$" ] || [ "$pid" = "$PPID" ]; then
            continue
        fi
        kill -9 "$pid" 2>/dev/null || true
    done
    echo "✅ 已清理旧进程。"
fi

echo "🚀 正在后台启动 OpenCode 服务 (端口: $PORT)..."

# 如果 opencode 支持 --host，则允许通过 HOST=0.0.0.0 绑定到公网网卡
HOST_ARGS=""
if [ -n "$HOST" ]; then
    if opencode serve --help 2>/dev/null | grep -q -- '--host'; then
        HOST_ARGS="--host $HOST"
    fi
fi

if [ "$BACKGROUND" = "1" ]; then
    nohup opencode serve $HOST_ARGS --port $PORT --log-level $LOG_LEVEL --print-logs > "$LOG_FILE" 2>&1 &
    NEW_PID=$!

    sleep 1
    if ps -p $NEW_PID > /dev/null; then
        echo "✅ 服务已成功在后台启动！"
        echo "📌 PID: $NEW_PID"
        echo "📂 日志文件: $(pwd)/$LOG_FILE"

        for _ in $(seq 1 20); do
            if curl -fsS --max-time 1 "http://localhost:${PORT}/global/health" >/dev/null 2>&1; then
                echo "✅ 健康检查通过: http://localhost:${PORT}/global/health"
                break
            fi
            sleep 0.5
        done

        if [ "$FOLLOW_LOGS" = "1" ]; then
            tail -f "$LOG_FILE"
        fi
        exit 0
    fi

    echo "❌ 启动失败，请检查 $LOG_FILE 中的错误信息。"
    exit 1
fi

opencode serve $HOST_ARGS --port $PORT --log-level $LOG_LEVEL --print-logs 2>&1 | tee "$LOG_FILE"
