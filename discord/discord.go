package discord

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

const (
	playCommand = "play"
	albumId     = "albumid"
)

type DiscordClient struct {
}

var commands = []api.CreateCommandData{
	{
		Name:        playCommand,
		Description: "play an album by album ID",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  albumId,
				Description: "ID of the album",
				Required:    true,
			},
		},
	},
}

func (dc *DiscordClient) Init(commandHandler interfaces.ICommandHandler) error {
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("DISCORD_BOT_TOKEN must be set")
	}

	hand := newHandler(state.New("Bot "+botToken), commandHandler)
	hand.state.AddInteractionHandler(hand)
	hand.state.AddIntents(gateway.IntentGuilds)
	hand.state.AddHandler(func(*gateway.ReadyEvent) {
		me, _ := hand.state.Me()
		log.Println("connected to the gateway as", me.Tag())
	})

	if err := cmdroute.OverwriteCommands(hand.state, commands); err != nil {
		log.Fatalln("cannot update commands:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := hand.state.Connect(ctx); err != nil {
		log.Fatalln("cannot connect:", err)
	}

	return nil
}

type handler struct {
	*cmdroute.Router
	state          *state.State
	commandHandler interfaces.ICommandHandler
}

func (h *handler) cmdPlay(ctx context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	commandAlbumId := cmd.Options.Find(albumId).Value.String()
	h.commandHandler.HandlePlay(commandAlbumId)
	message := fmt.Sprintf("Queueing album with ID of %v", commandAlbumId)
	return &api.InteractionResponseData{
		Content: option.NewNullableString(message),
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
