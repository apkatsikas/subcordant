package playlist

import (
	"slices"

	"github.com/apkatsikas/subcordant/types"
)

type PlaylistService struct {
	playlist []types.Track
}

func (ps *PlaylistService) Add(track types.Track) {
	ps.playlist = append(ps.playlist, track)
}

func (ps *PlaylistService) GetPlaylist() []types.Track {
	return ps.playlist
}

func (ps *PlaylistService) FinishTrack() {
	if len(ps.playlist) > 0 {
		ps.playlist = slices.Delete(ps.playlist, 0, 1)
	}
}

func (ps *PlaylistService) Clear() {
	ps.playlist = []types.Track{}
}
