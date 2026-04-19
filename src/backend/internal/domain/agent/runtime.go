package agent

import "sync"

// Runtime holds channels and state for an active agent.
type Runtime struct {
	ID     string
	In     chan AgentMessage
	Out    chan AgentMessage
	Status Status
	mu     sync.RWMutex
	closed bool
}

func (r *Runtime) SetStatus(status Status) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Status = status
}

func (r *Runtime) CurrentStatus() Status {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Status
}

func (r *Runtime) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return
	}
	close(r.In)
	close(r.Out)
	r.closed = true
}

// ChannelRegistry controls runtime channel lifecycle.
type ChannelRegistry struct {
	mu         sync.RWMutex
	bufferSize int
	runtimes   map[string]*Runtime
}

func NewChannelRegistry(bufferSize int) *ChannelRegistry {
	if bufferSize <= 0 {
		bufferSize = 10
	}
	return &ChannelRegistry{bufferSize: bufferSize, runtimes: make(map[string]*Runtime)}
}

func (c *ChannelRegistry) Create(agentID string) *Runtime {
	c.mu.Lock()
	defer c.mu.Unlock()

	rt := &Runtime{
		ID:     agentID,
		In:     make(chan AgentMessage, c.bufferSize),
		Out:    make(chan AgentMessage, c.bufferSize),
		Status: StatusIdle,
	}
	c.runtimes[agentID] = rt
	return rt
}

func (c *ChannelRegistry) Get(agentID string) (*Runtime, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	rt, ok := c.runtimes[agentID]
	return rt, ok
}

func (c *ChannelRegistry) Destroy(agentID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	rt, ok := c.runtimes[agentID]
	if !ok {
		return
	}
	rt.Close()
	delete(c.runtimes, agentID)
}
