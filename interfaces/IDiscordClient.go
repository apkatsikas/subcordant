package interfaces

import "io"

type IDiscordClient interface {
	Init(commandHandler ICommandHandler) error
	JoinVoiceChat() (io.Writer, error)
}
