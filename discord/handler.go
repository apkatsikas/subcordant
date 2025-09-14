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

func (h *handler) cmdPlay(_ context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	var subsonicId string
	if err := json.Unmarshal(cmd.Options.Find(optionalSubsonicId).Value, &subsonicId); err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("ERROR: Failed to unmarshal JSON: %v", err)),
		}
	}

	return h.runPlayCommand(
		cmd,
		func(guildId discord.GuildID, channelId discord.ChannelID) error {
			_, err := h.commandHandler.Play(subsonicId, guildId, channelId)
			return err
		},
		fmt.Sprintf("Received %v command with subsonic ID of %v", cmd.Name, subsonicId),
	)
}

func (h *handler) cmdPlayTrackFromAlbum(_ context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	var albumId string
	if err := json.Unmarshal(cmd.Options.Find("albumId").Value, &albumId); err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("ERROR: Failed to unmarshal albumId: %v", err)),
		}
	}

	var trackNumber int
	if err := json.Unmarshal(cmd.Options.Find("trackNumber").Value, &trackNumber); err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("ERROR: Failed to unmarshal trackNumber: %v", err)),
		}
	}

	return h.runPlayCommand(
		cmd,
		func(guildId discord.GuildID, channelId discord.ChannelID) error {
			_, err := h.commandHandler.PlayTrackFromAlbum(albumId, trackNumber, guildId, channelId)
			return err
		},
		fmt.Sprintf("Received %v command with album ID %v and track number %d", cmd.Name, albumId, trackNumber),
	)
}

func (h *handler) cmdPlayAlbumByName(_ context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	var albumName string
	if err := json.Unmarshal(cmd.Options.Find("albumName").Value, &albumName); err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("ERROR: Failed to unmarshal albumName: %v", err)),
		}
	}

	return h.runPlayCommand(
		cmd,
		func(guildId discord.GuildID, channelId discord.ChannelID) error {
			_, err := h.commandHandler.PlayAlbumByName(albumName, guildId, channelId)
			return err
		},
		fmt.Sprintf("Received %v command with album name %v", cmd.Name, albumName),
	)
}

func (h *handler) cmdPlayTrackByName(_ context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	var trackName string
	if err := json.Unmarshal(cmd.Options.Find("trackName").Value, &trackName); err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("ERROR: Failed to unmarshal trackName: %v", err)),
		}
	}

	return h.runPlayCommand(
		cmd,
		func(guildId discord.GuildID, channelId discord.ChannelID) error {
			_, err := h.commandHandler.PlayTrackByName(trackName, guildId, channelId)
			return err
		},
		fmt.Sprintf("Received %v command with track name %v", cmd.Name, trackName),
	)
}

func (h *handler) runPlayCommand(
	cmd cmdroute.CommandData,
	playFn func(guildId discord.GuildID, channelId discord.ChannelID) error,
	responseMsg string,
) *api.InteractionResponseData {
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

	go func() {
		if err := playFn(cmd.Event.GuildID, vs.ChannelID); err != nil {
			log.Printf("\nERROR: playFn resulted in %v", err)
		}
	}()

	return &api.InteractionResponseData{
		Content: option.NewNullableString(responseMsg),
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

func (h *handler) cmdSkip(_ context.Context, _ cmdroute.CommandData) *api.InteractionResponseData {
	go h.commandHandler.Skip()
	return &api.InteractionResponseData{
		Content: option.NewNullableString("Skipping track..."),
	}
}

func (h *handler) cmdHelp(_ context.Context, _ cmdroute.CommandData) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content: option.NewNullableString(prettyPrintCommands()),
	}
}

func newHandler(state *state.State, commandHandler interfaces.ICommandHandler) *handler {
	hand := &handler{state: state}
	hand.commandHandler = commandHandler

	hand.Router = cmdroute.NewRouter()
	// Automatically defer handles if they're slow.
	hand.Use(cmdroute.Deferrable(state, cmdroute.DeferOpts{}))
	hand.AddFunc(playCommand, hand.cmdPlay)
	hand.AddFunc(playAlbumTrackCommand, hand.cmdPlayTrackFromAlbum)
	hand.AddFunc(playTrackByNameCommand, hand.cmdPlayTrackByName)
	hand.AddFunc(playAlbumByNameCommand, hand.cmdPlayAlbumByName)
	hand.AddFunc(clearCommand, hand.cmdClear)
	hand.AddFunc(disconnectCommand, hand.cmdDisconnect)
	hand.AddFunc(skipCommand, hand.cmdSkip)
	hand.AddFunc(helpCommand, hand.cmdHelp)

	return hand
}
