package discord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
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
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/diamondburned/arikawa/v3/voice/udp"
)

const (
	playCommand   = "play"
	optionAlbumId = "albumid"

	// Optional to tweak the Opus stream.
	timeIncrement = 2880
)

type DiscordClient struct {
	*handler
	voiceSession   *voice.Session
	selfDisconnect bool
	mu             sync.Mutex
}

var commands = []api.CreateCommandData{
	{
		Name:        playCommand,
		Description: "play an album by album ID",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  optionAlbumId,
				Description: "ID of the album",
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
			if dc.voiceSession != nil {
				err := dc.voiceSession.Leave(context.Background())
				if err != nil {
					log.Printf("\nERROR: failed to leave voice session: %v", err)
				}
				dc.voiceSession = nil
			}
		}
	})
}

// Returning a channel ID indicates that the user is asking the bot to change channels
func (dc *DiscordClient) JoinVoiceChat(guildId discord.GuildID, channelId discord.ChannelID) (discord.ChannelID, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bot, err := dc.handler.state.Me()
	if err != nil {
		return discord.NullChannelID, fmt.Errorf("was not able to join a voice channel."+
			"Bot could not get info on itself. Error was %v", err)
	}

	botVoiceState, botVoiceStateErr := dc.handler.state.VoiceState(guildId, bot.ID)
	botNotInVoice := botVoiceStateErr != nil
	if botNotInVoice {
		if err := dc.newSessionAndJoin(ctx, channelId); err != nil {
			return discord.NullChannelID, fmt.Errorf("failed to create new session and join: %w", err)
		}
		return discord.NullChannelID, nil
	}
	if !botVoiceState.ChannelID.IsValid() {
		dc.voiceSession.Leave(ctx)
		dc.voiceSession = nil
		return discord.NullChannelID, fmt.Errorf("channel ID of bot was not valid")
	}

	alreadyInChannel := botVoiceState.ChannelID == channelId
	if alreadyInChannel {
		return discord.NullChannelID, nil
	}

	return channelId, nil
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

	if err := dc.newSessionAndJoin(ctx, channelId); err != nil {
		return fmt.Errorf("failed to create new session and join: %w", err)
	}
	return nil
}

func (dc *DiscordClient) newSessionAndJoin(ctx context.Context, channelId discord.ChannelID) error {
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

type handler struct {
	*cmdroute.Router
	state          *state.State
	commandHandler interfaces.ICommandHandler
	LastChannelId  discord.ChannelID
}

func (h *handler) cmdPlay(ctx context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	var albumId string
	err := json.Unmarshal(cmd.Options.Find(optionAlbumId).Value, &albumId)
	if err != nil {
		errorMessage := fmt.Sprintf("ERROR: Failed to unmarshal JSON: %v", err)

		return &api.InteractionResponseData{
			Content: option.NewNullableString(errorMessage),
		}
	}

	vs, err := h.state.VoiceState(cmd.Event.GuildID, cmd.Event.Member.User.ID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return &api.InteractionResponseData{
				Content: option.NewNullableString("User sending command must be in a voice channel, but was not found in one."),
			}
		}
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("Failed to get voice state: %v", err)),
		}
	}
	if !vs.ChannelID.IsValid() {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("Channel ID was not valid: %v", err)),
		}
	}

	h.LastChannelId = cmd.Event.ChannelID

	go h.play(albumId, cmd.Event.GuildID, vs.ChannelID)

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("Recieved %v command with albumid of %v", cmd.Name, albumId)),
	}
}

func (h *handler) play(albumId string, guildId discord.GuildID, channelId discord.ChannelID) {
	if _, err := h.commandHandler.Play(albumId, guildId, channelId); err != nil {
		log.Printf("\nERROR: Play resulted in %v", err)
	}
}

func newHandler(state *state.State, commandHandler interfaces.ICommandHandler) *handler {
	hand := &handler{state: state}
	hand.commandHandler = commandHandler

	hand.Router = cmdroute.NewRouter()
	// Automatically defer handles if they're slow.
	hand.Use(cmdroute.Deferrable(state, cmdroute.DeferOpts{}))
	hand.AddFunc(playCommand, hand.cmdPlay)

	return hand
}
