package runner

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/apkatsikas/subcordant/playlist"
	"github.com/apkatsikas/subcordant/subsonic"
	"github.com/apkatsikas/subcordant/types"
	flagutil "github.com/apkatsikas/subcordant/util/flag"
	"github.com/diamondburned/arikawa/v3/discord"
)

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
		sr.PlaylistService.Add(subsonic.ToTrack(song))
	}
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

// TODO Test last song, test after 2 plays, test no songs
// no user in voice, user in different channel
// skip after skip
// Run with race
// add tests (race)
func (sr *SubcordantRunner) Skip() {
	sr.mu.Lock()
	if sr.cancelPlay != nil {
		sr.cancelPlay()
	}
	time.Sleep(time.Millisecond * 2000)
	// TODO - fix issue
	//Got an error trying to play after skip playing track resulted in: failed to decode ogg: read |0: file already closed
	sr.playing = false
	sr.FinishTrack()
	// TODO - variable
	time.Sleep(time.Millisecond * 750)
	log.Printf("Skip()")
	sr.mu.Unlock()
	err := sr.playFromSkip()
	if err != nil {
		sr.discordClient.SendMessage(fmt.Sprintf("Got an error trying to play after skip %v", err))
	}
}

func (sr *SubcordantRunner) Disconnect() {
	sr.Reset()
	sr.discordClient.LeaveVoiceSession()
}

func (sr *SubcordantRunner) Play(albumId string, guildId discord.GuildID, switchToChannel discord.ChannelID) (types.PlaybackState, error) {
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
		err := sr.discordClient.SwitchVoiceChannel(switchToChannel)
		if err != nil {
			sr.Reset()
			sr.discordClient.SendMessage(fmt.Sprintf("Failed to switch channels, error is %v", err))
			return types.Invalid, err
		}
	}

	if err := sr.queue(albumId); err != nil {
		return types.Invalid, err
	}
	if sr.checkAndSetPlayingMutex() {
		return types.AlreadyPlaying, nil
	}

	for {
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

func (sr *SubcordantRunner) playFromSkip() error {
	log.Printf("playFromSkip()")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Printf("going to checkAndSetPlayingMutex")
	sr.checkAndSetPlayingMutex()
	log.Printf("checked mutex and set")
	for {
		log.Printf("In for")
		sr.mu.Lock()
		sr.cancelPlay = cancel
		if len(sr.PlaylistService.GetPlaylist()) == 0 {
			sr.playing = false
			sr.mu.Unlock()
			return fmt.Errorf("tried to play but no songs are left on the playlist")
		}

		trackId := sr.PlaylistService.GetPlaylist()[0]
		sr.mu.Unlock()
		log.Printf("gonna play from skip!")
		if err := sr.play(ctx, trackId); err != nil {
			sr.PlaylistService.FinishTrack()
			return fmt.Errorf("playing track %s resulted in: %v", trackId, err)
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
