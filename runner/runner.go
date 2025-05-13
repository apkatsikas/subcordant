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

func (sr *SubcordantRunner) Init() {
	sr.SubsonicClient = &subsonic.SubsonicClient{}
	sr.DiscordClient = &discord.DiscordClient{}
	err := sr.SubsonicClient.Init()

	if err != nil {
		log.Fatalln(err)
	}

	err = sr.DiscordClient.Init(sr)
	if err != nil {
		log.Fatalln(err)
	}
}

func (sr *SubcordantRunner) HandlePlay(albumId string) {
	fmt.Printf("Called with %v - IMPLEMENT ME!", albumId)
}
