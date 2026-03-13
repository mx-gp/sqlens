package web

import (
	"encoding/json"
	"net/http"

	"github.com/sqlens/sqlens/store"
)

type Server struct {
	addr  string
	store *store.MemoryStore
}

func NewServer(addr string, s *store.MemoryStore) *Server {
	return &Server{
		addr:  addr,
		store: s,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/api/queries", s.handleQueries)
	mux.HandleFunc("/api/n1", s.handleN1)

	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>SQLens Dashboard</title>
    <style>
        body { font-family: monospace; background: #1e1e1e; color: #d4d4d4; padding: 20px; }
        .query { background: #2d2d2d; margin-bottom: 10px; padding: 10px; border-radius: 4px; }
        .n1 { border-left: 4px solid #f44336; }
        .latency { color: #4caf50; }
    </style>
    <script>
        async function fetchQueries() {
            const res = await fetch('/api/queries');
            const queries = await res.json();
            const container = document.getElementById('queries');
            container.innerHTML = queries.map(q => 
                '<div class="query ' + (q.N1Flag ? 'n1' : '') + '">' +
                '<strong>[' + new Date(q.Timestamp).toLocaleTimeString() + ']</strong> ' +
                '<span class="latency">' + (q.Latency / 1000000) + 'ms</span><br>' +
                '<code>' + q.RawQuery + '</code>' +
                (q.N1Flag ? '<br><strong style="color:#f44336">N+1 Detected!</strong>' : '') +
                '</div>'
            ).join('');
        }
        setInterval(fetchQueries, 1000);
        window.onload = fetchQueries;
    </script>
</head>
<body>
    <h1>SQLens Dashboard</h1>
    <h2>Recent Queries</h2>
    <div id="queries">Loading...</div>
</body>
</html>
`))
}

func (s *Server) handleQueries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	events := s.store.GetRecent(50)
	json.NewEncoder(w).Encode(events)
}

func (s *Server) handleN1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Get N1 incidents from the last 5 minutes
	events := s.store.GetN1Incidents(5 * 60 * 1000000000)
	json.NewEncoder(w).Encode(events)
}
