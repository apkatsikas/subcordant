package interfaces

import (
	"net/url"

	gosubsonic "github.com/apkatsikas/go-subsonic"
	"github.com/apkatsikas/subcordant/subsonic"
)

type ISubsonicClient interface {
	Init() error
	GetTracks(id string) (*subsonic.TracksResult, error)
	StreamUrl(trackId string) (*url.URL, error)
	GetTrackFromAlbum(albumId string, trackNumber int) (*gosubsonic.Child, error)
	GetTrackByName(query string) (*gosubsonic.Child, error)
	GetAlbumByName(query string) (*gosubsonic.AlbumID3, error)
}
