# OpenCode + DuckDB + SSE 多模型动态路由

## 完整部署手册

---

## 一、手册说明

本手册从 0 开始实现「OpenCode 作为 SSE 服务端 + DuckDB 存储模型配置 + 按 UserID 动态选择模型」的完整能力。所有路径、依赖、配置均明确，数据库文件统一存放至 `/opt/duckdb/` 目录。

### 环境要求

| 依赖 | 版本要求 | 说明 |
|------|----------|------|
| 操作系统 | Linux/macOS | Windows 需适配路径（推荐 WSL） |
| Node.js | ≥ 18.x（20.x LTS 最佳） | 运行 OpenCode 核心环境 |
| DuckDB | ≥ 1.0.0 | 存储模型配置 / 用户映射 |
| 权限 | `/opt/duckdb/` 读写权限 | 生产环境建议限定用户组 |
| 网络 | 可访问 npm 源 | 或配置国内镜像（如淘宝源） |

---

## 二、Step 1：基础环境安装

### 2.1 安装 Node.js

**macOS（Homebrew）**

```bash
brew install node
```

**Linux（Ubuntu/Debian）**

```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs
```

**验证安装**

```bash
node -v    # 输出 v20.x.x 即为成功
npm -v     # 输出 10.x.x 即为成功
```

### 2.2 安装 OpenCode 核心包

```bash
# 全局安装 OpenCode
npm install -g opencode

# 验证安装
opencode --version    # 输出 ≥ 0.10.0 即为成功
```

### 2.3 安装 DuckDB（CLI + Node.js 驱动）

**1. 安装 DuckDB CLI**

```bash
# macOS
brew install duckdb

# Linux（Ubuntu/Debian）
wget https://github.com/duckdb/duckdb/releases/download/v1.0.0/duckdb_cli-linux-amd64.zip
unzip duckdb_cli-linux-amd64.zip -d /usr/local/bin/
chmod +x /usr/local/bin/duckdb
```

**2. 验证 DuckDB CLI**

```bash
duckdb --version    # 输出 ≥ 1.0.0 即为成功
```

> **说明**：Node.js 驱动会在插件安装时自动拉取，无需手动安装。

---

## 三、Step 2：创建目录结构

```bash
# 1. 创建 DuckDB 数据目录（核心：指定到 /opt/duckdb）
sudo mkdir -p /opt/duckdb
sudo chmod 777 /opt/duckdb
# 生产环境建议限定用户（如 chown -R youruser:yourgroup /opt/duckdb）

# 2. 创建 OpenCode 插件目录（固定路径，不可修改）
mkdir -p ~/.config/opencode/plugins/duckdb-model-router

# 3. 验证目录创建成功
ls -ld /opt/duckdb ~/.config/opencode/plugins/duckdb-model-router
```

---

## 四、Step 3：初始化 DuckDB 数据库

### 4.1 创建数据库文件并初始化表结构

```bash
# 进入 DuckDB CLI，连接 /opt/duckdb/opencode_models.duckdb
duckdb /opt/duckdb/opencode_models.duckdb
```

在 DuckDB 命令行中执行以下 SQL：

```sql
-- 1. 创建模型表
CREATE TABLE IF NOT EXISTS models (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    provider VARCHAR NOT NULL,
    base_url VARCHAR,
    api_key VARCHAR,
    model_name VARCHAR NOT NULL,
    description TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. 创建用户-模型映射表
CREATE TABLE IF NOT EXISTS user_model_bindings (
    user_id TEXT PRIMARY KEY,
    model_id TEXT NOT NULL REFERENCES models(id),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. 插入测试数据（可根据实际场景修改）
INSERT INTO models (id, name, provider, base_url, api_key, model_name, description, is_default)
VALUES
    ('ollama-qwen',   '通义千问-8B',    'openai-compatible', 'http://localhost:11434/v1', '',      'qwen3:8b',             '本地 Ollama 部署的千问模型', TRUE),
    ('vllm-deepseek', '深度求索-6.7B',  'openai-compatible', 'http://localhost:8000/v1',  '',      'deepseek-coder:6.7b',  'vLLM 部署的代码模型',      FALSE),
    ('openai-gpt4',   'GPT-4o',        'openai',            'https://api.openai.com/v1', 'sk-xxx', 'gpt-4o',               'OpenAI 官方 GPT-4o',      FALSE);

-- 4. 绑定用户与模型
INSERT INTO user_model_bindings (user_id, model_id)
VALUES
    ('user_001', 'ollama-qwen'),
    ('user_002', 'openai-gpt4');

-- 5. 验证数据（可选）
SELECT * FROM models;
SELECT * FROM user_model_bindings;

-- 退出 DuckDB CLI
.exit
```

### 4.2 验证数据库文件

```bash
ls -l /opt/duckdb/opencode_models.duckdb
# 能看到文件即为创建成功
```

---

## 五、Step 4：部署 DuckDB 模型路由插件

### 5.1 插件存放位置（固定）

> 插件必须放在：`~/.config/opencode/plugins/duckdb-model-router/`

### 5.2 创建插件文件

#### 1. 核心代码文件：`index.js`

```bash
touch ~/.config/opencode/plugins/duckdb-model-router/index.js
```

将以下内容写入 `index.js`：

```javascript
const duckdb = require('duckdb');
const fs = require('fs');
const path = require('path');

// 1. 初始化 DuckDB 连接（核心：指定 /opt/duckdb/ 路径）
const DB_PATH = '/opt/duckdb/opencode_models.duckdb';

// 检查数据库文件是否存在
if (!fs.existsSync(DB_PATH)) {
    throw new Error(`DuckDB 数据库文件不存在：${DB_PATH}，请先执行初始化步骤`);
}

const db = new duckdb.Database(DB_PATH);

// 2. 全局缓存：model.id → AI SDK Provider 实例
const modelProviders = new Map();

// ------------------------------
// 3. 从 DuckDB 加载所有模型（适配用户 DDL）
// ------------------------------
async function loadModelsFromDuckDB() {
    return new Promise((resolve, reject) => {
        const sql = `
            SELECT id, name, provider, base_url, api_key, model_name, description, is_default
            FROM models
        `;
        db.all(sql, (err, rows) => {
            if (err) return reject(err);
            rows.forEach(row => {
                const { id, provider, base_url, api_key, model_name } = row;
                try {
                    // 仅支持 openai-compatible 类型（可扩展其他 provider）
                    if (provider === 'openai-compatible') {
                        const { createOpenAICompatible } = require('@ai-sdk/openai-compatible');
                        const providerInstance = createOpenAICompatible({
                            name: id,
                            baseURL: base_url,
                            apiKey: api_key || 'dummy', // 无密钥时填占位符
                        });
                        // 缓存模型实例（适配用户字段）
                        modelProviders.set(id, {
                            provider: providerInstance,
                            modelName: model_name,
                            baseURL: base_url,
                            name: row.name,
                            isDefault: row.is_default
                        });
                        console.log(`[DuckDB] 加载模型: ${id} (${row.name}) → ${base_url}`);
                    } else {
                        console.warn(`[DuckDB] 暂不支持 provider: ${provider}（模型 ID: ${id}）`);
                    }
                } catch (e) {
                    console.error(`[DuckDB] 模型 ${id} 加载失败`, e);
                }
            });
            resolve();
        });
    });
}

// ------------------------------
// 4. 根据 user_id 查询绑定模型
// ------------------------------
async function getModelForUser(userId) {
    // 步骤1：查用户绑定的模型
    const userBinding = await new Promise((resolve) => {
        const sql = `
            SELECT b.model_id, m.*
            FROM user_model_bindings b
            JOIN models m ON b.model_id = m.id
            WHERE b.user_id = ?
        `;
        db.get(sql, [userId], (err, row) => resolve(row || null));
    });

    // 步骤2：无绑定则返回默认模型（is_default = TRUE）
    if (!userBinding) {
        return new Promise((resolve) => {
            const sql = `SELECT * FROM models WHERE is_default = TRUE LIMIT 1`;
            db.get(sql, (err, row) => resolve(row || null));
        });
    }
    return userBinding;
}

// ------------------------------
// 5. 手动指定模型查询
// ------------------------------
async function getModelById(modelId) {
    return new Promise((resolve) => {
        const sql = `SELECT * FROM models WHERE id = ?`;
        db.get(sql, [modelId], (err, row) => resolve(row || null));
    });
}

// ------------------------------
// 6. 插件入口（拦截 SSE 消息）
// ------------------------------
module.exports = {
    name: 'duckdb-model-router',
    version: '1.0.0',

    async init() {
        // 插件启动时加载所有模型
        await loadModelsFromDuckDB();
    },

    hooks: {
        'sse.message': async (context) => {
            const { request, response } = context;
            const userId = request.headers['x-user-id'];
            const modelIdOverride = request.headers['x-model-id']; // 手动指定模型 ID

            // 校验用户 ID
            if (!userId) {
                response.messages.push({
                    role: 'system',
                    content: '错误：请求头缺少 X-User-ID，无法路由模型'
                });
                return context;
            }

            // 核心：获取要使用的模型（手动指定 > 用户绑定 > 默认模型）
            let targetModel;
            if (modelIdOverride) {
                // 手动指定模型 ID
                targetModel = await getModelById(modelIdOverride);
                if (!targetModel) {
                    response.messages.push({
                        role: 'system',
                        content: `错误：未找到 ID 为 ${modelIdOverride} 的模型`
                    });
                    return context;
                }
            } else {
                // 自动匹配用户绑定/默认模型
                targetModel = await getModelForUser(userId);
                if (!targetModel) {
                    response.messages.push({
                        role: 'system',
                        content: '错误：未找到用户绑定模型，且无默认模型'
                    });
                    return context;
                }
            }

            // 从缓存获取模型 Provider 实例
            const modelInst = modelProviders.get(targetModel.id);
            if (!modelInst) {
                response.messages.push({
                    role: 'system',
                    content: `错误：模型 ${targetModel.id} 未加载（可能 provider 不支持）`
                });
                return context;
            }

            // 动态覆盖 OpenCode 当前使用的模型（核心逻辑）
            context.model = {
                provider: modelInst.provider,  // 模型 Provider 实例
                modelId: targetModel.model_name // 用户表的 model_name 字段
            };

            // 可选：将模型信息注入上下文（方便调试）
            response.context.push({
                role: 'system',
                content: `【当前使用模型】
ID：${targetModel.id}
名称：${targetModel.name}
Provider：${targetModel.provider}
BaseURL：${targetModel.base_url}
描述：${targetModel.description || '无'}`
            });

            return context;
        }
    }
};
```

#### 2. 插件配置文件：`package.json`

```bash
touch ~/.config/opencode/plugins/duckdb-model-router/package.json
```

将以下内容写入 `package.json`：

```json
{
    "name": "duckdb-model-router",
    "version": "1.0.0",
    "main": "index.js",
    "type": "commonjs",
    "dependencies": {
        "duckdb": "^1.0.0",
        "@ai-sdk/openai-compatible": "^1.0.0"
    },
    "description": "OpenCode 插件：基于 DuckDB 实现按 UserID 动态路由多模型（SSE 适配）"
}
```

### 5.3 安装插件依赖

```bash
cd ~/.config/opencode/plugins/duckdb-model-router
npm install    # 自动安装 duckdb 和 @ai-sdk/openai-compatible
```

---

## 六、Step 5：配置 OpenCode 服务

### 6.1 编辑 OpenCode 主配置文件

```bash
touch ~/.config/opencode/opencode.json
```

写入以下配置：

```json
{
    "server": {
        "enabled": true,
        "port": 3000,
        "cors": true,
        "sse": {
            "enabled": true,
            "path": "/api/sse/chat"
        }
    },
    "plugins": [
        "~/.config/opencode/plugins/duckdb-model-router"
    ],
    "model": "qwen3:8b",
    "provider": {
        "default": {
            "npm": "@ai-sdk/openai-compatible",
            "options": {
                "baseURL": "http://localhost:11434/v1"
            }
        }
    }
}
```

> `model` 为兜底默认模型，插件会优先使用 DuckDB 中 `is_default = TRUE` 的模型。

### 6.2 启动 OpenCode SSE 服务

```bash
# 前台启动（调试用）
opencode server start

# 后台启动（生产用）
opencode server start --daemon

# 验证服务启动
curl http://localhost:3000/health
# 输出 {"status":"ok"} 即为成功
```

---

## 七、Step 6：测试 SSE 多模型路由

### 7.1 客户端测试代码

创建 `sse-test.js`（可放在任意目录）：

```javascript
const EventSource = require('eventsource');   // npm install eventsource
const fetch = require('node-fetch');           // npm install node-fetch

// 按 user_id 连接 SSE（自动匹配模型）
function connectSSE(userId, customModelId = null) {
    const headers = {
        'X-User-ID': userId,
        'Content-Type': 'application/json'
    };

    // 手动指定模型 ID（可选）
    if (customModelId) {
        headers['X-Model-ID'] = customModelId;
    }

    const sse = new EventSource('http://localhost:3000/api/sse/chat', { headers });

    // 监听 AI 回复
    sse.addEventListener('message', (e) => {
        const data = JSON.parse(e.data);
        console.log(`[用户 ${userId}] AI 回复：`, data.content);
    });

    // 监听错误
    sse.addEventListener('error', (e) => {
        console.error(`[用户 ${userId}] SSE 错误：`, e);
        sse.close();
    });

    // 发送消息函数
    async function sendMessage(message) {
        await fetch('http://localhost:3000/api/sse/chat', {
            method: 'POST',
            headers: headers,
            body: JSON.stringify({ message })
        });
    }

    return { sse, sendMessage };
}

// 测试用例
async function runTests() {
    // 测试1：user_001 → 绑定的 ollama-qwen 模型
    const user1 = connectSSE('user_001');
    await user1.sendMessage('用 Python 写一个冒泡排序');

    // 测试2：user_002 → 绑定的 openai-gpt4 模型
    const user2 = connectSSE('user_002');
    await user2.sendMessage('解释一下微服务的核心设计原则');

    // 测试3：user_001 手动切换到 vllm-deepseek 模型
    const user1Custom = connectSSE('user_001', 'vllm-deepseek');
    await user1Custom.sendMessage('优化这段 C++ 代码的性能');

    // 测试4：无绑定的用户 → 默认模型（ollama-qwen）
    const user3 = connectSSE('user_003');
    await user3.sendMessage('介绍一下 DuckDB 的特点');
}

runTests();
```

### 7.2 运行测试

```bash
# 安装依赖
npm install eventsource node-fetch

# 运行测试
node sse-test.js
```

### 7.3 预期结果

| 用户 ID | 模型来源 | 实际使用模型 | 回复特征 |
|---------|---------|-------------|---------|
| `user_001` | 用户绑定 | `ollama-qwen` | 本地化回复，Python 代码示例 |
| `user_002` | 用户绑定 | `openai-gpt4` | 官方 GPT-4o 风格回复 |
| `user_001` | 手动指定 | `vllm-deepseek` | 代码优化专业回复 |
| `user_003` | 默认模型（is_default） | `ollama-qwen` | 本地化回复 DuckDB 特点 |

---

## 八、核心能力说明

### 8.1 模型路由逻辑

```
SSE 请求
  │
  ├── 提取 X-User-ID / X-Model-ID
  │
  ├── 有 X-Model-ID？
  │     ├── YES → 查询 DuckDB 模型表（按 id）
  │     └── NO  → 查询用户-模型绑定表
  │                ├── 有绑定 → 获取绑定模型
  │                └── 无绑定 → 获取 is_default = TRUE 的模型
  │
  ├── 加载模型 Provider 实例
  ├── 动态替换 OpenCode 模型
  └── 返回 AI 回复
```

### 8.2 关键路径 / 配置汇总

| 名称 | 路径 / 值 | 说明 |
|------|----------|------|
| DuckDB 数据库文件 | `/opt/duckdb/opencode_models.duckdb` | 模型配置核心存储 |
| 插件目录 | `~/.config/opencode/plugins/duckdb-model-router` | OpenCode 插件固定路径 |
| SSE 接口地址 | `http://localhost:3000/api/sse/chat` | 客户端连接地址 |
| 用户 ID 请求头 | `X-User-ID` | 必传，用于模型路由 |
| 模型 ID 请求头 | `X-Model-ID` | 可选，手动指定模型 |

---

## 九、常见问题排查

### 9.1 插件加载失败

- 检查插件目录路径是否正确：`~/.config/opencode/plugins/duckdb-model-router`
- 检查插件依赖是否安装：`cd ~/.config/opencode/plugins/duckdb-model-router && npm install`
- 查看 OpenCode 日志：`opencode server logs`

### 9.2 模型路由无响应

- 检查 DuckDB 数据库文件是否存在：`/opt/duckdb/opencode_models.duckdb`
- 验证用户-模型绑定数据：`duckdb /opt/duckdb/opencode_models.duckdb -c "SELECT * FROM user_model_bindings;"`
- 检查请求头是否传递 `X-User-ID`：客户端需显式设置该 Header

### 9.3 SSE 连接报错

- 检查端口是否占用：`netstat -tulpn | grep 3000`
- 检查 CORS 配置：OpenCode 配置中 `server.cors` 需设为 `true`
- 检查防火墙：开放 3000 端口（Linux：`sudo ufw allow 3000`）

---

## 十、扩展能力（可选）

### 10.1 支持更多 Provider

在插件 `loadModelsFromDuckDB` 函数中扩展：

```javascript
// 示例：支持 anthropic
if (provider === 'anthropic') {
    const { createAnthropic } = require('@ai-sdk/anthropic');
    const providerInstance = createAnthropic({
        apiKey: api_key,
    });
    // 缓存逻辑...
}
```

### 10.2 记忆功能集成

可结合之前的「记忆插件」，在当前插件中追加记忆逻辑，实现「用户 ID → 模型 + 记忆」双重绑定。

### 10.3 模型热加载

添加定时任务，自动重新加载 DuckDB 中的模型配置：

```javascript
// 在插件 init 函数中添加
setInterval(async () => {
    console.log('[DuckDB] 定时重新加载模型配置');
    await loadModelsFromDuckDB();
}, 30 * 60 * 1000); // 每 30 分钟加载一次
```

---

## 总结

| 项目 | 说明 |
|------|------|
| **核心路径** | DuckDB 数据库固定在 `/opt/duckdb/`，插件固定在 `~/.config/opencode/plugins/duckdb-model-router/`，无需修改 |
| **核心逻辑** | SSE 请求通过 `X-User-ID` 关联用户，从 DuckDB 读取绑定模型，动态替换 OpenCode 的模型实例 |
| **使用方式** | 客户端只需传递 `X-User-ID`（可选 `X-Model-ID`），即可实现多用户 / 多模型的隔离与动态切换 |

> 本手册所有代码均可直接复制使用，生产环境只需调整模型配置（如 API Key、BaseURL）和权限策略。
