package analyzer

import (
	"context"
	"testing"
)

func TestFingerprintAnalyzer(t *testing.T) {
	fa := NewFingerprintAnalyzer()
	ctx := context.Background()

	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{
			"Simple Select",
			"SELECT * FROM users WHERE id = 123",
			"SELECT * FROM users WHERE id = ?",
		},
		{
			"Select with String",
			"SELECT * FROM users WHERE email = 'test@example.com'",
			"SELECT * FROM users WHERE email = ?",
		},
		{
			"Multiple Params",
			"INSERT INTO users (name, age) VALUES ('Alice', 30)",
			"INSERT INTO users (name, age) VALUES (?, ?)",
		},
		{
			"Update query",
			"UPDATE users SET age = 31 WHERE id = 1",
			"UPDATE users SET age = ? WHERE id = ?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := QueryEvent{RawQuery: tt.query}
			result, err := fa.Analyze(ctx, event)
			if err != nil {
				t.Fatalf("Analyze failed: %v", err)
			}
			if result.Event.Fingerprint != tt.expected {
				t.Errorf("Expected fingerprint %q, got %q", tt.expected, result.Event.Fingerprint)
			}
		})
	}
}
