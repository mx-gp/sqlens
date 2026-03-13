package analyzer

import (
	"context"
	"testing"
	"time"
)

func TestN1DetectorAnalyzer(t *testing.T) {
	// 5 second window, 3 hits threshold
	n1 := NewN1DetectorAnalyzer(5*time.Second, 3)
	ctx := context.Background()
	connID := "test-conn"
	fp := "SELECT * FROM users WHERE id = ?"

	// First two calls should NOT flag N+1
	for i := 1; i <= 2; i++ {
		event := QueryEvent{ConnectionID: connID, Fingerprint: fp}
		result, _ := n1.Analyze(ctx, event)
		if result.Event.N1Flag {
			t.Errorf("Iteration %d: expected N1Flag false, got true", i)
		}
	}

	// Third call should flag N+1
	event := QueryEvent{ConnectionID: connID, Fingerprint: fp}
	result, _ := n1.Analyze(ctx, event)
	if !result.Event.N1Flag {
		t.Error("Third hit: expected N1Flag true, got false")
	}

	// After flag, it should reset for that fingerprint
	event2 := QueryEvent{ConnectionID: connID, Fingerprint: fp}
	result2, _ := n1.Analyze(ctx, event2)
	if result2.Event.N1Flag {
		t.Error("After reset: expected N1Flag false, got true")
	}
}
