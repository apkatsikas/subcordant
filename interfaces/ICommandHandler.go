package interfaces

import (
	"github.com/apkatsikas/subcordant/types"
	"github.com/disgoorg/snowflake/v2"
)

type ICommandHandler interface {
	Play(subsonicId string, channelId snowflake.ID) (types.PlaybackState, error)
	PlayTrackFromAlbum(subsonicId string, trackNumber int,
		switchToChannel snowflake.ID) (types.PlaybackState, error)
	PlayTrackByName(
		query string, switchToChannel snowflake.ID) (types.PlaybackState, error)
	PlayAlbumByName(query string, switchToChannel snowflake.ID) (types.PlaybackState, error)
	Reset()
	Disconnect()
	Skip()
}
