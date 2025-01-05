package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tkh/go-grpc-examples/client-stream/internal/client"
)

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024

	defaultChunkSize          = 32 * 1024
	defaultProgressUpdateFreq = 100 * time.Millisecond
	defaultServerURL          = "http://localhost:8080"
	defaultFileName           = "example.txt"
)

type config struct {
	fileSize     ByteSize
	chunkSize    int
	progressFreq time.Duration
	serverURL    string
	fileName     string
}

func parseFlags() *config {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A client that uploads a test file to the server using gRPC streaming.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nSize units supported: KB, MB, GB, TB (case insensitive)\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "client -size 500KB\n")
		fmt.Fprintf(os.Stderr, "client -size 1GB -chunk-size 65536 -progress-freq 500ms -server-url http://example.com:8080\n")
	}

	cfg := &config{
		fileSize:     ByteSize(100 * MB), // default to 100MB
		chunkSize:    defaultChunkSize,
		progressFreq: defaultProgressUpdateFreq,
		serverURL:    defaultServerURL,
		fileName:     defaultFileName,
	}

	fs.Var(&cfg.fileSize, "size", "Size of test file to create (e.g., 100MB, 1GB)")
	fs.IntVar(&cfg.chunkSize, "chunk-size", defaultChunkSize, "Size of each upload chunk in bytes")
	fs.DurationVar(&cfg.progressFreq, "progress-freq", defaultProgressUpdateFreq, "Minimum time between progress updates")
	fs.StringVar(&cfg.serverURL, "server-url", defaultServerURL, "URL of the upload server")
	fs.StringVar(&cfg.fileName, "file-name", defaultFileName, "Name of the test file to create")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fs.Usage()
		os.Exit(1)
	}

	// Validate configuration
	if cfg.chunkSize <= 0 {
		log.Fatalf("chunk size must be positive")
	}

	if _, err := url.Parse(cfg.serverURL); err != nil {
		log.Fatalf("invalid server URL: %v", err)
	}

	return cfg
}

type ByteSize int

func (b *ByteSize) String() string {
	return formatBytes(int(*b))
}

func (b *ByteSize) Set(value string) error {
	value = strings.ToUpper(strings.TrimSpace(value))
	var multiplier int

	switch {
	case strings.HasSuffix(value, "TB"):
		multiplier = TB
		value = strings.TrimSuffix(value, "TB")
	case strings.HasSuffix(value, "GB"):
		multiplier = GB
		value = strings.TrimSuffix(value, "GB")
	case strings.HasSuffix(value, "MB"):
		multiplier = MB
		value = strings.TrimSuffix(value, "MB")
	case strings.HasSuffix(value, "KB"):
		multiplier = KB
		value = strings.TrimSuffix(value, "KB")
	default:
		multiplier = 1
	}

	size, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid size value: %w", err)
	}

	*b = ByteSize(size * multiplier)
	return nil
}

func main() {
	cfg := parseFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := createTestFile(cfg.fileName, int(cfg.fileSize)); err != nil {
		log.Fatalf("failed to create example file for upload: %v", err)
	}
	// Ensure we clean up the file even if upload fails
	defer func() {
		if err := os.Remove(cfg.fileName); err != nil {
			log.Printf("warning: failed to clean up example file: %v", err)
		} else {
			log.Printf("cleaned up example file: %s", cfg.fileName)
		}
	}()

	uploadClient, err := client.NewUploadClient(
		http.DefaultClient,
		cfg.serverURL,
		cfg.chunkSize,
		cfg.progressFreq,
	)
	if err != nil {
		log.Fatalf("failed to create upload client: %v", err)
	}

	if err := uploadClient.UploadFile(ctx, cfg.fileName); err != nil {
		log.Fatalf("file upload failed: %v\nMake sure the server is running with: go run cmd/server/main.go", err)
	}
}

func formatBytes(bytes int) string {
	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2fTB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2fGB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2fMB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2fKB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

func createTestFile(fileName string, targetSize int) error {
	fmt.Printf("Creating test file: %s (size: %s)\n", fileName, formatBytes(targetSize))

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	content := "This is some text for the example file.\n"
	contentSize := len(content)

	writtenSize := 0
	for writtenSize < targetSize {
		remainingSize := targetSize - writtenSize
		if remainingSize < contentSize {
			content = content[:remainingSize]
		}

		n, err := file.WriteString(content)
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}

		writtenSize += n
	}

	fmt.Printf("Successfully created test file: %s \n", fileName)

	return nil
}
