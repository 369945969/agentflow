package agentcore

import "sync"

type SubAgentPool struct {
	mu     sync.RWMutex
	agents map[string]SubAgentConfig
	order  []string
}

func NewSubAgentPool(agents []SubAgentConfig) *SubAgentPool {
	p := &SubAgentPool{
		agents: make(map[string]SubAgentConfig),
	}
	for _, a := range agents {
		if a.Endpoint == "" || a.AgentID == "" {
			continue
		}
		p.agents[a.AgentID] = a
		p.order = append(p.order, a.AgentID)
	}
	return p
}

func (p *SubAgentPool) Default() (SubAgentConfig, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if len(p.order) == 0 {
		return SubAgentConfig{}, false
	}
	a, ok := p.agents[p.order[0]]
	return a, ok
}

func (p *SubAgentPool) Get(agentID string) (SubAgentConfig, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	a, ok := p.agents[agentID]
	return a, ok
}

func (p *SubAgentPool) GetAll() []SubAgentConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]SubAgentConfig, 0, len(p.order))
	for _, id := range p.order {
		a, ok := p.agents[id]
		if ok {
			out = append(out, a)
		}
	}
	return out
}
