package interfaces

import (
	"context"
	"net/url"

	"github.com/disgoorg/disgo/voice"
)

type IStreamer interface {
	PrepStreamFromStream(streamUrl *url.URL) error
	PrepStreamFromFile(file string) error
	Stream(ctx context.Context, setFrameProvider func(voice.OpusFrameProvider)) error
}
