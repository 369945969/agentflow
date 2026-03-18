import { Router } from 'express';
import { query, exec } from '../db/duckdb.js';
import { randomUUID } from 'crypto';

const router = Router();

router.get('/', async (req, res) => {
  try {
    const skills = await query('SELECT * FROM skills ORDER BY created_at DESC');
    res.json(skills);
  } catch (error) {
    res.status(500).json({ error: 'Failed to fetch skills' });
  }
});

router.post('/', async (req, res) => {
  const { name, type, description, logic, version, icon, color } = req.body;
  const id = randomUUID();
  try {
    await exec('INSERT INTO skills (id, name, type, description, logic, version, icon, color) VALUES (?, ?, ?, ?, ?, ?, ?, ?)', [id, name, type, description, logic, version, icon, color]);
    const newSkill = await query('SELECT * FROM skills WHERE id = ?', [id]);
    res.status(201).json(newSkill[0]);
  } catch (error) {
    res.status(500).json({ error: 'Failed to create skill' });
  }
});

router.post('/install', async (req, res) => {
  const { url } = req.body;
  const id = randomUUID();
  try {
    const skillData = { id, name: 'Remote Skill', type: 'HTTP API', description: `Installed from ${url}`, version: 'v1.0.0', icon: 'lucide:link', logic: `# Fetched from ${url}` };
    await exec('INSERT INTO skills (id, name, type, description, version, icon, logic) VALUES (?, ?, ?, ?, ?, ?, ?)', [skillData.id, skillData.name, skillData.type, skillData.description, skillData.version, skillData.icon, skillData.logic]);
    const newSkill = await query('SELECT * FROM skills WHERE id = ?', [id]);
    res.json(newSkill[0]);
  } catch (error) {
    res.status(500).json({ error: 'Installation failed' });
  }
});

router.delete('/:id', async (req, res) => {
  const { id } = req.params;
  try {
    await exec('DELETE FROM edges WHERE source_id = ? OR target_id = ?', [id, id]);
    await exec('DELETE FROM skills WHERE id = ?', [id]);
    res.status(204).send();
  } catch (error) {
    res.status(500).json({ error: 'Failed to delete skill' });
  }
});

export default router;