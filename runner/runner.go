package runner

import (
	"context"
	"log"

	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/apkatsikas/subcordant/playlist"
)

type SubcordantRunner struct {
	subsonicClient  interfaces.ISubsonicClient
	discordClient   interfaces.IDiscordClient
	ffmpegCommander interfaces.IFfmpegCommander
	*playlist.PlaylistService
}

func (sr *SubcordantRunner) Init(
	subsonicClient interfaces.ISubsonicClient, discordClient interfaces.IDiscordClient,
	ffmpegCommander interfaces.IFfmpegCommander) {
	sr.PlaylistService = &playlist.PlaylistService{}
	sr.subsonicClient = subsonicClient
	sr.discordClient = discordClient
	sr.ffmpegCommander = ffmpegCommander

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

	stream, err := sr.subsonicClient.Stream(sr.GetPlaylist()[0])
	if err != nil {
		return err
	}
	defer stream.Close()

	// TODO - context from somewhere else
	sr.ffmpegCommander.Start(context.Background(), stream)

	voiceSession, err := sr.discordClient.JoinVoiceChat()
	if err != nil {
		return err
	}

	err = sr.ffmpegCommander.Stream(voiceSession)
	if err != nil {
		return err
	}
	return nil
}
