package runner

import (
	"context"

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
	ffmpegCommander interfaces.IFfmpegCommander) error {
	sr.PlaylistService = &playlist.PlaylistService{}
	sr.subsonicClient = subsonicClient
	sr.discordClient = discordClient
	sr.ffmpegCommander = ffmpegCommander

	err := sr.subsonicClient.Init()
	if err != nil {
		return err
	}

	err = sr.discordClient.Init(sr)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SubcordantRunner) HandlePlay(albumId string) error {
	album, err := sr.subsonicClient.GetAlbum(albumId)

	if err != nil {
		return err
	}

	for _, song := range album.Song {
		sr.PlaylistService.Add(song.ID)
	}

	firstTrack := sr.GetPlaylist()[0]
	stream, err := sr.subsonicClient.Stream(firstTrack)
	if err != nil {
		return err
	}
	defer stream.Close()

	// TODO - context from somewhere else
	err = sr.ffmpegCommander.Start(context.Background(), stream, "temptrack")
	if err != nil {
		return err
	}

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
