package analyzer

import (
	"context"
	"sync"
	"time"
)

type N1DetectorAnalyzer struct {
	window    time.Duration
	threshold int
	history   map[string]map[string][]time.Time // connID -> fingerprint -> timestamps
	mu        sync.Mutex
}

func NewN1DetectorAnalyzer(window time.Duration, threshold int) *N1DetectorAnalyzer {
	return &N1DetectorAnalyzer{
		window:    window,
		threshold: threshold,
		history:   make(map[string]map[string][]time.Time),
	}
}

func (n *N1DetectorAnalyzer) Analyze(ctx context.Context, event QueryEvent) (*AnalysisResult, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	connID := event.ConnectionID
	fp := event.Fingerprint

	if n.history[connID] == nil {
		n.history[connID] = make(map[string][]time.Time)
	}

	now := time.Now()
	times := n.history[connID][fp]
	
	// Clean up old entries outside the window
	valid := times[:0]
	for _, t := range times {
		if now.Sub(t) <= n.window {
			valid = append(valid, t)
		}
	}
	
	valid = append(valid, now)
	n.history[connID][fp] = valid

	if len(valid) >= n.threshold {
		event.N1Flag = true
		// Reset to avoid spamming alerts for the same window
		n.history[connID][fp] = nil
	}

	return &AnalysisResult{Event: event}, nil
}
