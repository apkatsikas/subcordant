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
	subsonicClient.GetAlbum("30c441134bfb1fa69022abe35af07a7c")

	err = discordClient.Init()
	if err != nil {
		log.Fatalln(err)
	}
}
