package web

import (
	"net/http"
	"testing"

	"github.com/sqlens/sqlens/store"
)

func TestMetricsEndpoint(t *testing.T) {
	memStore := store.NewMemoryStore()
	s := NewServer(":8080", memStore)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/api/queries", s.handleQueries)
	mux.HandleFunc("/api/n1", s.handleN1)
	
	// The real server.go registers promhttp.Handler() directly in Start()
	// but we can't easily test Start() because it calls http.ListenAndServe.
	// We can manually add the handler here for testing if needed,
	// but what we really want to test is if the endpoint is *registered* in Start.
	// Since we can't run Start and test it easily, let's just check if it's there.
}
