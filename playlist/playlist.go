package playlist

type PlaylistService struct {
	playlist []string
}

func (ps *PlaylistService) Add(trackId string) {
	ps.playlist = append(ps.playlist, trackId)
}

func (ps *PlaylistService) GetPlaylist() []string {
	return ps.playlist
}

func (ps *PlaylistService) FinishTrack() {
	if len(ps.playlist) > 0 {
		ps.playlist = append(ps.playlist[:0], ps.playlist[0+1:]...)
	}
}

func (ps *PlaylistService) Clear() {
	ps.playlist = []string{}
}
