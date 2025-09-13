package interfaces

import (
	"github.com/apkatsikas/subcordant/types"
	"github.com/diamondburned/arikawa/v3/discord"
)

type ICommandHandler interface {
	Play(subsonicId string, guildId discord.GuildID, channelId discord.ChannelID) (types.PlaybackState, error)
	PlayTrackFromAlbum(subsonicId string, trackNumber int, guildId discord.GuildID,
		switchToChannel discord.ChannelID) (types.PlaybackState, error)
	PlayTrackByName(
		query string, guildId discord.GuildID, switchToChannel discord.ChannelID) (types.PlaybackState, error)
	Reset()
	Disconnect()
	Skip()
}
