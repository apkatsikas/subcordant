package interfaces

import (
	"context"
	"io"
	"net/url"
)

type IStreamer interface {
	PrepStreamFromStream(streamUrl *url.URL) error
	PrepStreamFromFile(file string) error
	Stream(ctx context.Context, voice io.Writer) error
}
