package analyzer

import (
	"context"
	"strings"
)

type GuardrailAnalyzer struct{}

func NewGuardrailAnalyzer() *GuardrailAnalyzer {
	return &GuardrailAnalyzer{}
}

func (g *GuardrailAnalyzer) Analyze(ctx context.Context, event QueryEvent) (*AnalysisResult, error) {
	upperQuery := strings.ToUpper(event.RawQuery)
	
	// Guardrail 1: DELETE/UPDATE without WHERE (Dangerous!)
	if (strings.Contains(upperQuery, "DELETE") || strings.Contains(upperQuery, "UPDATE")) && 
	   !strings.Contains(upperQuery, "WHERE") {
		event.Violations = append(event.Violations, "DANGEROUS: Modification query without WHERE clause detected!")
	}

	// Guardrail 2: SELECT * (Performance Antipattern)
	if strings.Contains(upperQuery, "SELECT *") {
		event.Violations = append(event.Violations, "SLOW: SELECT * usage (fetch only needed columns)")
	}

	// Guardrail 3: SELECT without LIMIT (Scale Risk)
	if strings.Contains(upperQuery, "SELECT") && !strings.Contains(upperQuery, "LIMIT") {
		event.Violations = append(event.Violations, "UNBOUNDED: SELECT without LIMIT (potential large result set)")
	}

	return &AnalysisResult{Event: event}, nil
}
