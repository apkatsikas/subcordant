package runner

import (
	"fmt"
	"log"

	"github.com/apkatsikas/subcordant/discord"
	"github.com/apkatsikas/subcordant/subsonic"
)

type SubcordantRunner struct {
	*subsonic.SubsonicClient
	*discord.DiscordClient
}

func (sr *SubcordantRunner) Run() {
	sr.SubsonicClient = &subsonic.SubsonicClient{}
	sr.DiscordClient = &discord.DiscordClient{}
	err := sr.SubsonicClient.Init()
	if err != nil {
		log.Fatalln(err)
	}
	album := sr.SubsonicClient.GetAlbum("30c441134bfb1fa69022abe35af07a7c")

	for _, song := range album.Song {
		fmt.Println(song.ID)
	}

	err = sr.DiscordClient.Init()
	if err != nil {
		log.Fatalln(err)
	}
}
