package playlist

import sub "github.com/apkatsikas/go-subsonic"

type PlaylistService struct {
	playlist []*sub.Child
}

func (ps *PlaylistService) Add(trackId *sub.Child) {
	ps.playlist = append(ps.playlist, trackId)
}

func (ps *PlaylistService) GetPlaylist() []*sub.Child {
	return ps.playlist
}

func (ps *PlaylistService) FinishTrack() {
	if len(ps.playlist) > 0 {
		ps.playlist = append(ps.playlist[:0], ps.playlist[0+1:]...)
	}
}

func (ps *PlaylistService) Clear() {
	ps.playlist = []*sub.Child{}
}
