package runner

import (
	"log"

	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/apkatsikas/subcordant/playlist"
)

type SubcordantRunner struct {
	subsonicClient interfaces.ISubsonicClient
	discordClient  interfaces.IDiscordClient
	*playlist.PlaylistService
}

func (sr *SubcordantRunner) Init(subsonicClient interfaces.ISubsonicClient, discordClient interfaces.IDiscordClient) {
	sr.PlaylistService = &playlist.PlaylistService{}
	sr.subsonicClient = subsonicClient
	sr.discordClient = discordClient

	err := sr.subsonicClient.Init()
	if err != nil {
		log.Fatalln(err)
	}

	err = sr.discordClient.Init(sr)
	if err != nil {
		log.Fatalln(err)
	}
}

func (sr *SubcordantRunner) HandlePlay(albumId string) {
	album, err := sr.subsonicClient.GetAlbum(albumId)

	if err != nil {
		log.Fatalln(err)
	}

	for _, song := range album.Song {
		sr.PlaylistService.Add(song.ID)
	}

	err = sr.discordClient.JoinVoiceChat()
	if err != nil {
		log.Fatalln(err)
	}
}
