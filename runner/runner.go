package runner

import (
	"fmt"
	"io"

	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/apkatsikas/subcordant/playlist"
)

type SubcordantRunner struct {
	subsonicClient interfaces.ISubsonicClient
	discordClient  interfaces.IDiscordClient
	streamer       interfaces.IStreamer
	*playlist.PlaylistService
	voiceSession io.Writer
	playing      bool
}

func (sr *SubcordantRunner) Init(
	subsonicClient interfaces.ISubsonicClient, discordClient interfaces.IDiscordClient,
	streamer interfaces.IStreamer) error {
	sr.PlaylistService = &playlist.PlaylistService{}
	sr.subsonicClient = subsonicClient
	sr.discordClient = discordClient
	sr.streamer = streamer

	if err := sr.subsonicClient.Init(); err != nil {
		return err
	}

	if err := sr.discordClient.Init(sr); err != nil {
		return err
	}
	return nil
}

func (sr *SubcordantRunner) queue(albumId string) error {
	album, err := sr.subsonicClient.GetAlbum(albumId)
	if err != nil {
		return err
	}

	for _, song := range album.Song {
		sr.PlaylistService.Add(song.ID)
	}
	return nil
}

func (sr *SubcordantRunner) Play(albumId string) error {
	if err := sr.queue(albumId); err != nil {
		return err
	}
	// TODO - mutex
	if sr.playing {
		return nil
	}
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
	streamUrl, err := sr.subsonicClient.StreamUrl(trackId)
	if err != nil {
		return err
	}

	if err := sr.streamer.PrepStream(streamUrl); err != nil {
		return err
	}

	if sr.voiceSession == nil {
		voiceSession, err := sr.discordClient.JoinVoiceChat()
		if err != nil {
			return err
		}
		sr.voiceSession = voiceSession
	}

	return sr.streamer.Stream(sr.voiceSession)
}
