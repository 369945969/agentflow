#!/bin/bash

# OpenCode + DuckDB 自动化安装脚本 (CentOS/RHEL)
set -e

PROJECT_ROOT=$(pwd)
PLUGIN_DIR="$HOME/.config/opencode/plugins/duckdb-model-router"
CONFIG_DIR="$HOME/.config/opencode"
DATE_SUFFIX=$(date +%Y%m%d_%H%M%S)

echo "🚀 开始在 CentOS 上安装 OpenCode + DuckDB 环境..."

# 1. 检查并安装 Node.js (使用 NodeSource 20.x)
if command -v node &> /dev/null; then
    echo "✅ Node.js 已安装: $(node -v)"
else
    echo "📦 安装 Node.js..."
    curl -fsSL https://rpm.nodesource.com/setup_20.x | sudo bash -
    sudo yum install -y nodejs
fi

# 2. 检查并安装 OpenCode
if command -v opencode &> /dev/null; then
    echo "✅ OpenCode 已安装: $(opencode --version)"
else
    echo "📦 安装 OpenCode..."
    sudo npm install -g opencode
fi

# 3. 检查并安装 DuckDB CLI
if command -v duckdb &> /dev/null; then
    echo "✅ DuckDB CLI 已安装: $(duckdb --version)"
else
    echo "📦 安装 DuckDB CLI..."
    sudo yum install -y unzip wget
    wget https://github.com/duckdb/duckdb/releases/download/v1.0.0/duckdb_cli-linux-amd64.zip
    unzip -o duckdb_cli-linux-amd64.zip -d /usr/local/bin/
    sudo chmod +x /usr/local/bin/duckdb
    rm duckdb_cli-linux-amd64.zip
fi

# 4. 创建系统目录结构
echo "📂 检查目录结构..."
sudo mkdir -p /opt/duckdb
sudo chmod 777 /opt/duckdb
mkdir -p "$PLUGIN_DIR"
mkdir -p "$CONFIG_DIR"

# 5. 迁移现有的 DuckDB 数据库 (保持名称一致)
TARGET_DB="/opt/duckdb/agentflow.duckdb"
if [ -f "$TARGET_DB" ]; then
    echo "💾 备份旧数据库: $TARGET_DB -> $TARGET_DB.$DATE_SUFFIX"
    mv "$TARGET_DB" "$TARGET_DB.$DATE_SUFFIX"
fi

echo "🗄️ 迁移现有的 AgentFlow 数据库..."
SOURCE_DB="$PROJECT_ROOT/../backend/apiServer/data/agentflow.duckdb"
if [ -f "$SOURCE_DB" ]; then
    cp "$SOURCE_DB" "$TARGET_DB"
    echo "✅ 已成功将 $SOURCE_DB 复制到 $TARGET_DB"
else
    echo "⚠️  警告: 未找到源数据库文件 $SOURCE_DB"
fi

# 6. 部署插件和配置文件 (备份并覆盖)
echo "🚚 部署插件和配置文件..."

deploy_file() {
    local src=$1
    local dest=$2
    if [ -f "$dest" ]; then
        echo "💾 备份旧文件: $dest -> $dest.$DATE_SUFFIX"
        mv "$dest" "$dest.$DATE_SUFFIX"
    fi
    cp "$src" "$dest"
    echo "✅ 已更新: $dest"
}

deploy_file "$PROJECT_ROOT/plugin/index.js" "$PLUGIN_DIR/index.js"
deploy_file "$PROJECT_ROOT/plugin/package.json" "$PLUGIN_DIR/package.json"
deploy_file "$PROJECT_ROOT/opencode.json" "$CONFIG_DIR/opencode.json"

# 7. 安装插件依赖
echo "📦 更新插件依赖..."
cd "$PLUGIN_DIR"
npm install

echo "✨ CentOS 安装与更新完成！"
echo "💡 运行 'opencode server start' 启动服务。"
