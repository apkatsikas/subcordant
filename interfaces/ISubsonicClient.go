package interfaces

import (
	"net/url"

	"github.com/apkatsikas/subcordant/subsonic"
)

type ISubsonicClient interface {
	Init() error
	GetTracks(id string) (*subsonic.TracksResult, error)
	StreamUrl(trackId string) (*url.URL, error)
}
