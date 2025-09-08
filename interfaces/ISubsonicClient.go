package interfaces

import (
	"net/url"

	gosubonic "github.com/apkatsikas/go-subsonic"
	"github.com/apkatsikas/subcordant/subsonic"
)

type ISubsonicClient interface {
	Init() error
	GetTracks(id string) (*subsonic.TracksResult, error)
	StreamUrl(trackId string) (*url.URL, error)
	GetTrackFromAlbum(albumId string, trackNumber int) (*gosubonic.Child, error)
}
