package subsonic

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	sub "github.com/apkatsikas/go-subsonic"
	"github.com/apkatsikas/subcordant/types"
)

type SubsonicClient struct {
	client *sub.Client
}

type TracksResult struct {
	Name   string
	Tracks []*sub.Child
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

	sc.client = &sub.Client{}
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
	tracks := &TracksResult{}
	tracks.Tracks = []*sub.Child{}
	albumResult, err := sc.client.GetAlbum(id)

	if err != nil || albumResult != nil {
		playlistResult, err := sc.client.GetPlaylist(id)

		if err != nil {
			track, err := sc.client.GetSong(id)
			if err != nil {
				return nil, fmt.Errorf("could not find an album, playlist or track with id of %v", id)
			}
			tracks.Tracks = append(tracks.Tracks, track)
			tracks.Name = fmt.Sprintf("track - %v", track.Title)
			return tracks, nil
		}
		tracks.Tracks = append(tracks.Tracks, playlistResult.Entry...)
		tracks.Name = fmt.Sprintf("playlist - %v", playlistResult.Name)
		return tracks, nil
	}
	log.Printf("Got %v", albumResult)
	tracks.Tracks = append(tracks.Tracks, albumResult.Song...)
	tracks.Name = fmt.Sprintf("album - %v", albumResult.Name)
	return tracks, nil
}

func (sc *SubsonicClient) StreamUrl(trackId string) (*url.URL, error) {
	streamUrl, err := sc.client.StreamUrl(trackId, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to stream with track ID of %v", trackId)
	}
	return streamUrl, nil
}

func ToTrack(c *sub.Child) types.Track {
	return types.Track{
		ID:   c.ID,
		Path: c.Path,
	}
}
