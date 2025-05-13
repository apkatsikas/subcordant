package runner

import (
	"fmt"
	"log"

	"github.com/apkatsikas/subcordant/interfaces"
)

type SubcordantRunner struct {
	subsonicClient interfaces.ISubsonicClient
	discordClient  interfaces.IDiscordClient
}

func (sr *SubcordantRunner) Init(subsonicClient interfaces.ISubsonicClient, discordClient interfaces.IDiscordClient) {
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
	fmt.Printf("Called with %v - IMPLEMENT ME!", albumId)
}
