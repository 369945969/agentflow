import neo4j, { Driver } from 'neo4j-driver';
import dotenv from 'dotenv';

dotenv.config();

let driver: Driver;

export const getMemgraphDriver = () => {
  if (!driver) {
    const uri = process.env.MEMGRAPH_URI || 'bolt://localhost:7687';
    const user = process.env.MEMGRAPH_USER || '';
    const password = process.env.MEMGRAPH_PASSWORD || '';
    
    driver = neo4j.driver(uri, neo4j.auth.basic(user, password));
  }
  return driver;
};

export const runQuery = async (query: string, params: any = {}) => {
  const driver = getMemgraphDriver();
  const session = driver.session();
  try {
    const result = await session.run(query, params);
    return result;
  } finally {
    await session.close();
  }
};
