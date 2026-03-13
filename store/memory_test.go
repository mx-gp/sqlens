package store

import (
	"testing"
	"time"

	"github.com/sqlens/sqlens/analyzer"
)

func TestMemoryStore(t *testing.T) {
	s := NewMemoryStore()

	// Save 10 events
	for i := 0; i < 10; i++ {
		s.Save(analyzer.QueryEvent{
			RawQuery:  "SELECT 1",
			N1Flag:    i%2 == 0, // Even indexes have N1Flag
			Timestamp: time.Now(),
		})
	}

	// Test GetRecent
	recent := s.GetRecent(5)
	if len(recent) != 5 {
		t.Errorf("Expected 5 recent events, got %d", len(recent))
	}

	// Test GetN1Incidents
	// All 10 events are within the last minute
	incidents := s.GetN1Incidents(time.Minute)
	if len(incidents) != 5 {
		t.Errorf("Expected 5 N1 incidents, got %d", len(incidents))
	}
}
