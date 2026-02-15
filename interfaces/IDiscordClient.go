package interfaces

import (
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

type IDiscordClient interface {
	Init(commandHandler ICommandHandler) error
	JoinVoiceChat(channelId snowflake.ID) (snowflake.ID, error)
	SwitchVoiceChannel(channelId snowflake.ID) error
	SendMessage(message string)
	LeaveVoiceSession()
	Shutdown()
	SetFrameProvider(frameProvider voice.OpusFrameProvider)
}
