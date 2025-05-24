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
	godotenv.Load() // ignore errors if file does not exist

	runner := runner.SubcordantRunner{}
	err := runner.Init(&subsonic.SubsonicClient{}, &discord.DiscordClient{}, &streamer.Streamer{})
	if err != nil {
		log.Fatalf("failed to init runner: %v", err)
	}
}
