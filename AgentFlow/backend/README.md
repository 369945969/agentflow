# AgentFlow Backend

基于 TypeScript + Express 的后端服务。

## 技术栈
- **Relational DB**: SQLite (via Prisma)
- **Graph DB**: Memgraph (via Neo4j Driver)
- **API**: Express

## 快速启动

1. **安装依赖**:
   ```bash
   npm install
   ```

2. **初始化数据库**:
   ```bash
   # 生成 Prisma 客户端并同步 SQLite
   npx prisma migrate dev --name init
   ```

3. **运行 Memgraph (Docker)**:
   ```bash
   docker run -p 7687:7687 -p 7444:7444 memgraph/memgraph
   ```

4. **启动开发服务器**:
   ```bash
   npm run dev
   ```

## 目录说明
- `/ddl`: 数据库初始化脚本 (SQL & Cypher)
- `/src/db`: 数据库连接配置
- `/src/routes`: API 路由实现
- `/prisma`: SQLite 模式定义
