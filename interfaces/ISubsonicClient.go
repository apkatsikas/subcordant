package interfaces

import (
	"net/url"

	sub "github.com/apkatsikas/go-subsonic"
)

type ISubsonicClient interface {
	Init() error
	GetAlbum(albumId string) (*sub.AlbumID3, error)
	StreamUrl(trackId string) (*url.URL, error)
}
