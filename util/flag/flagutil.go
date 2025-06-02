package flagutil

import (
	"flag"
	"fmt"
	"sync"
)

type StreamFrom string

func (s *StreamFrom) String() string {
	return string(*s)
}

func (s *StreamFrom) Set(value string) error {
	switch value {
	case string(StreamFromFile), string(StreamFromStream):
		*s = StreamFrom(value)
		return nil
	default:
		return fmt.Errorf("invalid value for streamFrom: %s", value)
	}
}

const (
	StreamFromFile   StreamFrom = "file"
	StreamFromStream StreamFrom = "stream"
)

type FlagUtil struct {
	StreamFrom StreamFrom
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
