import { Router } from 'express';
import prisma from '../db/sqlite';
import { runQuery } from '../db/memgraph';

const router = Router();

// GET all skills
router.get('/', async (req, res) => {
  try {
    const skills = await prisma.skill.findMany({
      orderBy: { createdAt: 'desc' }
    });
    res.json(skills);
  } catch (error) {
    res.status(500).json({ error: 'Failed to fetch skills' });
  }
});

// POST - Add new skill (Manual)
router.get('/add', async (req, res) => {
  // Use GET for simple demonstration if needed, but normally POST
});

router.post('/', async (req, res) => {
  const { name, type, description, logic, version, icon, color } = req.body;
  try {
    const skill = await prisma.skill.create({
      data: { name, type, description, logic, version, icon, color }
    });

    // Sync to Memgraph
    await runQuery(
      'CREATE (s:Skill {id: $id, name: $name, type: $type})',
      { id: skill.id, name: skill.name, type: skill.type }
    );

    res.status(201).json(skill);
  } catch (error) {
    res.status(500).json({ error: 'Failed to create skill' });
  }
});

// POST - Install skill (via URL)
router.post('/install', async (req, res) => {
  const { url } = req.body;
  try {
    // Mocking the installation process
    console.log(`Fetching skill from ${url}...`);
    
    // In a real app, you would fetch(url) and parse the JSON
    const installedSkill = await prisma.skill.create({
      data: {
        name: 'Remote Skill',
        type: 'HTTP API',
        description: `Successfully installed from ${url}`,
        version: 'v1.0.0',
        icon: 'lucide:link',
        logic: '# Remote installation logic'
      }
    });

    await runQuery(
      'CREATE (s:Skill {id: $id, name: $name, type: $type, installed_from: $url})',
      { id: installedSkill.id, name: installedSkill.name, type: installedSkill.type, url }
    );

    res.json(installedSkill);
  } catch (error) {
    res.status(500).json({ error: 'Installation failed' });
  }
});

// DELETE skill
router.delete('/:id', async (req, res) => {
  const { id } = req.params;
  try {
    await prisma.skill.delete({ where: { id } });
    
    // Remove from Memgraph
    await runQuery('MATCH (s:Skill {id: $id}) DETACH DELETE s', { id });
    
    res.status(204).send();
  } catch (error) {
    res.status(500).json({ error: 'Failed to delete skill' });
  }
});

export default router;
