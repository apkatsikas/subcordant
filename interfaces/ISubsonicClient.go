package interfaces

import (
	sub "github.com/apkatsikas/subcordant/go-subsonic"
)

type ISubsonicClient interface {
	Init() error
	GetAlbum(albumId string) (*sub.AlbumID3, error)
	Stream(trackId string) (string, error)
}
