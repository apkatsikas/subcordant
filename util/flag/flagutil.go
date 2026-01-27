package flagutil

import (
	"flag"
	"fmt"
	"strconv"
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

type IdleDisconnectTimeout int

func (idt *IdleDisconnectTimeout) String() string {
	return strconv.Itoa(int(*idt))
}

func (idt *IdleDisconnectTimeout) Set(value string) error {
	timeout, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid value for idleDisconnectTimeout: %s", value)
	}
	*idt = IdleDisconnectTimeout(timeout)
	return nil
}

const (
	StreamFromFile   StreamFrom = "file"
	StreamFromStream StreamFrom = "stream"
)

type FlagUtil struct {
	StreamFrom            StreamFrom
	IdleDisconnectTimeout IdleDisconnectTimeout
}

func (fu *FlagUtil) Setup() {
	flag.Var(&fu.StreamFrom, "streamFrom",
		"Source from which to stream - valid values are 'file' or 'stream'")
	flag.Var(&fu.IdleDisconnectTimeout, "idleDisconnectTimeout",
		"Duration in minutes after which bot will disconnect when no music is playing")
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
