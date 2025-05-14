package main

import (
	"github.com/apkatsikas/subcordant/discord"
	"github.com/apkatsikas/subcordant/ffmpeg"
	"github.com/apkatsikas/subcordant/runner"
	"github.com/apkatsikas/subcordant/subsonic"
)

func main() {
	runner := runner.SubcordantRunner{}
	runner.Init(&subsonic.SubsonicClient{}, &discord.DiscordClient{}, &ffmpeg.FfmpegCommander{})
}
