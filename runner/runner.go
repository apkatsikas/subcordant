package runner

import (
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

	err = sr.DiscordClient.Init()
	if err != nil {
		log.Fatalln(err)
	}
}
