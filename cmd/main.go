package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/apkatsikas/subcordant/discord"
	"github.com/apkatsikas/subcordant/runner"
	"github.com/apkatsikas/subcordant/streamer"
	"github.com/apkatsikas/subcordant/subsonic"
	flagutil "github.com/apkatsikas/subcordant/util/flag"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load() // ignore errors if file does not exist

	fu := flagutil.Get()
	fu.Setup()

	runner := runner.SubcordantRunner{}
	err := runner.Init(&subsonic.SubsonicClient{}, &discord.DiscordClient{}, &streamer.Streamer{},
		fu.StreamFrom, fu.IdleDisconnectTimeout)
	if err != nil {
		log.Fatalf("failed to init runner: %v", err)
	}
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	runner.Shutdown()
	log.Println("ty for jammin w/ me <3")
}
