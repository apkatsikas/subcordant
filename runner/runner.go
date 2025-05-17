package runner

import (
	"context"
	"io"
	"log"

	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/apkatsikas/subcordant/playlist"
)

const trackName = "temptrack"

type SubcordantRunner struct {
	subsonicClient  interfaces.ISubsonicClient
	discordClient   interfaces.IDiscordClient
	ffmpegCommander interfaces.IExecCommander
	*playlist.PlaylistService
	voiceSession io.Writer
	Playing      bool
}

func (sr *SubcordantRunner) Init(
	subsonicClient interfaces.ISubsonicClient, discordClient interfaces.IDiscordClient,
	ffmpegCommander interfaces.IExecCommander) error {
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

func (sr *SubcordantRunner) HandlePlay(albumId string) {
	go func() {
		album, err := sr.subsonicClient.GetAlbum(albumId)

		if err != nil {
			log.Println(err)
		}

		for _, song := range album.Song {
			sr.PlaylistService.Add(song.ID)
		}

		if !sr.Playing {
			go sr.playTracks()
		}
	}()
}

func (sr *SubcordantRunner) playTracks() {
	for {
		sr.Playing = true
		playlist := sr.PlaylistService.GetPlaylist()
		if len(playlist) == 0 {
			sr.Playing = false
			return
		}

		trackId := playlist[0]
		if err := sr.doPlay(trackId); err != nil {
			log.Printf("ERROR: playing track %s resulted in: %v", trackId, err)
			// TODO - Optionally handle errors (e.g., skip to the next track)
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

	err = sr.ffmpegCommander.Start(ctx, stream, trackName, cancel)
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

	return sr.ffmpegCommander.Stream(sr.voiceSession, cancel)
}
