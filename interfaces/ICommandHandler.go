package interfaces

import (
	"github.com/apkatsikas/subcordant/types"
	"github.com/diamondburned/arikawa/v3/discord"
)

type ICommandHandler interface {
	Play(albumId string, guildId discord.GuildID, channelId discord.ChannelID) (types.PlaybackState, error)
	Reset()
	Disconnect()
	Skip()
}
