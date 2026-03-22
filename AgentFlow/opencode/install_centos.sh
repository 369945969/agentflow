#!/bin/bash

# OpenCode + DuckDB + Daytona + Mem 自动化安装脚本 (CentOS/RHEL)
set -e

PROJECT_ROOT=$(pwd)
PLUGIN_DIR="$HOME/.config/opencode/plugins/duckdb-model-router"
DAYTONA_PLUGIN_DIR="$HOME/.config/opencode/plugins/daytona-sandbox"
CONFIG_DIR="$HOME/.config/opencode"
DATE_SUFFIX=$(date +%Y%m%d_%H%M%S)

echo "🚀 开始在 CentOS 上安装 OpenCode 环境..."

# 1. 检查并安装 Node.js
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

# 3. 检查并安装 opencode-mem 插件
if npm list -g opencode-mem &> /dev/null; then
    echo "✅ opencode-mem 已安装"
else
    echo "📦 安装 opencode-mem..."
    sudo npm install -g opencode-mem
fi

# 4. 检查并安装 DuckDB CLI
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

# 5. 检查并安装 Daytona
if command -v daytona &> /dev/null; then
    echo "✅ Daytona 已安装: $(daytona --version)"
else
    echo "📦 安装 Daytona..."
    curl -sfL https://download.daytona.io/daytona/install.sh | sudo bash
fi

# 6. 创建系统目录结构
echo "📂 检查目录结构..."
sudo mkdir -p /opt/duckdb
sudo chmod 777 /opt/duckdb
mkdir -p "$PLUGIN_DIR"
mkdir -p "$DAYTONA_PLUGIN_DIR"
mkdir -p "$CONFIG_DIR"

# 7. 迁移现有的 DuckDB 数据库
echo "🗄️ 检查 DuckDB 数据库状态..."
TARGET_DB="/opt/duckdb/agentflow.duckdb"
SOURCE_DB="$PROJECT_ROOT/../backend/apiServer/data/agentflow.duckdb"

check_table_exists() {
    if [ ! -f "$1" ]; then return 1; fi
    duckdb "$1" "SELECT count(*) FROM models;" >/dev/null 2>&1
    return $?
}

if [ -f "$SOURCE_DB" ]; then
    FORCE_REPLACE=false
    if [ -f "$TARGET_DB" ]; then
        if ! check_table_exists "$TARGET_DB"; then
            echo "⚠️  目标数据库缺少 models 表，将强制更新。"
            FORCE_REPLACE=true
        elif ! cmp -s "$SOURCE_DB" "$TARGET_DB"; then
            echo "💾 数据库有差异，备份旧数据库: $TARGET_DB -> $TARGET_DB.$DATE_SUFFIX"
            mv "$TARGET_DB" "$TARGET_DB.$DATE_SUFFIX"
            FORCE_REPLACE=true
        else
            echo "ℹ️  数据库文件一致且结构完整，无需迁移。"
        fi
    else
        FORCE_REPLACE=true
    fi

    if [ "$FORCE_REPLACE" = true ]; then
        cp "$SOURCE_DB" "$TARGET_DB"
        echo "✅ 已同步数据库: $TARGET_DB"
    fi
else
    echo "⚠️  警告: 未找到源数据库文件 $SOURCE_DB"
fi

# 8. 部署插件和配置文件
echo "🚚 部署插件和配置文件..."

deploy_file() {
    local src=$1
    local dest=$2
...
    fi
    cp "$src" "$dest"
    echo "✅ 已更新: $dest"
}

if grep -q "PLUGIN_PATH_PLACEHOLDER" "$PROJECT_ROOT/opencode.json"; then
    deploy_file "$PROJECT_ROOT/plugin/index.js" "$PLUGIN_DIR/index.js"
    deploy_file "$PROJECT_ROOT/plugin/package.json" "$PLUGIN_DIR/package.json"
fi

deploy_file "$PROJECT_ROOT/plugin/daytona-sandbox/index.js" "$DAYTONA_PLUGIN_DIR/index.js"
deploy_file "$PROJECT_ROOT/plugin/daytona-sandbox/package.json" "$DAYTONA_PLUGIN_DIR/package.json"
deploy_file "$PROJECT_ROOT/opencode-mem.jsonc" "$CONFIG_DIR/opencode-mem.jsonc"

echo "🚚 部署并配置 opencode.json..."
sed "s|PLUGIN_PATH_PLACEHOLDER|$PLUGIN_DIR|g; s|DAYTONA_PLUGIN_PATH_PLACEHOLDER|$DAYTONA_PLUGIN_DIR|g" "$PROJECT_ROOT/opencode.json" > "$PROJECT_ROOT/opencode.json.tmp"

if [ -f "$CONFIG_DIR/opencode.json" ]; then
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
    mv "$PROJECT_ROOT/opencode.json.tmp" "$CONFIG_DIR/opencode.json"
    echo "✅ 已初始化: $CONFIG_DIR/opencode.json"
fi

# 9. 安装插件依赖
echo "📦 更新插件依赖..."
[ -d "$PLUGIN_DIR" ] && cd "$PLUGIN_DIR" && npm install
[ -d "$DAYTONA_PLUGIN_DIR" ] && cd "$DAYTONA_PLUGIN_DIR" && npm install

echo "✨ CentOS 安装与更新完成！"
echo "💡 运行 'bash start_opencode.sh' 启动服务。"
