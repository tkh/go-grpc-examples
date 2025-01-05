package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	uploadv1 "github.com/tkh/go-grpc-examples/client-stream/gen/upload/v1"
	"github.com/tkh/go-grpc-examples/client-stream/gen/upload/v1/uploadv1connect"
)

type uploadHandler struct{}

func (s uploadHandler) UploadFile(
	ctx context.Context,
	stream *connect.ClientStream[uploadv1.UploadFileRequest],
) (
	*connect.Response[uploadv1.UploadFileResponse],
	error,
) {
	slog.Info("request", "header", stream.RequestHeader())

	var (
		metadata   *uploadv1.Metadata
		file       *os.File
		totalBytes int64
		fileID     string
	)

	for stream.Receive() {
		msg := stream.Msg()
		// Get file metadata from first message.
		if metadata == nil {
			meta := msg.GetMetadata()
			if meta == nil {
				return nil, connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf("first message must contain metadata"),
				)
			}

			metadata = meta
			fileID, err := newFileID()
			if err != nil {
				return nil, connect.NewError(
					connect.CodeInternal,
					fmt.Errorf("failed to generate new file ID: %w", err),
				)
			}

			file, err = createFile(fileID)
			if err != nil {
				return nil, connect.NewError(
					connect.CodeInternal,
					fmt.Errorf("failed to create file: %w", err),
				)
			}
			defer file.Close()

			continue
		}

		chunk := msg.GetChunk()
		if len(chunk) == 0 {
			continue
		}

		n, err := file.Write(chunk)
		if err != nil {
			cleanupFile(ctx, fileID)

			return nil, connect.NewError(
				connect.CodeInternal,
				fmt.Errorf("failed to write chunk: %w", err),
			)
		}
		totalBytes += int64(n)
	}

	if err := stream.Err(); err != nil {
		cleanupFile(ctx, fileID)

		return nil, connect.NewError(
			connect.CodeInternal,
			fmt.Errorf("error while processing stream: %w", err),
		)
	}

	if totalBytes != metadata.GetSize() {
		cleanupFile(ctx, fileID)

		return nil, connect.NewError(
			connect.CodeDataLoss,
			fmt.Errorf(
				"final file size mismatch: got: %d, want: %d",
				totalBytes, metadata.GetSize(),
			),
		)
	}

	return connect.NewResponse(&uploadv1.UploadFileResponse{
		Id:       fileID,
		Filename: metadata.Filename,
		Size:     totalBytes,
	}), nil
}

func main() {
	mux := http.NewServeMux()
	path, handler := uploadv1connect.NewUploadServiceHandler(&uploadHandler{})
	mux.Handle(path, handler)
	log.Fatal(http.ListenAndServe(
		"localhost:8080",
		// h2c for HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	))
}

// newFileID generates a string identrifier for a file.
func newFileID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func createFile(fileID string) (*os.File, error) {
	if err := os.MkdirAll("uploads", 0o7500); err != nil {
		return nil, fmt.Errorf("failed to create uploads directory: %w", err)
	}

	return os.Create(filepath.Join("uploads", fileID))
}

func cleanupFile(ctx context.Context, fileID string) {
	err := os.Remove(filepath.Join("uploads", fileID))
	if err != nil {
		slog.ErrorContext(ctx, "file cleanup failed", "fileID", fileID)
	}
}
