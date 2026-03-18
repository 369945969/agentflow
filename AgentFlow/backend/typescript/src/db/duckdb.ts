import duckdb from 'duckdb';
import path from 'path';
import fs from 'fs';

const dbPath = path.resolve(process.cwd(), './data/agentflow.duckdb');

const dataDir = path.dirname(dbPath);
if (!fs.existsSync(dataDir)) {
    fs.mkdirSync(dataDir, { recursive: true });
}

const db = new duckdb.Database(dbPath);

export const query = (sql: string, params: any[] = []): Promise<any[]> => {
    return new Promise((resolve, reject) => {
        db.all(sql, ...params, (err: any, res: any) => {
            if (err) reject(err);
            else resolve(res);
        });
    });
};

export const exec = (sql: string, params: any[] = []): Promise<void> => {
    return new Promise((resolve, reject) => {
        db.run(sql, ...params, (err: any) => {
            if (err) reject(err);
            else resolve();
        });
    });
};

export const initDb = async () => {
    await exec(`CREATE TABLE IF NOT EXISTS skills (id VARCHAR PRIMARY KEY, name VARCHAR, type VARCHAR, description TEXT, version VARCHAR, icon VARCHAR, color VARCHAR, logic TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`);
    await exec(`CREATE TABLE IF NOT EXISTS agents (id VARCHAR PRIMARY KEY, name VARCHAR, description TEXT, model VARCHAR, system_prompt TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`);
    await exec(`CREATE TABLE IF NOT EXISTS edges (id VARCHAR PRIMARY KEY, from_id VARCHAR, to_id VARCHAR, type VARCHAR, metadata JSON, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`);
    console.log('✅ DuckDB Initialized (Relational + Graph)');
};

export default db;