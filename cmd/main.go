package main

import (
	"log"

	"github.com/apkatsikas/subcordant/discord"
	"github.com/apkatsikas/subcordant/runner"
	"github.com/apkatsikas/subcordant/streamer"
	"github.com/apkatsikas/subcordant/subsonic"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	runner := runner.SubcordantRunner{}
	err = runner.Init(&subsonic.SubsonicClient{}, &discord.DiscordClient{}, &streamer.Streamer{})
	if err != nil {
		log.Fatalf("failed to init runner: %v", err)
	}
}
