package triad

import (
	"context"
	"sync"
)

// Broker keeps an in-memory pub/sub for triad execution events.
type Broker struct {
	mu          sync.RWMutex
	subscribers map[string][]chan Event
}

func NewBroker() *Broker {
	return &Broker{subscribers: make(map[string][]chan Event)}
}

func (b *Broker) Subscribe(projectID string, buffer int) <-chan Event {
	if buffer <= 0 {
		buffer = 16
	}
	ch := make(chan Event, buffer)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[projectID] = append(b.subscribers[projectID], ch)
	return ch
}

func (b *Broker) Emit(_ context.Context, event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, sub := range b.subscribers[event.ProjectID] {
		select {
		case sub <- event:
		default:
		}
	}
}

func (b *Broker) CloseProject(projectID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	subs := b.subscribers[projectID]
	for _, sub := range subs {
		close(sub)
	}
	delete(b.subscribers, projectID)
}
