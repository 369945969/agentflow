#!/bin/bash

# WebSocket 服务自动化测试脚本
set -e

# 获取脚本所在目录
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

echo "🧪 开始验证 WebSocket 服务..."

# 1. 启动服务 (如果没启动)
if ! lsof -i:3002 >/dev/null; then
    echo "🚀 服务未运行，正在通过 start.sh 启动..."
    bash start.sh
    sleep 2
fi

# 2. 运行 Go Test
echo "🏃 运行 Go 集成测试..."
go test -v ws_test.go main.go

echo "✨ 验证完成！"
