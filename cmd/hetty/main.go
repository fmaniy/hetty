package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// defaultAddr is the default address the proxy and API server listens on.
	defaultAddr = ":8080"
	// defaultCertsDir is the default directory for storing CA certificates.
	defaultCertsDir = "~/.hetty/certs"
	// defaultDBPath is the default path for the bbolt database file.
	defaultDBPath = "~/.hetty/db"
)

// version is set at build time via ldflags.
var version = "dev"

func main() {
	if err := run(); err != nil {
		log.Fatalf("[FATAL] %v", err)
	}
}

func run() error {
	// Parse flags.
	addr := flag.String("addr", defaultAddr, "Address to listen on (e.g. \"127.0.0.1:8080\")")
	certsDir := flag.String("certs", defaultCertsDir, "Directory for storing CA certificates")
	dbPath := flag.String("db", defaultDBPath, "Path for the bbolt database file")
	upstreamProxy := flag.String("upstream-proxy", "", "Optional upstream proxy URL (e.g. \"http://localhost:8888\")")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	printVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *printVersion {
		fmt.Printf("hetty %v\n", version)
		return nil
	}

	if *verbose {
		log.Printf("[DEBUG] Starting hetty %v", version)
		log.Printf("[DEBUG] Listening on %v", *addr)
		log.Printf("[DEBUG] Certificates directory: %v", *certsDir)
		log.Printf("[DEBUG] Database path: %v", *dbPath)
		if *upstreamProxy != "" {
			log.Printf("[DEBUG] Upstream proxy: %v", *upstreamProxy)
		}
	}

	// Set up a base context that is cancelled on OS interrupt signals.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create a TCP listener.
	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %v: %w", *addr, err)
	}

	// Configure TLS (placeholder — real CA/cert loading will be added later).
	_ = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Set up the HTTP server (mux and handlers will be wired in subsequent files).
	srv := &http.Server{
		Handler:      http.DefaultServeMux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Serve in a goroutine so we can listen for shutdown signals.
	serveErr := make(chan error, 1)
	go func() {
		log.Printf("[INFO] hetty %v is running on %v", version, listener.Addr())
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- fmt.Errorf("http server error: %w", err)
		}
		close(serveErr)
	}()

	// Block until a signal or server error is received.
	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		log.Println("[INFO] Shutting down gracefully…")
	}

	// Allow up to 10 seconds for in-flight requests to complete.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Println("[INFO] Shutdown complete.")
	return nil
}
