package interfaces

import "github.com/apkatsikas/subcordant/types"

type ICommandHandler interface {
	Play(albumId string) (types.PlaybackState, error)
	Reset()
}
