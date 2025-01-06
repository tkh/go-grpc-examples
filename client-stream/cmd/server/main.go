package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/tkh/go-grpc-examples/client-stream/gen/upload/v1/uploadv1connect"
	"github.com/tkh/go-grpc-examples/client-stream/internal/server"
)

const defaultShutdownTimeout = 5 * time.Second

type config struct {
	uploadsDir      string
	addr            string
	shutdownTimeout time.Duration
}

func parseFlags() *config {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A server that receives file uploads via gRPC streaming.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "server -uploads-dir /tmp/uploads -addr localhost:9090\n")
	}

	cfg := &config{
		shutdownTimeout: defaultShutdownTimeout,
	}
	fs.StringVar(&cfg.uploadsDir, "uploads-dir", "uploads", "Directory for storing uploaded files")
	fs.StringVar(&cfg.addr, "addr", "localhost:8080", "Address to listen on")
	fs.DurationVar(&cfg.shutdownTimeout, "shutdown-timeout", defaultShutdownTimeout, "Server shutdown timeout")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fs.Usage()
		os.Exit(1)
	}

	// Clean and validate uploads directory path
	cfg.uploadsDir = filepath.Clean(cfg.uploadsDir)
	if filepath.IsAbs(cfg.uploadsDir) && !strings.HasPrefix(cfg.uploadsDir, os.TempDir()) {
		log.Printf("Warning: Using absolute path outside of temp directory: %s", cfg.uploadsDir)
	}

	return cfg
}

func main() {
	cfg := parseFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := ensureUploadDir(ctx, cfg.uploadsDir); err != nil {
		log.Fatalf("failed to ensure upload directory: %v", err)
	}

	srv := newServer(cfg.addr, cfg.uploadsDir)
	if err := runServer(ctx, srv, cfg.shutdownTimeout); err != nil {
		slog.ErrorContext(ctx, "run server error", "error", err)
		os.Exit(1)
	}
}

func newServer(addr string, uploadsDir string) *http.Server {
	mux := http.NewServeMux()
	path, handler := uploadv1connect.NewUploadServiceHandler(server.NewUploadServer(uploadsDir))
	mux.Handle(path, handler)

	return &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
}

func runServer(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) error {
	errChan := make(chan error, 1)

	go func() {
		slog.Info("starting server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		slog.InfoContext(ctx, "server shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}

		return nil

	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	}
}

func ensureUploadDir(ctx context.Context, uploadDir string) error {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, 0o755)
		if err != nil {
			return fmt.Errorf("failed to create upload directory: %w", err)
		}
		slog.InfoContext(ctx, "upload directory created", "path", uploadDir)
	}

	testFile := uploadDir + "/.test"
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		return fmt.Errorf("upload directory is not writable: %w", err)
	}
	if err := os.Remove(testFile); err != nil {
		return fmt.Errorf("failed to remove test file: %w", err)
	}

	return nil
}
