#!/bin/bash

# OpenCode + DuckDB 自动化安装脚本 (Ubuntu/Debian)
# 逻辑：程序检测安装，内容有差异才备份并替换
set -e

PROJECT_ROOT=$(pwd)
PLUGIN_DIR="$HOME/.config/opencode/plugins/duckdb-model-router"
CONFIG_DIR="$HOME/.config/opencode"
DATE_SUFFIX=$(date +%Y%m%d_%H%M%S)

echo "🚀 开始在 Ubuntu 上安装 OpenCode + DuckDB 环境..."

# 1. 更新包列表
sudo apt-get update -y

# 2. 检查并安装 Node.js (使用 NodeSource 20.x)
if command -v node &> /dev/null; then
    echo "✅ Node.js 已安装: $(node -v)"
else
    echo "📦 安装 Node.js..."
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
    sudo apt-get install -y nodejs
fi

# 3. 检查并安装 OpenCode
if command -v opencode &> /symbol/null; then
    echo "✅ OpenCode 已安装: $(opencode --version)"
else
    echo "📦 安装 OpenCode..."
    sudo npm install -g opencode
fi

# 4. 检查并安装 DuckDB CLI
if command -v duckdb &> /dev/null; then
    echo "✅ DuckDB CLI 已安装: $(duckdb --version)"
else
    echo "📦 安装 DuckDB CLI..."
    sudo apt-get install -y unzip wget
    wget https://github.com/duckdb/duckdb/releases/download/v1.0.0/duckdb_cli-linux-amd64.zip
    unzip -o duckdb_cli-linux-amd64.zip -d /usr/local/bin/
    sudo chmod +x /usr/local/bin/duckdb
    rm duckdb_cli-linux-amd64.zip
fi

# 5. 创建系统目录结构
echo "📂 检查目录结构..."
sudo mkdir -p /opt/duckdb
sudo chmod 777 /opt/duckdb
mkdir -p "$PLUGIN_DIR"
mkdir -p "$CONFIG_DIR"

# 6. 迁移现有的 DuckDB 数据库 (内容比对)
echo "🗄️ 检查 DuckDB 数据库状态..."
TARGET_DB="/opt/duckdb/agentflow.duckdb"
SOURCE_DB="$PROJECT_ROOT/../backend/apiServer/data/agentflow.duckdb"

if [ -f "$SOURCE_DB" ]; then
    if [ -f "$TARGET_DB" ]; then
        if cmp -s "$SOURCE_DB" "$TARGET_DB"; then
            echo "ℹ️  数据库文件一致，无需迁移。"
        else
            echo "💾 数据库有差异，备份旧数据库: $TARGET_DB -> $TARGET_DB.$DATE_SUFFIX"
            mv "$TARGET_DB" "$TARGET_DB.$DATE_SUFFIX"
            cp "$SOURCE_DB" "$TARGET_DB"
            echo "✅ 已更新数据库: $TARGET_DB"
        fi
    else
        cp "$SOURCE_DB" "$TARGET_DB"
        echo "✅ 已初始化数据库: $TARGET_DB"
    fi
else
    echo "⚠️  警告: 未找到源数据库文件 $SOURCE_DB"
fi

# 7. 部署插件和配置文件 (内容比对)
echo "🚚 部署插件和配置文件..."

deploy_file() {
    local src=$1
    local dest=$2
    if [ -f "$dest" ]; then
        if cmp -s "$src" "$dest"; then
            echo "ℹ️  文件一致，无需部署: $dest"
            return 0
        else
            echo "💾 文件有差异，备份旧文件: $dest -> $dest.$DATE_SUFFIX"
            mv "$dest" "$dest.$DATE_SUFFIX"
        fi
    fi
    cp "$src" "$dest"
    echo "✅ 已更新: $dest"
}

deploy_file "$PROJECT_ROOT/plugin/index.js" "$PLUGIN_DIR/index.js"
deploy_file "$PROJECT_ROOT/plugin/package.json" "$PLUGIN_DIR/package.json"

# 特殊处理 opencode.json 以替换绝对路径
echo "🚚 部署并配置 opencode.json..."
if [ -f "$CONFIG_DIR/opencode.json" ]; then
    # 临时生成一个处理过的版本进行比对
    sed "s|PLUGIN_PATH_PLACEHOLDER|$PLUGIN_DIR|g" "$PROJECT_ROOT/opencode.json" > "$PROJECT_ROOT/opencode.json.tmp"
    if cmp -s "$PROJECT_ROOT/opencode.json.tmp" "$CONFIG_DIR/opencode.json"; then
        echo "ℹ️  opencode.json 已是最新，无需部署。"
        rm "$PROJECT_ROOT/opencode.json.tmp"
    else
        echo "💾 opencode.json 有差异，备份旧文件..."
        mv "$CONFIG_DIR/opencode.json" "$CONFIG_DIR/opencode.json.$DATE_SUFFIX"
        mv "$PROJECT_ROOT/opencode.json.tmp" "$CONFIG_DIR/opencode.json"
        echo "✅ 已更新: $CONFIG_DIR/opencode.json"
    fi
else
    sed "s|PLUGIN_PATH_PLACEHOLDER|$PLUGIN_DIR|g" "$PROJECT_ROOT/opencode.json" > "$CONFIG_DIR/opencode.json"
    echo "✅ 已初始化: $CONFIG_DIR/opencode.json"
fi

# 8. 安装插件依赖
echo "📦 更新插件依赖..."
cd "$PLUGIN_DIR"
npm install

echo "✨ Ubuntu 安装与更新完成！"
echo "💡 运行 'opencode serve' 启动服务。"
