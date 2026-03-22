#!/bin/bash

# OpenCode + Daytona 后台启动脚本
# 逻辑：杀掉旧进程 -> 启动 Daytona Server -> 启动 OpenCode

# --- 配置区 ---
PORT=3001
LOG_LEVEL="INFO"
LOG_FILE="opencode.log"
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
RUN_DIR="$HOME"
export DEEPSEEK_API_KEY=sk-6696ebd5e41f4fd9a1a218b57a85ad6b

resolve_daytona_bin() {
    local candidate
    for candidate in "${DAYTONA_BIN:-}" "$(command -v daytona 2>/dev/null)" "/usr/local/bin/daytona" "/opt/homebrew/bin/daytona" "$HOME/.daytona/bin/daytona" "$HOME/bin/daytona"; do
        [ -n "$candidate" ] || continue
        [ -x "$candidate" ] || continue
        echo "$candidate"
        return 0
    done
    return 1
}

daytona_ready() {
    local daytona_bin=$1
    "$daytona_bin" list > /dev/null 2>&1
}

echo "🛑 正在停止现有的进程..."

# 1. 杀掉旧的 OpenCode serve 进程，避免误杀当前脚本
PIDS=$(pgrep -f "opencode serve" || true)
if [ -n "$PIDS" ]; then
    kill $PIDS 2>/dev/null || true
    sleep 1
    kill -9 $PIDS 2>/dev/null || true
    echo "✅ 已清理 OpenCode 旧进程。"
fi

# 2. 检查并启动 Daytona Server
echo "🚀 检查 Daytona Server..."
if DAYTONA_BIN="$(resolve_daytona_bin)"; then
    export DAYTONA_BIN
    export PATH="$(dirname "$DAYTONA_BIN"):$PATH"
    if daytona_ready "$DAYTONA_BIN"; then
        if ! "$DAYTONA_BIN" server status &> /dev/null; then
            echo "📦 正在启动 Daytona Server..."
            nohup "$DAYTONA_BIN" server > /dev/null 2>&1 &
            sleep 3
        fi
        echo "✅ Daytona Server 已就绪。"
    else
        echo "⚠️  Daytona 已安装，但当前未登录或没有 profile，沙箱功能将被跳过。"
    fi
else
    echo "⚠️  未检测到 daytona，沙箱功能将被跳过。"
fi

# 3. 启动 OpenCode 服务
echo "🚀 正在后台启动 OpenCode 服务 (端口: $PORT)..."
cd "$RUN_DIR"
nohup opencode serve --port $PORT --log-level $LOG_LEVEL --print-logs > "$PROJECT_ROOT/$LOG_FILE" 2>&1 &

NEW_PID=$!
HEALTH_URL="http://127.0.0.1:$PORT/global/health"
for _ in $(seq 1 20); do
    if curl -fsS "$HEALTH_URL" > /dev/null 2>&1; then
        break
    fi
    sleep 1
done

if ps -p $NEW_PID > /dev/null && curl -fsS "$HEALTH_URL" > /dev/null 2>&1; then
    echo "✅ 服务已成功在后台启动！"
    echo "📌 PID: $NEW_PID"
    echo "📂 日志文件: $PROJECT_ROOT/$LOG_FILE"
    echo "📝 您可以运行 'tail -f $LOG_FILE' 实时查看日志。"
else
    echo "❌ 启动失败，请检查 $LOG_FILE 中的错误信息。"
    exit 1
fi
