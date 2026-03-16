package proxy

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/sqlens/sqlens/analyzer"
	"github.com/sqlens/sqlens/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	listenAddr      string
	targetAddr      string
	pipeline        *analyzer.Pipeline
	store           EventStore
	redactSensitive bool
	tracer          trace.Tracer
}

// EventStore represents an interface to save and retrieve query events
type EventStore interface {
	Save(event analyzer.QueryEvent)
}

func NewServer(listen, target string, p *analyzer.Pipeline, s EventStore, redact bool) *Server {
	return &Server{
		listenAddr:      listen,
		targetAddr:      target,
		pipeline:        p,
		store:           s,
		redactSensitive: redact,
		tracer:          otel.Tracer("sqlens-proxy"),
	}
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	slog.Info("Proxy listening", "addr", s.listenAddr, "target", s.targetAddr)

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			slog.Error("Failed to accept connection", "err", err)
			continue
		}

		go s.handleConnection(ctx, clientConn)
	}
}

func (s *Server) handleConnection(ctx context.Context, clientConn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from panic in connection handler", "panic", r)
		}
		clientConn.Close()
	}()

	// Set a deadline for the initial connection to avoid hanging on zombie clients
	clientConn.SetDeadline(time.Now().Add(5 * time.Second))
	targetConn, err := net.DialTimeout("tcp", s.targetAddr, 3*time.Second)
	if err != nil {
		slog.Error("Failed to connect to target", "err", err)
		return
	}
	defer targetConn.Close()
	clientConn.SetDeadline(time.Time{}) // Remove deadline for active session

	connID := clientConn.RemoteAddr().String()
	slog.Debug("Connection established", "id", connID)

	var lastQueryStart time.Time
	var lastQueryStr string
	var lastQueryCtx context.Context
	var lastQuerySpan trace.Span
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(2)

	// Client to Target (Requests)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic in request sniffer", "panic", r)
			}
			wg.Done()
		}()
		buf := make([]byte, 8192) // Larger buffer for complex queries
		for {
			n, err := clientConn.Read(buf)
			if err != nil {
				return
			}
			
			// Sniff PostgreSQL 'Q' (Simple Query)
			if n > 5 && buf[0] == 'Q' {
				mu.Lock()
				if lastQuerySpan != nil {
					lastQuerySpan.End()
				}
				lastQueryStart = time.Now()
				// Basic safety: limit query string size to avoid huge allocations
				end := n - 1
				if end > 2048 {
					end = 2048
				}
				lastQueryStr = string(buf[5:end])
				
				// Start OpenTelemetry span
				qCtx, span := s.tracer.Start(ctx, "sql_query", 
					trace.WithAttributes(attribute.String("sql.query", lastQueryStr)))
				lastQueryCtx = qCtx
				lastQuerySpan = span
				
				mu.Unlock()
			}

			_, err = targetConn.Write(buf[:n])
			if err != nil {
				return
			}
		}
	}()

	// Target to Client (Responses)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic in response sniffer", "panic", r)
			}
			wg.Done()
		}()
		buf := make([]byte, 8192)
		for {
			n, err := targetConn.Read(buf)
			if err != nil {
				return
			}

			mu.Lock()
			if !lastQueryStart.IsZero() {
				latency := time.Since(lastQueryStart)
				query := lastQueryStr
				qCtx := lastQueryCtx
				qSpan := lastQuerySpan
				lastQueryStart = time.Time{}
				lastQueryCtx = nil
				lastQuerySpan = nil
				mu.Unlock()

				go func(q string, l time.Duration, ctx context.Context, span trace.Span) {
					defer func() { recover() }() // Async processing safety
					event := analyzer.QueryEvent{
						ConnectionID: connID,
						RawQuery:     q,
						Timestamp:    time.Now(),
						Latency:      l,
					}
					
					// Reassign event to the one that's been through the pipeline
					event = s.pipeline.Process(ctx, event)

					// Record Prometheus metrics
					metrics.RecordEvent(event.N1Flag, event.Violations, l.Seconds(), event.Fingerprint)

					// Update and End Span
					if span != nil {
						span.SetAttributes(
							attribute.String("sql.fingerprint", event.Fingerprint),
							attribute.Bool("sqlens.n1_flag", event.N1Flag),
							attribute.StringSlice("sqlens.violations", event.Violations),
						)
						span.End()
					}
					
					// If redaction is enabled, mask the raw SQL
					if s.redactSensitive {
						event.RawQuery = event.Fingerprint
					}
					
					s.store.Save(event)
				}(query, latency, qCtx, qSpan)
			} else {
				mu.Unlock()
			}

			_, err = clientConn.Write(buf[:n])
			if err != nil {
				return
			}
		}
	}()

	wg.Wait()
	slog.Debug("Connection closed", "id", connID)
}
