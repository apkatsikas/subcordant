package interfaces

import (
	"io"

	sub "github.com/delucks/go-subsonic"
)

type ISubsonicClient interface {
	Init() error
	GetAlbum(albumId string) (*sub.AlbumID3, error)
	Stream(trackId string) (io.ReadCloser, error)
}
