package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"connectrpc.com/connect"

	uploadv1 "github.com/tkh/go-grpc-examples/client-stream/gen/upload/v1"
	"github.com/tkh/go-grpc-examples/client-stream/gen/upload/v1/uploadv1connect"
)

// Common errors.
var (
	ErrInvalidChunkSize = errors.New("chunk size must be positive")
	ErrInvalidURL       = errors.New("invalid server URL")
)

// UploadClient handles file uploads to the server.
type UploadClient struct {
	client       uploadv1connect.UploadServiceClient
	chunkSize    int
	progressFreq time.Duration
}

// NewUploadClient creates a new upload client instance.
func NewUploadClient(httpClient *http.Client, baseURL string, chunkSize int, progressFreq time.Duration) (*UploadClient, error) {
	if chunkSize <= 0 {
		return nil, fmt.Errorf("%w: got %d", ErrInvalidChunkSize, chunkSize)
	}

	if _, err := url.Parse(baseURL); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &UploadClient{
		client: uploadv1connect.NewUploadServiceClient(
			httpClient,
			baseURL,
		),
		chunkSize:    chunkSize,
		progressFreq: progressFreq,
	}, nil
}

// UploadFile streams a file to the server in chunks.
func (c *UploadClient) UploadFile(ctx context.Context, filePath string) error {
	stream := c.client.UploadFile(ctx)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("getting file info: %w", err)
	}

	if err := c.sendMetadata(stream, fileInfo); err != nil {
		return fmt.Errorf("sending metadata: %w", err)
	}

	if err := c.streamFileContent(stream, file, fileInfo.Size()); err != nil {
		return fmt.Errorf("streaming file content: %w", err)
	}

	response, err := stream.CloseAndReceive()
	if err != nil {
		return fmt.Errorf("completing upload: %w", err)
	}

	fmt.Printf("\nUpload complete: %s\n", response.Msg.Id)
	return nil
}

// sendMetadata sends the initial metadata message containing file information.
func (c *UploadClient) sendMetadata(
	stream *connect.ClientStreamForClient[uploadv1.UploadFileRequest, uploadv1.UploadFileResponse],
	fileInfo os.FileInfo,
) error {
	return stream.Send(&uploadv1.UploadFileRequest{
		Data: &uploadv1.UploadFileRequest_Metadata{
			Metadata: &uploadv1.Metadata{
				Filename: fileInfo.Name(),
				Size:     fileInfo.Size(),
			},
		},
	})
}

// streamFileContent reads the file in chunks and streams them to the server.
func (c *UploadClient) streamFileContent(
	stream *connect.ClientStreamForClient[uploadv1.UploadFileRequest, uploadv1.UploadFileResponse],
	file *os.File,
	totalSize int64,
) error {
	buffer := make([]byte, c.chunkSize)
	totalSent := int64(0)
	lastLog := time.Now()

	for {
		n, err := file.Read(buffer)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}

		if err := stream.Send(&uploadv1.UploadFileRequest{
			Data: &uploadv1.UploadFileRequest_Chunk{
				Chunk: buffer[:n],
			},
		}); err != nil {
			return fmt.Errorf("sending chunk: %w", err)
		}

		totalSent += int64(n)
		lastLog = logProgress(totalSent, totalSize, lastLog, c.progressFreq)
	}

	return nil
}

// logProgress updates the progress display if enough time has passed since the last update.
func logProgress(sent, total int64, lastLog time.Time, progressFreq time.Duration) time.Time {
	if time.Since(lastLog) < progressFreq {
		return lastLog
	}

	percent := float64(sent) / float64(total) * 100
	fmt.Printf("\rUploading: %.1f%% (%d/%d bytes)", percent, sent, total)

	if sent == total {
		fmt.Println() // Add newline when complete.
	}

	return time.Now()
}
