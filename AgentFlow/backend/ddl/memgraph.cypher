// Memgraph Cypher DDL (Graph Schema)

// Create Constraints
CREATE CONSTRAINT ON (s:Skill) ASSERT s.id IS UNIQUE;
CREATE CONSTRAINT ON (a:Agent) ASSERT a.id IS UNIQUE;

// Example nodes and relationships setup
// (:Agent {id: 'agent-1'})-[:USES_SKILL {installed_at: datetime()}]->(:Skill {id: 'skill-1'})
// (:Skill {id: 'skill-1'})-[:DEPENDS_ON]->(:Skill {id: 'skill-2'})
