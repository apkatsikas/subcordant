package interfaces

import (
	"context"
	"io"
)

type IFfmpegCommander interface {
	Start(ctx context.Context, file string) error
	Stream(voice io.Writer) error
}
