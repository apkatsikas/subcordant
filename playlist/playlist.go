package playlist

import (
	"slices"
	"sync"

	"github.com/apkatsikas/subcordant/types"
)

type PlaylistService struct {
	mu       sync.Mutex
	playlist []types.Track
}

func (ps *PlaylistService) Add(track types.Track) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.playlist = append(ps.playlist, track)
}

func (ps *PlaylistService) GetPlaylist() []types.Track {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.playlist
}

func (ps *PlaylistService) FinishTrack() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if len(ps.playlist) > 0 {
		ps.playlist = slices.Delete(ps.playlist, 0, 1)
	}
}

func (ps *PlaylistService) Clear() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.playlist = []types.Track{}
}
