const { execSync } = require('child_process');

// 缓存已创建的工作区 ID
const activeWorkspaces = new Set();

async function ensureWorkspace(sessionId) {
    const workspaceId = `opencode-${sessionId}`;
    if (activeWorkspaces.has(workspaceId)) return workspaceId;

    try {
        console.log(`[Daytona] 正在为会话 ${sessionId} 创建沙箱工作区...`);
        // 检查工作区是否已存在
        const list = execSync('daytona list').toString();
        if (!list.includes(workspaceId)) {
            // 使用默认的 starter 模板创建工作区
            // 注意：这里假设 daytona 已经配置好默认 target
            execSync(`daytona create ${workspaceId} --non-interactive`);
        }
        activeWorkspaces.add(workspaceId);
        return workspaceId;
    } catch (error) {
        console.error(`[Daytona] 创建工作区失败: ${error.message}`);
        return null;
    }
}

module.exports = function(pluginContext) {
    return {
        name: 'daytona-sandbox',
        version: '1.0.0',
        hooks: {
            'sse.message': async (context) => {
                const sessionId = context.request.headers['x-session-id'] || context.sessionID;
                if (!sessionId) return context;

                const workspaceId = await ensureWorkspace(sessionId);
                if (workspaceId) {
                    // 告诉 OpenCode 后续的代码执行应该在这个 Daytona 工作区中进行
                    // 注意：这需要 OpenCode 的执行引擎支持从 context.execution 获取 target
                    context.execution = {
                        type: 'daytona',
                        workspaceId: workspaceId
                    };
                    
                    if (context.response && context.response.messages) {
                        context.response.messages.push({
                            role: 'system',
                            content: `[SANDBOX] 隔离环境已就绪: ${workspaceId}`
                        });
                    }
                }
                return context;
            }
        }
    };
};
