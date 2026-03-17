# AgentFlow Backend (Unified DuckDB)

基于 TypeScript + Express 的后端服务。

## 技术栈
- **Unified Database**: DuckDB (for both Relational and Graph data)
- **API**: Express

## 快速启动

1. **安装依赖**:
   ```bash
   npm install
   ```

2. **启动开发服务器**:
   ```bash
   npm run dev
   ```
   *(数据库将在首次启动时自动在 `data/agentflow.duckdb` 创建)*

## 目录说明
- `/src/db/duckdb.ts`: 数据库连接、查询及初始化逻辑
- `/src/routes`: API 路由实现
