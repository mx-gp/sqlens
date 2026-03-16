package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	QueriesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sqlens_queries_total",
		Help: "Total number of SQL queries intercepted",
	}, []string{"n1_flag", "has_violations"})

	QueryLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "sqlens_query_latency_seconds",
		Help:    "Latency of SQL queries in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"fingerprint"})

	N1IncidentsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sqlens_n1_incidents_total",
		Help: "Total number of N+1 incidents detected",
	}, []string{"fingerprint"})

	ViolationsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sqlens_violations_total",
		Help: "Total number of SQL guardrail violations detected",
	}, []string{"violation_type"})
)

func RecordEvent(n1 bool, violations []string, latency float64, fingerprint string) {
	hasViolations := "false"
	if len(violations) > 0 {
		hasViolations = "true"
	}

	n1Flag := "false"
	if n1 {
		n1Flag = "true"
		N1IncidentsTotal.WithLabelValues(fingerprint).Inc()
	}

	QueriesTotal.WithLabelValues(n1Flag, hasViolations).Inc()
	QueryLatency.WithLabelValues(fingerprint).Observe(latency)

	for _, v := range violations {
		// Basic violation type extraction (e.g., from "SLOW: ...")
		vType := "unknown"
		if i := len(v); i > 0 {
			if parts := splitViolation(v); len(parts) > 0 {
				vType = parts[0]
			}
		}
		ViolationsTotal.WithLabelValues(vType).Inc()
	}
}

func splitViolation(v string) []string {
	for i := 0; i < len(v); i++ {
		if v[i] == ':' {
			return []string{v[:i]}
		}
	}
	return nil
}
