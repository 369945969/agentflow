#!/bin/bash

export DEEPSEEK_API_KEY=sk-6696ebd5e41f4fd9a1a218b57a85ad6b
# OpenCode 后台启动脚本
# 逻辑：杀掉旧进程 -> 并在后台启动新进程 -> 输出日志到文件

# --- 配置区 ---
PORT=3001
LOG_LEVEL="INFO"
LOG_FILE="opencode.log"

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

# 使用 nohup 后台启动
# 2>&1 将错误输出也合并到日志文件中
nohup opencode serve --port $PORT --log-level $LOG_LEVEL --print-logs > "$LOG_FILE" 2>&1 &

# 获取新进程 PID
NEW_PID=$!

# 等待一秒检查进程是否还在
sleep 1
if ps -p $NEW_PID > /dev/null; then
    echo "✅ 服务已成功在后台启动！"
    echo "📌 PID: $NEW_PID"
    echo "📂 日志文件: $(pwd)/$LOG_FILE"
    echo "📝 您可以运行 'tail -f $LOG_FILE' 实时查看日志。"

    for _ in $(seq 1 20); do
        if curl -fsS --max-time 1 "http://localhost:${PORT}/global/health" >/dev/null 2>&1; then
            echo "✅ 健康检查通过: http://localhost:${PORT}/global/health"
            exit 0
        fi
        sleep 0.5
    done
    echo "⚠️  服务进程已启动，但健康检查未通过，请检查日志: $LOG_FILE"
else
    echo "❌ 启动失败，请检查 $LOG_FILE 中的错误信息。"
    exit 1
fi
