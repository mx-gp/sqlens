package analyzer

import (
	"context"
	"regexp"
)

var (
	// Very basic regex to replace literals with placeholders
	numberRegex = regexp.MustCompile(`\b\d+\b`)
	stringRegex = regexp.MustCompile(`'[^']*'`)
)

type FingerprintAnalyzer struct{}

func NewFingerprintAnalyzer() *FingerprintAnalyzer {
	return &FingerprintAnalyzer{}
}

func (f *FingerprintAnalyzer) Analyze(ctx context.Context, event QueryEvent) (*AnalysisResult, error) {
	fp := stringRegex.ReplaceAllString(event.RawQuery, "?")
	fp = numberRegex.ReplaceAllString(fp, "?")
	
	event.Fingerprint = fp
	return &AnalysisResult{Event: event}, nil
}
