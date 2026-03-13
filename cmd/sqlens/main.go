package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sqlens/sqlens/analyzer"
	"github.com/sqlens/sqlens/config"
	"github.com/sqlens/sqlens/proxy"
	"github.com/sqlens/sqlens/store"
	"github.com/sqlens/sqlens/web"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize Storage
	memStore := store.NewMemoryStore()

	// Initialize Analyzer Pipeline
	pipeline := analyzer.NewPipeline(
		analyzer.NewFingerprintAnalyzer(),
		analyzer.NewN1DetectorAnalyzer(
			time.Duration(cfg.N1WindowSecs)*time.Second,
			cfg.N1Threshold,
		),
		analyzer.NewGuardrailAnalyzer(), // New Feature: Real-time SQL Guardrails
	)

	// Initialize TCP Proxy
	proxyServer := proxy.NewServer(cfg.ListenAddr, cfg.TargetAddr, pipeline, memStore, cfg.RedactSensitive)

	// Initialize Web Dashboard
	webServer := web.NewServer(":8080", memStore)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start servers
	go func() {
		if err := webServer.Start(); err != nil {
			slog.Error("Web server failed", "err", err)
		}
	}()

	go func() {
		if err := proxyServer.Start(ctx); err != nil {
			slog.Error("Proxy server failed", "err", err)
		}
	}()

	slog.Info("SQLens started", "web", ":8080", "proxy", cfg.ListenAddr, "target", cfg.TargetAddr)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down SQLens...")
	cancel()
	time.Sleep(time.Second) // wait for cleanup
}
