package interfaces

import (
	"context"
	"io"
	"net/url"
)

type IStreamer interface {
	PrepStream(ctx context.Context, inputUrl *url.URL,
		cancelFunc context.CancelFunc) error
	Stream(voice io.Writer, cancelFunc context.CancelFunc) error
}
