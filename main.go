package main

import (
	"log"

	"github.com/apkatsikas/subcordant/discord"
	"github.com/apkatsikas/subcordant/subsonic"
)

func main() {
	subsonicClient := subsonic.SubsonicClient{}
	discordClient := discord.DiscordClient{}
	err := subsonicClient.Init()
	if err != nil {
		log.Fatalln(err)
	}
	subsonicClient.ArtistSearch("my bloody valentine")

	err = discordClient.Init()
	if err != nil {
		log.Fatalln(err)
	}
}
