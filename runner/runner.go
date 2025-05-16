package runner

import (
	"context"
	"io"

	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/apkatsikas/subcordant/playlist"
)

const trackName = "temptrack"

type SubcordantRunner struct {
	subsonicClient  interfaces.ISubsonicClient
	discordClient   interfaces.IDiscordClient
	ffmpegCommander interfaces.IFfmpegCommander
	*playlist.PlaylistService
	voiceSession io.Writer
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
	err = sr.doPlay(firstTrack, false)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SubcordantRunner) doPlay(trackId string, nextTrack bool) error {
	stream, err := sr.subsonicClient.Stream(trackId)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-ctx.Done()
		stream.Close()
		sr.PlaylistService.FinishTrack()
		playlist := sr.PlaylistService.GetPlaylist()
		if len(playlist) > 0 {
			sr.doPlay(playlist[0], true)
		}
	}()

	err = sr.ffmpegCommander.Start(ctx, stream, trackName, cancel)
	if err != nil {
		cancel()
		return err
	}

	if !nextTrack {
		voiceSession, err := sr.discordClient.JoinVoiceChat(cancel)
		if err != nil {
			cancel()
			return err
		}
		sr.voiceSession = voiceSession
	}

	err = sr.ffmpegCommander.Stream(sr.voiceSession, cancel)
	if err != nil {
		cancel()
		return err
	}
	return nil
}
