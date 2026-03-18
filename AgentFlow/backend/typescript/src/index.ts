import express from 'express';
import cors from 'cors';
import dotenv from 'dotenv';
import skillRoutes from './routes/skills.js';
import { initDb } from './db/duckdb.js';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3000;

app.use(cors());
app.use(express.json());

app.use('/api/skills', skillRoutes);

app.get('/health', (req, res) => {
  res.json({ status: 'ok', message: 'Backend server is running' });
});

const startServer = async () => {
  await initDb();
  app.listen(PORT, () => {
    console.log(`🚀 Server is listening at http://localhost:${PORT}`);
  });
};

startServer();