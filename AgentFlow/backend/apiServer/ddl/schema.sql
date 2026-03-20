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
