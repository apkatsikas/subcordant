package subsonic

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	sub "github.com/apkatsikas/go-subsonic"
	"github.com/apkatsikas/subcordant/types"
)

type SubsonicClient struct {
	client *sub.Client
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
	sc.client.ClientName = "coolhacker"

	authErr := sc.client.Authenticate(subsonicPassword)
	if authErr != nil {
		return authErr
	}
	return nil
}

func (sc *SubsonicClient) GetAlbum(albumId string) (*sub.AlbumID3, error) {
	albumResult, err := sc.client.GetAlbum(albumId)

	if err != nil {
		return nil, fmt.Errorf("failed to get album with ID %v - %v", albumId, err)
	}
	return albumResult, nil
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
