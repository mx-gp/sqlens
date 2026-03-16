package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestRecordEvent(t *testing.T) {
	// We can't easily reset promauto-registered metrics in a simple test,
	// but we can check if they increment.
	
	initialQueries := getCounterValue(QueriesTotal.WithLabelValues("false", "false"))
	
	RecordEvent(false, nil, 0.1, "test-fingerprint")
	
	finalQueries := getCounterValue(QueriesTotal.WithLabelValues("false", "false"))
	
	if finalQueries != initialQueries+1 {
		t.Errorf("Expected QueriesTotal to increment, got %f -> %f", initialQueries, finalQueries)
	}
}

func getCounterValue(counter prometheus.Counter) float64 {
	var m dto.Metric
	counter.Write(&m)
	return m.GetCounter().GetValue()
}
