CREATE TABLE IF NOT EXISTS models (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    provider VARCHAR NOT NULL, -- e.g., OpenAI, Anthropic, Local
    base_url VARCHAR,           -- API base URL
    api_key VARCHAR,            -- API Key (should be handled securely)
    model_name VARCHAR NOT NULL, -- The specific model identifier (e.g., gpt-4)
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
