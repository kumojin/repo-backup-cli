package storage

import (
	"context"
	"io"
)

type BlobRepository interface {
	Upload(ctx context.Context, blobName string, in io.Reader) (string, error)
}
