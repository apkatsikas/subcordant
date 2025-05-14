package runner

import (
	"context"
	"log"
	"os"

	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/apkatsikas/subcordant/playlist"
)

type SubcordantRunner struct {
	subsonicClient  interfaces.ISubsonicClient
	discordClient   interfaces.IDiscordClient
	ffmpegCommander interfaces.IFfmpegCommander
	*playlist.PlaylistService
	removeMeFileToStream string
}

func (sr *SubcordantRunner) Init(
	subsonicClient interfaces.ISubsonicClient, discordClient interfaces.IDiscordClient,
	ffmpegCommander interfaces.IFfmpegCommander) {
	sr.PlaylistService = &playlist.PlaylistService{}
	sr.subsonicClient = subsonicClient
	sr.discordClient = discordClient
	sr.ffmpegCommander = ffmpegCommander

	// TODO - delete
	ffmpegFile := os.Getenv("FFMPEG_FILE")
	if ffmpegFile == "" {
		log.Fatalln("FFMPEG_FILE must be set.")
	}
	sr.removeMeFileToStream = ffmpegFile

	err := sr.subsonicClient.Init()
	if err != nil {
		log.Fatalln(err)
	}

	err = sr.discordClient.Init(sr)
	if err != nil {
		log.Fatalln(err)
	}
}

func (sr *SubcordantRunner) HandlePlay(albumId string) error {
	album, err := sr.subsonicClient.GetAlbum(albumId)

	if err != nil {
		return err
	}

	for _, song := range album.Song {
		sr.PlaylistService.Add(song.ID)
	}

	// TODO - context from somewhere else
	sr.ffmpegCommander.Start(context.Background(), sr.removeMeFileToStream)

	voiceSession, err := sr.discordClient.JoinVoiceChat()
	if err != nil {
		return err
	}

	sr.ffmpegCommander.Stream(voiceSession)
	return nil
}
