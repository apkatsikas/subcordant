package interfaces

import (
	"context"
	"io"
)

type IStreamer interface {
	PrepStream(inputUrl string) error
	Stream(ctx context.Context, voice io.Writer) error
}
