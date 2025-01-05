package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"connectrpc.com/connect"

	uploadv1 "github.com/tkh/go-grpc-examples/client-stream/gen/upload/v1"
)

// Configuration constants.
const (
	defaultUploadsDir  = "uploads"        // Default directory for storing uploaded files.
	artificialDelay    = time.Millisecond // Simulated network latency for demonstration.
	defaultPermissions = 0o755            // Default Unix permissions for created directories.
	defaultContextTime = 30 * time.Second // Default timeout for file operations.
)

// Configuration flags.
var (
	uploadsDir = flag.String("uploads-dir", defaultUploadsDir, "Directory for storing uploaded files")
)

// Domain-specific errors.
var (
	ErrMissingMetadata = errors.New("first message must contain metadata")
	ErrSizeMismatch    = errors.New("final file size mismatch")
)

// UploadServer handles file upload requests.
type UploadServer struct{}

// NewUploadServer creates a new upload server instance.
func NewUploadServer() *UploadServer {
	return &UploadServer{}
}

// ensureUploadsDir creates the uploads directory if it doesn't exist.
func ensureUploadsDir() error {
	if _, err := os.Stat(*uploadsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(*uploadsDir, defaultPermissions); err != nil {
			return fmt.Errorf("creating uploads directory: %w", err)
		}
		slog.Info("uploads directory created", "path", *uploadsDir)
	}
	return nil
}

// UploadFile handles the file upload stream from clients.
func (s *UploadServer) UploadFile(
	ctx context.Context,
	stream *connect.ClientStream[uploadv1.UploadFileRequest],
) (*connect.Response[uploadv1.UploadFileResponse], error) {
	slog.Info("starting file upload", "header", stream.RequestHeader())

	if err := ensureUploadsDir(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("ensuring uploads directory: %w", err))
	}

	var (
		metadata   *uploadv1.Metadata
		file       *os.File
		totalBytes int64
		fileID     string
	)

	for stream.Receive() {
		msg := stream.Msg()
		if metadata == nil {
			meta := msg.GetMetadata()
			if meta == nil {
				return nil, connect.NewError(connect.CodeInvalidArgument, ErrMissingMetadata)
			}

			metadata = meta
			var err error
			if fileID, err = newFileID(); err != nil {
				return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("generating file ID: %w", err))
			}

			if file, err = createFile(fileID); err != nil {
				return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("creating file: %w", err))
			}
			defer file.Close()

			continue
		}

		chunk := msg.GetChunk()
		if len(chunk) == 0 {
			continue
		}

		time.Sleep(artificialDelay) // Simulate network latency.

		// Create a context with timeout for file operations.
		writeCtx, cancel := context.WithTimeout(ctx, defaultContextTime)
		defer cancel()

		// Use a channel to handle the write operation with timeout.
		errCh := make(chan error, 1)
		go func() {
			n, err := file.Write(chunk)
			if err != nil {
				errCh <- err
				return
			}
			totalBytes += int64(n)
			errCh <- nil
		}()

		// Wait for write completion or timeout.
		select {
		case err := <-errCh:
			if err != nil {
				cleanupFile(ctx, fileID)
				return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("writing chunk: %w", err))
			}
		case <-writeCtx.Done():
			cleanupFile(ctx, fileID)
			return nil, connect.NewError(connect.CodeDeadlineExceeded, fmt.Errorf("write operation timed out"))
		}
	}

	if err := stream.Err(); err != nil {
		cleanupFile(ctx, fileID)

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("processing stream: %w", err))
	}

	if totalBytes != metadata.GetSize() {
		cleanupFile(ctx, fileID)

		return nil, connect.NewError(
			connect.CodeDataLoss,
			fmt.Errorf("%w: got %d, want %d", ErrSizeMismatch, totalBytes, metadata.GetSize()),
		)
	}

	return connect.NewResponse(&uploadv1.UploadFileResponse{
		Id:       fileID,
		Filename: metadata.Filename,
		Size:     totalBytes,
	}), nil
}

// newFileID generates a string identifier for a file.
func newFileID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// createFile creates a new file in the uploads directory.
func createFile(fileID string) (*os.File, error) {
	return os.Create(filepath.Join(*uploadsDir, fileID))
}

// cleanupFile removes a file from the uploads directory.
func cleanupFile(ctx context.Context, fileID string) {
	err := os.Remove(filepath.Join(*uploadsDir, fileID))
	if err != nil {
		slog.ErrorContext(ctx, "file cleanup failed", "fileID", fileID)
	}
}
