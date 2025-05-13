package subsonic

import (
	"fmt"
	"log"
	"net/http"
	"os"

	sub "github.com/delucks/go-subsonic"
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

func (sc *SubsonicClient) GetAlbum(albumId string) *sub.AlbumID3 {
	albumResult, err := sc.client.GetAlbum(albumId)
	if err != nil {
		log.Fatalln(err)
	}

	return albumResult
}
