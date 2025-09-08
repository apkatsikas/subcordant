package subsonic

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	gosubsonic "github.com/apkatsikas/go-subsonic"
	"github.com/apkatsikas/subcordant/types"
)

type SubsonicClient struct {
	client *gosubsonic.Client
}

type TracksResult struct {
	Name   string
	Tracks []*gosubsonic.Child
}

func (sc *SubsonicClient) Init() error {
	subsonicUrl := os.Getenv("SUBSONIC_URL")
	subsonicUser := os.Getenv("SUBSONIC_USER")
	subsonicPassword := os.Getenv("SUBSONIC_PASSWORD")

	if subsonicUrl == "" {
		return fmt.Errorf("SUBSONIC_URL must be set")
	}
	if subsonicUser == "" {
		return fmt.Errorf("SUBSONIC_USER must be set")
	}
	if subsonicPassword == "" {
		return fmt.Errorf("SUBSONIC_PASSWORD must be set")
	}

	sc.client = &gosubsonic.Client{}
	sc.client.Client = &http.Client{}

	sc.client.BaseUrl = subsonicUrl
	sc.client.User = subsonicUser
	sc.client.PasswordAuth = true
	sc.client.ClientName = "subcordant bot"

	authErr := sc.client.Authenticate(subsonicPassword)
	if authErr != nil {
		return authErr
	}
	return nil
}

func (sc *SubsonicClient) GetTracks(id string) (*TracksResult, error) {
	tracks := &TracksResult{Tracks: []*gosubsonic.Child{}}

	if album, err := sc.client.GetAlbum(id); err == nil && album != nil {
		tracks.Tracks = append(tracks.Tracks, album.Song...)
		tracks.Name = fmt.Sprintf("album - %v", album.Name)
		return tracks, nil
	}

	if playlist, err := sc.client.GetPlaylist(id); err == nil && playlist != nil {
		tracks.Tracks = append(tracks.Tracks, playlist.Entry...)
		tracks.Name = fmt.Sprintf("playlist - %v", playlist.Name)
		return tracks, nil
	}

	if track, err := sc.client.GetSong(id); err == nil && track != nil {
		tracks.Tracks = append(tracks.Tracks, track)
		tracks.Name = fmt.Sprintf("track - %v", track.Title)
		return tracks, nil
	}

	return nil, fmt.Errorf("could not find an album, playlist or track with id of %v", id)
}

func (sc *SubsonicClient) GetTrackFromAlbum(albumId string, trackNumber int) (*gosubsonic.Child, error) {
	if trackNumber <= 0 {
		return nil, fmt.Errorf("track number must be greater than 0")
	}

	album, err := sc.client.GetAlbum(albumId)
	if err != nil || album == nil {
		return nil, fmt.Errorf("could not find album with ID of %v", albumId)
	}

	if int(trackNumber) > album.SongCount {
		return nil, fmt.Errorf(
			"track number %v was greater than album song count of %v", trackNumber, album.SongCount)
	}

	return album.Song[trackNumber], nil
}

func (sc *SubsonicClient) StreamUrl(trackId string) (*url.URL, error) {
	streamUrl, err := sc.client.StreamUrl(trackId, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to stream with track ID of %v", trackId)
	}
	return streamUrl, nil
}

func ToTrack(c *gosubsonic.Child) types.Track {
	return types.Track{
		ID:   c.ID,
		Path: c.Path,
	}
}
