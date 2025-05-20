package interfaces

import (
	"io"
	"net/url"
)

type IStreamer interface {
	PrepStream(inputUrl *url.URL) error
	Stream(voice io.Writer) error
}
