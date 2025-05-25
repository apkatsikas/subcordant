package runner

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/apkatsikas/subcordant/playlist"
	"github.com/apkatsikas/subcordant/types"
)

type SubcordantRunner struct {
	subsonicClient interfaces.ISubsonicClient
	discordClient  interfaces.IDiscordClient
	streamer       interfaces.IStreamer
	*playlist.PlaylistService
	voiceSession io.Writer
	playing      bool
	mu           sync.Mutex
	cancelPlay   context.CancelFunc
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
		message := fmt.Sprintf("Could not find album with ID of %v", albumId)
		sr.discordClient.SendMessage(message)
		return err
	}

	message := fmt.Sprintf("Queued album: %v", album.Name)
	sr.discordClient.SendMessage(message)

	for _, song := range album.Song {
		sr.PlaylistService.Add(song.ID)
	}
	return nil
}

func (sr *SubcordantRunner) Reset() {
	sr.PlaylistService.Clear()
	sr.cancelPlay()
	sr.voiceSession = nil
	sr.playing = false
}

func (sr *SubcordantRunner) Play(albumId string) (types.PlaybackState, error) {
	ctx, cancel := context.WithCancel(context.Background())
	sr.cancelPlay = cancel
	defer cancel()
	if err := sr.queue(albumId); err != nil {
		return types.Invalid, err
	}
	if sr.checkAndSetPlayingMutex() {
		return types.AlreadyPlaying, nil
	}
	if err := sr.joinVoice(); err != nil {
		sr.PlaylistService.Clear()
		return types.Invalid, err
	}
	for {
		playlist := sr.PlaylistService.GetPlaylist()
		if len(playlist) == 0 {
			sr.playing = false
			return types.PlaybackComplete, nil
		}

		trackId := playlist[0]
		if err := sr.play(ctx, trackId); err != nil {
			sr.PlaylistService.FinishTrack()
			return types.Invalid, fmt.Errorf("playing track %s resulted in: %v", trackId, err)
		}
		sr.PlaylistService.FinishTrack()
	}
}

func (sr *SubcordantRunner) checkAndSetPlayingMutex() bool {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	if sr.playing {
		return true
	}
	sr.playing = true
	return false
}

func (sr *SubcordantRunner) play(context context.Context, trackId string) error {
	streamUrl, err := sr.subsonicClient.StreamUrl(trackId)
	if err != nil {
		return err
	}

	if err := sr.streamer.PrepStream(streamUrl); err != nil {
		return err
	}

	return sr.streamer.Stream(context, sr.voiceSession)
}

func (sr *SubcordantRunner) joinVoice() error {
	if sr.voiceSession == nil {

		voiceSession, err := sr.discordClient.JoinVoiceChat()
		if err != nil {
			return err
		}
		sr.voiceSession = voiceSession
	}
	return nil
}
