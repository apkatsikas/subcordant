package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/apkatsikas/subcordant/constants"
	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
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
	voiceChannelId discord.Snowflake
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

func getChannelId() (discord.Snowflake, error) {
	idStr := os.Getenv("DISCORD_VOICE_CHANNEL_ID")
	if idStr == "" {
		return 0, fmt.Errorf("DISCORD_VOICE_CHANNEL_ID must be set")
	}
	id, err := discord.ParseSnowflake(idStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert channel id %v to Snowflake", idStr)
	}
	return id, nil
}

func setUdpDialer(v *voice.Session) {
	// Optimize Opus frame duration. This step is optional, but it is
	// recommended.
	v.SetUDPDialer(udp.DialFuncWithFrequency(
		constants.FrameDuration*time.Millisecond,
		timeIncrement,
	))
}

func (dc *DiscordClient) Init(commandHandler interfaces.ICommandHandler) error {
	hand, err := createBotAndHandler(commandHandler)
	if err != nil {
		return err
	}
	voiceChannelId, err := getChannelId()
	if err != nil {
		return err
	}
	dc.voiceChannelId = voiceChannelId

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
			dc.commandHandler.Reset()
		}
	})
}

func (dc *DiscordClient) JoinVoiceChat() (io.Writer, error) {
	v, err := voice.NewSession(dc.state)
	if err != nil {
		return nil, fmt.Errorf("cannot make new voice session: %w", err)
	}

	setUdpDialer(v)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := v.JoinChannelAndSpeak(ctx, discord.ChannelID(dc.voiceChannelId), false, true); err != nil {
		v.Leave(ctx)
		return nil, fmt.Errorf("failed to join channel: %w", err)
	}

	return v, nil
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

	h.LastChannelId = cmd.Event.ChannelID

	go h.play(albumId)

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("Recieved %v command with albumid of %v", cmd.Name, albumId)),
	}
}

func (h *handler) play(albumId string) {
	if _, err := h.commandHandler.Play(albumId); err != nil {
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
