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
