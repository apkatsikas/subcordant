package interfaces

import (
	"context"
	"io"
)

type IFfmpegCommander interface {
	Start(ctx context.Context, input io.ReadCloser) error
	Stream(voice io.Writer) error
}
