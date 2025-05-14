package main

import (
	"log"

	"github.com/apkatsikas/subcordant/discord"
	"github.com/apkatsikas/subcordant/ffmpeg"
	"github.com/apkatsikas/subcordant/runner"
	"github.com/apkatsikas/subcordant/subsonic"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	runner := runner.SubcordantRunner{}
	runner.Init(&subsonic.SubsonicClient{}, &discord.DiscordClient{}, &ffmpeg.FfmpegCommander{})
}
