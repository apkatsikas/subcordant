package discord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/state/store"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

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

func (h *handler) cmdClear(_ context.Context, _ cmdroute.CommandData) *api.InteractionResponseData {
	h.commandHandler.Reset()
	return &api.InteractionResponseData{
		Content: option.NewNullableString("Clearing subcordant..."),
	}
}

func (h *handler) cmdDisconnect(_ context.Context, _ cmdroute.CommandData) *api.InteractionResponseData {
	h.commandHandler.Disconnect()
	return &api.InteractionResponseData{
		Content: option.NewNullableString("Disconnecting..."),
	}
}

func newHandler(state *state.State, commandHandler interfaces.ICommandHandler) *handler {
	hand := &handler{state: state}
	hand.commandHandler = commandHandler

	hand.Router = cmdroute.NewRouter()
	// Automatically defer handles if they're slow.
	hand.Use(cmdroute.Deferrable(state, cmdroute.DeferOpts{}))
	hand.AddFunc(playCommand, hand.cmdPlay)
	hand.AddFunc(clearCommand, hand.cmdClear)
	hand.AddFunc(disconnectCommand, hand.cmdDisconnect)

	return hand
}
