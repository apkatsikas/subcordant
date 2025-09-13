package discord

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/apkatsikas/subcordant/constants"
	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/state/store"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/diamondburned/arikawa/v3/voice/udp"
)

const (
	playCommand                  = "play"
	playCommandDescription       = "play an album, playlist or track by ID"
	clearCommand                 = "clear"
	clearCommandDescription      = "clears the playlist and stops playback"
	disconnectCommand            = "disconnect"
	disconnectCommandDescription = "disconnects the subcordant bot from the voice channel, " +
		"stopping playback and clearing the playlist"
	skipCommand                       = "skip"
	skipCommandDescription            = "skips the surrently playing song"
	helpCommand                       = "help"
	helpCommandDescription            = "describes all Subcordant commands"
	playAlbumTrackCommand             = "albumtrack"
	playAlbumTrackCommandDescription  = "play a track from an album by albumid and track number"
	playTrackByNameCommand            = "trackname"
	playTrackByNameCommandDescription = "play a track by name"

	optionalSubsonicId  = "subsonicid"
	optionalAlbumId     = "albumid"
	optionalTrackName   = "trackname"
	optionalTrackNumber = "tracknumber"

	// Optional to tweak the Opus stream.
	timeIncrement      = 2880
	dontSwitchChannels = discord.NullChannelID
)

type DiscordClient struct {
	*handler
	voiceSession   *voice.Session
	selfDisconnect bool
	mu             sync.Mutex
}

var commandMap = map[string]string{
	playCommand:            playCommandDescription,
	playAlbumTrackCommand:  playAlbumTrackCommandDescription,
	playTrackByNameCommand: playTrackByNameCommandDescription,
	clearCommand:           clearCommandDescription,
	disconnectCommand:      disconnectCommandDescription,
	skipCommand:            skipCommandDescription,
	helpCommand:            helpCommandDescription,
}

var commandOrder = []string{
	playCommand,
	playAlbumTrackCommand,
	playTrackByNameCommand,
	skipCommand,
	clearCommand,
	disconnectCommand,
	helpCommand,
}

func prettyPrintCommands() string {
	var sb strings.Builder
	sb.WriteString("**Available Commands:**\n\n")

	for _, cmd := range commandOrder {
		desc := commandMap[cmd]
		sb.WriteString(fmt.Sprintf("- **%s**: %s\n", cmd, desc))
	}

	return sb.String()
}

var commands = []api.CreateCommandData{
	{
		Name:        playCommand,
		Description: playCommandDescription,
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  optionalSubsonicId,
				Description: "ID of the subsonic album, playlist or track",
				Required:    true,
			},
		},
	},
	{
		Name:        clearCommand,
		Description: clearCommandDescription,
	},
	{
		Name:        disconnectCommand,
		Description: disconnectCommandDescription,
	},
	{
		Name:        skipCommand,
		Description: skipCommandDescription,
	},
	{
		Name:        helpCommand,
		Description: helpCommandDescription,
	},
	{
		Name:        playAlbumTrackCommand,
		Description: playAlbumTrackCommandDescription,
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  optionalAlbumId,
				Description: "ID of the subsonic album",
				Required:    true,
			},
			&discord.IntegerOption{
				OptionName:  optionalTrackNumber,
				Description: "Number of the track from the album",
				Required:    true,
			},
		},
	},
	{
		Name:        playTrackByNameCommand,
		Description: playTrackByNameCommandDescription,
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  optionalTrackName,
				Description: "Name of the track",
				Required:    true,
			},
		},
	},
}

func createBotAndHandler(commandHandler interfaces.ICommandHandler) (*handler, error) {
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("DISCORD_BOT_TOKEN must be set")
	}

	return newHandler(state.New("Bot "+botToken), commandHandler), nil
}

func setUdpDialer(v *voice.Session) {
	// Optimize Opus frame duration. This step is optional, but it is
	// recommended.
	v.SetUDPDialer(udp.DialFuncWithFrequency(
		constants.FrameDuration*time.Millisecond,
		timeIncrement,
	))
}

func (dc *DiscordClient) GetVoice() io.Writer {
	return dc.voiceSession
}

func (dc *DiscordClient) Init(commandHandler interfaces.ICommandHandler) error {
	hand, err := createBotAndHandler(commandHandler)
	if err != nil {
		return err
	}

	dc.setupHandler(hand)

	if err := cmdroute.OverwriteCommands(dc.handler.state, commands); err != nil {
		return err
	}

	err = dc.connect()
	if err != nil {
		return err
	}

	return nil
}

func (dc *DiscordClient) SendMessage(message string) {
	_, err := dc.state.SendMessage(dc.LastChannelId, message)
	if err != nil {
		log.Printf("\nERROR: send message resulted in %v", err)
	}
}

func (dc *DiscordClient) connect() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := dc.handler.state.Connect(ctx); err != nil {
		return fmt.Errorf("cannot connect: %v", err)
	}
	return nil
}

func (dc *DiscordClient) setupHandler(hand *handler) {
	dc.handler = hand
	dc.handler.state.AddInteractionHandler(dc.handler)
	dc.setupBotDisconnectHandler()
	voice.AddIntents(dc.handler.state)
}

func (dc *DiscordClient) setupBotDisconnectHandler() {
	dc.handler.state.AddHandler(func(event *gateway.VoiceStateUpdateEvent) {
		me, err := dc.handler.state.Me()
		if err != nil {
			log.Printf("ERROR: getting bot state resulted in %v", err)
		}
		isBot := me.ID == event.UserID
		isDisconnect := !event.ChannelID.IsValid()

		if isBot && isDisconnect {
			if dc.selfDisconnect {
				dc.selfDisconnect = false
				return
			}

			dc.commandHandler.Reset()
			dc.LeaveVoiceSession()
		}
	})
}

func (dc *DiscordClient) LeaveVoiceSession() {
	if dc.voiceSession != nil {
		err := dc.voiceSession.Leave(context.Background())
		if err != nil {
			log.Printf("\nERROR: failed to leave voice session: %v", err)
		}
		dc.voiceSession = nil
	}
}

func (dc *DiscordClient) JoinVoiceChat(guildId discord.GuildID, switchToChannel discord.ChannelID) (discord.ChannelID, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bot, err := dc.handler.state.Me()
	if err != nil {
		return dontSwitchChannels, fmt.Errorf("was not able to join a voice channel."+
			"Bot could not get info on itself. Error was %v", err)
	}

	botVoiceState, err := dc.handler.state.VoiceState(guildId, bot.ID)

	if err != nil {
		botNotInVoice := errors.Is(err, store.ErrNotFound)
		if botNotInVoice {
			if err := dc.enterVoice(ctx, switchToChannel); err != nil {
				return dontSwitchChannels, fmt.Errorf("failed to create new session and join: %w", err)
			}
			return dontSwitchChannels, nil
		}
		return dontSwitchChannels, fmt.Errorf("got an unexpected error getting voice state for bot %w", err)
	}

	if !botVoiceState.ChannelID.IsValid() {
		dc.voiceSession.Leave(ctx)
		dc.voiceSession = nil
		return dontSwitchChannels, fmt.Errorf("channel ID of bot was not valid")
	}

	alreadyInChannel := botVoiceState.ChannelID == switchToChannel
	if alreadyInChannel {
		return dontSwitchChannels, nil
	}

	return switchToChannel, nil
}

func (dc *DiscordClient) SwitchVoiceChannel(channelId discord.ChannelID) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	dc.selfDisconnect = true
	err := dc.voiceSession.Leave(ctx)
	if err != nil {
		return fmt.Errorf("bot failed to leave channel it was in %v", err)
	}

	if err := dc.enterVoice(ctx, channelId); err != nil {
		return fmt.Errorf("failed to create new session and join: %w", err)
	}
	return nil
}

func (dc *DiscordClient) enterVoice(ctx context.Context, channelId discord.ChannelID) error {
	v, err := voice.NewSession(dc.state)
	if err != nil {
		return fmt.Errorf("cannot make new voice session: %w", err)
	}

	setUdpDialer(v)
	dc.voiceSession = v

	if err := dc.voiceSession.JoinChannelAndSpeak(ctx, channelId, false, true); err != nil {
		dc.voiceSession.Leave(ctx)
		dc.voiceSession = nil
		return fmt.Errorf("failed to join channel: %w", err)
	}
	return nil
}
