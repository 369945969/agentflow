#!/bin/bash

# OpenCode + DuckDB + Daytona + Mem 自动化安装脚本 (macOS)
# 逻辑：程序检测安装，内容有差异才备份并替换
set -e

PROJECT_ROOT=$(pwd)
PLUGIN_DIR="$HOME/.config/opencode/plugins/duckdb-model-router"
DAYTONA_PLUGIN_DIR="$HOME/.config/opencode/plugins/daytona-sandbox"
CONFIG_DIR="$HOME/.config/opencode"
DATE_SUFFIX=$(date +%Y%m%d_%H%M%S)

echo "🚀 开始在 macOS 上安装 OpenCode 环境..."

# 1. 检查并安装 Homebrew
if ! command -v brew &> /dev/null; then
    echo "❌ 未找到 Homebrew，请先安装 Homebrew: https://brew.sh/"
    exit 1
fi

# 2. 检查并安装 Node.js
if command -v node &> /dev/null; then
    echo "✅ Node.js 已安装: $(node -v)"
else
    echo "📦 通过 Homebrew 安装 Node.js..."
    brew install node
fi

# 3. 检查并安装 OpenCode
if command -v opencode &> /dev/null; then
    echo "✅ OpenCode 已安装: $(opencode --version)"
else
    echo "📦 安装 OpenCode..."
    sudo npm install -g opencode
fi

# 4. 检查并安装 opencode-mem 插件
if npm list -g opencode-mem &> /dev/null; then
    echo "✅ opencode-mem 已安装"
else
    echo "📦 安装 opencode-mem..."
    sudo npm install -g opencode-mem
fi

# 5. 检查并安装 DuckDB CLI
if command -v duckdb &> /dev/null; then
    echo "✅ DuckDB CLI 已安装: $(duckdb --version)"
else
    echo "📦 通过 Homebrew 安装 DuckDB..."
    brew install duckdb
fi

# 6. 检查并安装 Daytona
if command -v daytona &> /dev/null; then
    echo "✅ Daytona 已安装: $(daytona --version)"
else
    echo "📦 安装 Daytona..."
    curl -sfL https://download.daytona.io/daytona/install.sh | sudo bash
fi

# 7. 创建系统目录结构
echo "📂 检查目录结构..."
sudo mkdir -p /opt/duckdb
sudo chmod 777 /opt/duckdb
mkdir -p "$PLUGIN_DIR"
mkdir -p "$DAYTONA_PLUGIN_DIR"
mkdir -p "$CONFIG_DIR"

# 8. 迁移现有的 DuckDB 数据库
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

# 9. 部署插件和配置文件
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

# 部署路由插件 (如果被使用)
# 检查 opencode.json 模板是否包含该插件
if grep -q "PLUGIN_PATH_PLACEHOLDER" "$PROJECT_ROOT/opencode.json"; then
    deploy_file "$PROJECT_ROOT/plugin/index.js" "$PLUGIN_DIR/index.js"
    deploy_file "$PROJECT_ROOT/plugin/package.json" "$PLUGIN_DIR/package.json"
else
    echo "🗑️  检测到 duckdb-model-router 未在配置中使用，跳过部署。"
    # 如果您想删除物理文件，可以解除下面注释
    # rm -rf "$PLUGIN_DIR"
fi

# 部署 Daytona 插件
deploy_file "$PROJECT_ROOT/plugin/daytona-sandbox/index.js" "$DAYTONA_PLUGIN_DIR/index.js"
deploy_file "$PROJECT_ROOT/plugin/daytona-sandbox/package.json" "$DAYTONA_PLUGIN_DIR/package.json"

# 部署 opencode-mem 配置
deploy_file "$PROJECT_ROOT/opencode-mem.jsonc" "$CONFIG_DIR/opencode-mem.jsonc"

# 特殊处理 opencode.json 以替换绝对路径
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

# 10. 安装插件依赖
echo "📦 更新插件依赖..."
[ -d "$PLUGIN_DIR" ] && cd "$PLUGIN_DIR" && npm install
[ -d "$DAYTONA_PLUGIN_DIR" ] && cd "$DAYTONA_PLUGIN_DIR" && npm install

echo "✨ macOS 安装与更新完成！"
echo "💡 运行 'bash start_opencode.sh' 启动服务。"
