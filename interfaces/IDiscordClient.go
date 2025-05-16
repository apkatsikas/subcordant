package interfaces

import (
	"context"
	"io"
)

type IDiscordClient interface {
	Init(commandHandler ICommandHandler) error
	JoinVoiceChat(cancelFunc context.CancelFunc) (io.Writer, error)
}
