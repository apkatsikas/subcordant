package interfaces

import (
	"context"
	"io"
	"net/url"
)

type IStreamer interface {
	PrepStream(inputUrl *url.URL) error
	Stream(ctx context.Context, voice io.Writer) error
}
