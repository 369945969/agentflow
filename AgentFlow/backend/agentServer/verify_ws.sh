#!/bin/bash

# WebSocket 服务自动化测试脚本
set -e

# 获取脚本所在目录
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

echo "🧪 开始验证 WebSocket 服务..."

echo "🏃 运行 Go 测试..."
go test -v ./...

echo "✨ 验证完成！"
