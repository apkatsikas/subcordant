package runner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/apkatsikas/subcordant/playlist"
	"github.com/apkatsikas/subcordant/subsonic"
	"github.com/apkatsikas/subcordant/types"
	flagutil "github.com/apkatsikas/subcordant/util/flag"
	"github.com/diamondburned/arikawa/v3/discord"
)

var timeBetweenSkips = time.Millisecond * 3500

type SubcordantRunner struct {
	subsonicClient interfaces.ISubsonicClient
	discordClient  interfaces.IDiscordClient
	streamer       interfaces.IStreamer
	*playlist.PlaylistService
	playing    bool
	mu         sync.Mutex
	cancelPlay context.CancelFunc
	flagutil.StreamFrom
}

func (sr *SubcordantRunner) Init(
	subsonicClient interfaces.ISubsonicClient, discordClient interfaces.IDiscordClient,
	streamer interfaces.IStreamer, streamFrom flagutil.StreamFrom) error {
	sr.PlaylistService = &playlist.PlaylistService{}
	sr.subsonicClient = subsonicClient
	sr.discordClient = discordClient
	sr.streamer = streamer
	sr.StreamFrom = streamFrom

	if err := sr.subsonicClient.Init(); err != nil {
		return err
	}

	if err := sr.discordClient.Init(sr); err != nil {
		return err
	}
	return nil
}

func (sr *SubcordantRunner) queue(subsonicId string) error {
	tracks, err := sr.subsonicClient.GetTracks(subsonicId)
	if err != nil {
		sr.discordClient.SendMessage(err.Error())
		return err
	}

	message := fmt.Sprintf("Queued tracks: %v", tracks.Name)
	sr.discordClient.SendMessage(message)

	for _, song := range tracks.Tracks {
		sr.PlaylistService.Add(subsonic.ToTrack(song))
	}
	return nil
}

func (sr *SubcordantRunner) queueTrackFromAlbum(subsonicId string, trackNumber int) error {
	track, err := sr.subsonicClient.GetTrackFromAlbum(subsonicId, trackNumber)
	if err != nil {
		sr.discordClient.SendMessage(err.Error())
		return err
	}

	message := fmt.Sprintf("Queued track: %v", track.Title)
	sr.discordClient.SendMessage(message)

	sr.PlaylistService.Add(subsonic.ToTrack(track))
	return nil
}

func (sr *SubcordantRunner) Reset() {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.PlaylistService.Clear()
	if sr.cancelPlay != nil {
		sr.cancelPlay()
	}
	sr.playing = false
}

func (sr *SubcordantRunner) Skip() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sr.mu.Lock()
	sr.FinishTrack()
	remaining := sr.GetPlaylist()
	sr.Clear()
	if sr.cancelPlay != nil {
		sr.cancelPlay()
	}
	sr.playing = false
	sr.mu.Unlock()
	time.Sleep(timeBetweenSkips)
	for _, track := range remaining {
		sr.Add(track)
	}
	_, err := sr.playLooper(ctx, cancel)
	if err != nil {
		sr.discordClient.SendMessage(fmt.Sprintf("Got an error trying to play after skip %v", err))
	}
}

func (sr *SubcordantRunner) Disconnect() {
	sr.Reset()
	sr.discordClient.LeaveVoiceSession()
}

func (sr *SubcordantRunner) playWithQueue(
	guildId discord.GuildID,
	switchToChannel discord.ChannelID,
	queueFn func() error,
) (types.PlaybackState, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switchToChannel, err := sr.discordClient.JoinVoiceChat(guildId, switchToChannel)
	if err != nil {
		sr.Reset()
		sr.discordClient.SendMessage(fmt.Sprintf("Could not join voice, error is %v", err))
		return types.Invalid, err
	}

	wantToSwitchChannels := switchToChannel.IsValid()
	if wantToSwitchChannels {
		sr.Reset()
		if err := sr.discordClient.SwitchVoiceChannel(switchToChannel); err != nil {
			sr.Reset()
			sr.discordClient.SendMessage(fmt.Sprintf("Failed to switch channels, error is %v", err))
			return types.Invalid, err
		}
	}

	if err := queueFn(); err != nil {
		return types.Invalid, err
	}

	return sr.playLooper(ctx, cancel)
}

func (sr *SubcordantRunner) Play(subsonicId string, guildId discord.GuildID, switchToChannel discord.ChannelID) (types.PlaybackState, error) {
	return sr.playWithQueue(guildId, switchToChannel, func() error {
		return sr.queue(subsonicId)
	})
}

func (sr *SubcordantRunner) PlayTrackFromAlbum(subsonicId string, trackNumber int, guildId discord.GuildID, switchToChannel discord.ChannelID) (types.PlaybackState, error) {
	return sr.playWithQueue(guildId, switchToChannel, func() error {
		return sr.queueTrackFromAlbum(subsonicId, trackNumber)
	})
}

func (sr *SubcordantRunner) playLooper(ctx context.Context, cancel context.CancelFunc) (types.PlaybackState, error) {
	if sr.checkAndSetPlayingMutex() {
		return types.AlreadyPlaying, nil
	}
	for {
		select {
		case <-ctx.Done():
			return types.PlaybackComplete, nil
		default:
		}
		sr.mu.Lock()
		// Set the cancel play function after any potential SubcordantRunner Reset functions are called
		// this way it cancels the previous function (ongoing playback) if needed
		// and we avoid overwriting it accidentally with the new value
		sr.cancelPlay = cancel
		if len(sr.PlaylistService.GetPlaylist()) == 0 {
			sr.playing = false
			sr.mu.Unlock()
			return types.PlaybackComplete, nil
		}

		trackId := sr.PlaylistService.GetPlaylist()[0]
		sr.mu.Unlock()
		if err := sr.play(ctx, trackId); err != nil {
			sr.PlaylistService.FinishTrack()
			return types.Invalid, fmt.Errorf("playing track %s resulted in: %v", trackId, err)
		}
		sr.mu.Lock()
		sr.PlaylistService.FinishTrack()
		sr.mu.Unlock()
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

func (sr *SubcordantRunner) play(context context.Context, track types.Track) error {
	if sr.StreamFrom == flagutil.StreamFromFile {
		if err := sr.streamer.PrepStreamFromFile(track.Path); err != nil {
			return err
		}
		return sr.streamer.Stream(context, sr.discordClient.GetVoice())
	}

	streamUrl, err := sr.subsonicClient.StreamUrl(track.ID)
	if err != nil {
		return err
	}

	if err := sr.streamer.PrepStreamFromStream(streamUrl); err != nil {
		return err
	}
	return sr.streamer.Stream(context, sr.discordClient.GetVoice())
}
