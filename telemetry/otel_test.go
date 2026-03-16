package telemetry

import (
	"context"
	"testing"
)

func TestInitTracer(t *testing.T) {
	tp, err := InitTracer()
	if err != nil {
		t.Fatalf("Failed to initialize tracer: %v", err)
	}
	if tp == nil {
		t.Fatal("Expected TracerProvider to be non-nil")
	}
	defer Shutdown(context.Background(), tp)
}
