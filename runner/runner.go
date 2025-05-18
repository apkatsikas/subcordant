package runner

import (
	"context"
	"fmt"
	"io"

	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/apkatsikas/subcordant/playlist"
)

const trackName = "/tmp/music_pipe"

type SubcordantRunner struct {
	subsonicClient interfaces.ISubsonicClient
	discordClient  interfaces.IDiscordClient
	execCommander  interfaces.IExecCommander
	*playlist.PlaylistService
	voiceSession io.Writer
	playing      bool
}

func (sr *SubcordantRunner) Init(
	subsonicClient interfaces.ISubsonicClient, discordClient interfaces.IDiscordClient,
	execCommander interfaces.IExecCommander) error {
	sr.PlaylistService = &playlist.PlaylistService{}
	sr.subsonicClient = subsonicClient
	sr.discordClient = discordClient
	sr.execCommander = execCommander

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

func (sr *SubcordantRunner) IsPlaying() bool {
	return sr.playing
}

func (sr *SubcordantRunner) Queue(albumId string) error {
	album, err := sr.subsonicClient.GetAlbum(albumId)

	if err != nil {
		return err
	}

	for _, song := range album.Song {
		sr.PlaylistService.Add(song.ID)
	}
	return nil
}

func (sr *SubcordantRunner) Play() error {
	for {
		sr.playing = true
		playlist := sr.PlaylistService.GetPlaylist()
		if len(playlist) == 0 {
			sr.playing = false
			return nil
		}

		trackId := playlist[0]
		if err := sr.doPlay(trackId); err != nil {
			sr.PlaylistService.FinishTrack()
			return fmt.Errorf("playing track %s resulted in: %v", trackId, err)
		}
		sr.PlaylistService.FinishTrack()
	}
}

func (sr *SubcordantRunner) doPlay(trackId string) error {
	stream, err := sr.subsonicClient.Stream(trackId)
	if err != nil {
		return err
	}
	defer stream.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = sr.execCommander.Start(ctx, stream, trackName, cancel)
	if err != nil {
		return err
	}

	if sr.voiceSession == nil {
		voiceSession, err := sr.discordClient.JoinVoiceChat(cancel)
		if err != nil {
			return err
		}
		sr.voiceSession = voiceSession
	}

	return sr.execCommander.Stream(sr.voiceSession, cancel)
}
