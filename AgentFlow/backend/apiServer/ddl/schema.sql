-- Models table
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

-- Skills table
CREATE TABLE IF NOT EXISTS skills (
    id VARCHAR PRIMARY KEY, 
    name VARCHAR, 
    type VARCHAR, 
    description TEXT, 
    version VARCHAR, 
    icon VARCHAR, 
    color VARCHAR, 
    logic TEXT, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Agents table
CREATE TABLE IF NOT EXISTS agents (
    id VARCHAR PRIMARY KEY, 
    name VARCHAR, 
    description TEXT, 
    model VARCHAR, 
    system_prompt TEXT, 
    thinking_enabled BOOLEAN DEFAULT TRUE,
    simplified_output BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Edges table
CREATE TABLE IF NOT EXISTS edges (
    id VARCHAR PRIMARY KEY, 
    from_id VARCHAR, 
    to_id VARCHAR, 
    type VARCHAR, 
    metadata JSON, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Groups table
CREATE TABLE IF NOT EXISTS groups (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    group_rule_mode VARCHAR DEFAULT 'free',
    thinking_enabled BOOLEAN DEFAULT TRUE,
    simplified_output BOOLEAN DEFAULT FALSE,
    custom_rule TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Group members table
CREATE TABLE IF NOT EXISTS group_members (
    group_id VARCHAR NOT NULL,
    user_id VARCHAR NOT NULL,
    PRIMARY KEY (group_id, user_id)
);

-- Agent skills table
CREATE TABLE IF NOT EXISTS agent_skills (
    agent_id VARCHAR NOT NULL,
    skill_id VARCHAR NOT NULL,
    PRIMARY KEY (agent_id, skill_id)
);

-- Workflows table
CREATE TABLE IF NOT EXISTS workflows (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT,
    graph JSON, -- 存储节点和连线的 JSON 结构
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
