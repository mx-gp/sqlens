package analyzer

import (
	"context"
	"time"
)

type QueryEvent struct {
	ConnectionID string
	RawQuery     string
	Latency      time.Duration
	Timestamp    time.Time
	Fingerprint  string
	N1Flag       bool
	Violations   []string // New: Performance/Safety guardrail warnings
}

type AnalysisResult struct {
	Event QueryEvent
}

type Analyzer interface {
	Analyze(ctx context.Context, event QueryEvent) (*AnalysisResult, error)
}

type Pipeline struct {
	analyzers []Analyzer
}

func NewPipeline(analyzers ...Analyzer) *Pipeline {
	return &Pipeline{analyzers: analyzers}
}

func (p *Pipeline) Process(ctx context.Context, event QueryEvent) {
	currentEvent := event
	for _, a := range p.analyzers {
		result, err := a.Analyze(ctx, currentEvent)
		if err != nil {
			// In a real app we'd log this, but we continue processing
			continue
		}
		if result != nil {
			currentEvent = result.Event
		}
	}
}
