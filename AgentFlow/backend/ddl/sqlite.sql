-- SQLite DDL
-- Table: Skills
CREATE TABLE IF NOT EXISTS skills (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- 'HTTP API', 'Script', 'Database', 'Plugin'
    description TEXT,
    version TEXT DEFAULT 'v1.0.0',
    icon TEXT DEFAULT 'lucide:zap',
    color TEXT DEFAULT '#3B9BFF',
    logic TEXT, -- Python code or JSON config
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Table: Agents
CREATE TABLE IF NOT EXISTS agents (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    model TEXT,
    system_prompt TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
