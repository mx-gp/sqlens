package store

import (
	"sync"
	"time"

	"github.com/sqlens/sqlens/analyzer"
)

type MemoryStore struct {
	mu     sync.RWMutex
	events []analyzer.QueryEvent
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		events: make([]analyzer.QueryEvent, 0),
	}
}

func (m *MemoryStore) Save(event analyzer.QueryEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
	
	// Keep the last 1000 events
	if len(m.events) > 1000 {
		m.events = m.events[1:]
	}
}

func (m *MemoryStore) GetRecent(limit int) []analyzer.QueryEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if limit > len(m.events) {
		limit = len(m.events)
	}
	
	result := make([]analyzer.QueryEvent, limit)
	copy(result, m.events[len(m.events)-limit:])
	return result
}

// GetN1Incidents returns all events that have the N1Flag set
func (m *MemoryStore) GetN1Incidents(since time.Duration) []analyzer.QueryEvent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var incidents []analyzer.QueryEvent
	cutoff := time.Now().Add(-since)
	
	for _, e := range m.events {
		if e.N1Flag && e.Timestamp.After(cutoff) {
			incidents = append(incidents, e)
		}
	}
	return incidents
}
