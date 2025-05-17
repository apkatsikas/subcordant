package interfaces

import (
	"context"
	"io"
)

type IExecCommander interface {
	Start(ctx context.Context, input io.ReadCloser, inputDestination string, cancelFunc context.CancelFunc) error
	Stream(voice io.Writer, cancelFunc context.CancelFunc) error
}
