# AgentFlow WebSocket Gateway (Golang)

这是一个基于 Go 语言实现的 WebSocket 代理服务端，专门用于在 **WebSocket 协议**与 **OpenCode SSE (Server-Sent Events) 接口**之间建立双向流式通信桥梁。

## 核心功能

1.  **协议转换**: 将前端的 WebSocket 消息转换为带自定义 Header（`X-User-ID`, `X-Model-ID`）的 HTTP POST 请求发送至 OpenCode。
2.  **流式转发**: 实时捕获 OpenCode SSE 输出流，并将其透明地转发给 WebSocket 客户端。
3.  **动态路由支持**: 允许客户端在 WebSocket 握手后的 JSON 消息中指定 `userId`，由 OpenCode 插件根据该 ID 从 DuckDB 中自动匹配模型。
4.  **异步并发**: 采用 Go 协程处理每个流式请求，支持高并发连接而不阻塞主线程。

## 技术栈

-   **语言**: Golang 1.20+
-   **库**: `github.com/gorilla/websocket` (工业级 WebSocket 实现)
-   **后端**: 转发至 OpenCode 服务 (`:3001`)

## 接口规范

### 连接地址
`ws://localhost:3002/ws`

### 客户端发送消息格式 (JSON)
```json
{
  "userId": "user_001",      // 必须，用于 DuckDB 路由
  "modelId": "custom_model", // 可选，强制覆盖模型
  "message": "你好，请介绍一下你自己" // 必须，用户提示词
}
```

### 服务端返回消息格式 (JSON)

#### 1. 流式内容 (Delta)
```json
{
  "type": "delta",
  "raw": {
    "content": "具体的回复文本片段...",
    "role": "assistant"
  }
}
```

#### 2. 完成信号 (Done)
```json
{
  "type": "done"
}
```

#### 3. 错误信息 (Error)
```json
{
  "type": "error",
  "error": "错误描述信息"
}
```

## 目录结构说明

-   `main.go`: 服务端核心逻辑，处理 WebSocket 升级、HTTP 请求转发及流式解析。
-   `ws_test.go`: 集成测试代码，模拟真实客户端进行端到端验证。
-   `start.sh`: 生产/开发环境启动脚本，支持后台运行和日志重定向。
-   `verify_ws.sh`: 自动化测试触发脚本。
-   `ws_server.log`: 服务运行日志。

## 快速开始

### 1. 安装依赖
```bash
go mod tidy
```

### 2. 启动服务
```bash
bash start.sh
```

### 3. 执行自动化测试
```bash
bash verify_ws.sh
```

## 注意事项

-   运行前请确保 **OpenCode 服务** 已在 `3001` 端口启动。
-   如需修改 OpenCode 的地址，请编辑 `main.go` 中的 `openCodeURL` 变量。
-   WebSocket 服务默认监听 `3002` 端口。
