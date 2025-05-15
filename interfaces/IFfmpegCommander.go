package interfaces

import (
	"context"
	"io"
)

type IFfmpegCommander interface {
	Start(ctx context.Context, input io.ReadCloser, inputDestination string) error
	Stream(voice io.Writer) error
}
