package flagutil

import (
	"flag"
	"fmt"
	"sync"
)

type StreamSource string

func (s *StreamSource) String() string {
	return string(*s)
}

func (s *StreamSource) Set(value string) error {
	switch value {
	case string(StreamSourceFile), string(StreamSourceStream):
		*s = StreamSource(value)
		return nil
	default:
		return fmt.Errorf("invalid value for streamFrom: %s", value)
	}
}

const (
	StreamSourceFile   StreamSource = "file"
	StreamSourceStream StreamSource = "stream"
)

type FlagUtil struct {
	StreamFrom StreamSource
}

func (fu *FlagUtil) Setup() {
	flag.Var(&fu.StreamFrom, "streamFrom",
		"Source from which to stream - valid values are 'file' or 'stream'")
	flag.Parse()
}

var (
	fu     *FlagUtil
	fuOnce sync.Once
)

func Get() *FlagUtil {
	if fu == nil {
		fuOnce.Do(func() {
			fu = &FlagUtil{}
		})
	}
	return fu
}
