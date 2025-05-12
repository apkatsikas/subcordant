package main

import (
	"github.com/apkatsikas/subcordant/subsonic"
)

func main() {
	subsonicClient := subsonic.SubsonicClient{}
	initErr := subsonicClient.Init()
	if initErr != nil {
		panic(initErr)
	}
	subsonicClient.ArtistSearch("my bloody valentine")
}
