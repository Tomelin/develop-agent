package interview

import "sync"

type Event struct {
	ProjectID string `json:"project_id"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}

type Broker struct {
	mu   sync.RWMutex
	subs map[string][]chan Event
}

func NewBroker() *Broker {
	return &Broker{subs: make(map[string][]chan Event)}
}

func (b *Broker) Subscribe(projectID string) <-chan Event {
	ch := make(chan Event, 8)
	b.mu.Lock()
	b.subs[projectID] = append(b.subs[projectID], ch)
	b.mu.Unlock()
	return ch
}

func (b *Broker) Emit(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subs[e.ProjectID] {
		select {
		case ch <- e:
		default:
		}
	}
}
