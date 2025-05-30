package interfaces

import (
	"io"

	"github.com/diamondburned/arikawa/v3/discord"
)

type IDiscordClient interface {
	Init(commandHandler ICommandHandler) error
	JoinVoiceChat(guildId discord.GuildID, channelId discord.ChannelID) (discord.ChannelID, error)
	SwitchVoiceChannel(channelId discord.ChannelID) error
	SendMessage(message string)
	GetVoice() io.Writer
}
