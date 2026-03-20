const duckdb = require('duckdb');
const fs = require('fs');

// 1. 初始化 DuckDB 连接
const DB_PATH = '/opt/duckdb/agentflow.duckdb';

if (!fs.existsSync(DB_PATH)) {
    throw new Error(`DuckDB 数据库文件不存在：${DB_PATH}，请先执行初始化步骤`);
}

const db = new duckdb.Database(DB_PATH);

// 2. 全局缓存：model.id → AI SDK Provider 实例
const modelProviders = new Map();

// 3. 从 DuckDB 加载所有模型
async function loadModelsFromDuckDB() {
    return new Promise((resolve, reject) => {
        const sql = `SELECT id, name, provider, base_url, api_key, model_name, is_default FROM models`;
        db.all(sql, (err, rows) => {
            if (err) return reject(err);
            rows.forEach(row => {
                const { id, provider, base_url, api_key, model_name } = row;
                try {
                    if (provider === 'openai-compatible') {
                        const { createOpenAICompatible } = require('@ai-sdk/openai-compatible');
                        const providerInstance = createOpenAICompatible({
                            name: id,
                            baseURL: base_url,
                            apiKey: api_key || 'dummy',
                        });
                        modelProviders.set(id, {
                            provider: providerInstance,
                            modelName: model_name,
                            name: row.name,
                            isDefault: row.is_default
                        });
                        console.log(`[DuckDB] 加载模型: ${id} (${row.name}) → ${base_url}`);
                    }
                } catch (e) {
                    console.error(`[DuckDB] 模型 ${id} 加载失败`, e);
                }
            });
            resolve();
        });
    });
}

// 4. 获取要使用的模型
async function getTargetModel(userId, modelIdOverride) {
    if (modelIdOverride) {
        return new Promise(resolve => {
            db.get("SELECT * FROM models WHERE id = ?", [modelIdOverride], (err, row) => resolve(row));
        });
    }
    return new Promise(resolve => {
        db.get(`
            SELECT m.* FROM user_model_bindings b 
            JOIN models m ON b.model_id = m.id 
            WHERE b.user_id = ?
        `, [userId], (err, row) => {
            if (row) resolve(row);
            else db.get("SELECT * FROM models WHERE is_default = TRUE LIMIT 1", (err, def) => resolve(def));
        });
    });
}

// 5. 插件入口
module.exports = {
    name: 'duckdb-model-router',
    version: '1.0.0',
    async init() {
        await loadModelsFromDuckDB();
    },
    hooks: {
        'sse.message': async (context) => {
            const userId = context.request.headers['x-user-id'];
            const modelIdOverride = context.request.headers['x-model-id'];

            if (!userId && !modelIdOverride) return context;

            const targetModel = await getTargetModel(userId, modelIdOverride);
            if (targetModel && modelProviders.has(targetModel.id)) {
                const inst = modelProviders.get(targetModel.id);
                context.model = {
                    provider: inst.provider,
                    modelId: inst.modelName
                };
                
                // 注入调试信息，方便验证脚本识别
                context.response.messages.push({
                    role: 'system',
                    content: `[DEBUG] ROUTED_TO: ${targetModel.id} (${targetModel.model_name})`
                });
            }
            return context;
        }
    }
};
