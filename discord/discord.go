package discord

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/apkatsikas/subcordant/interfaces"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/godave/golibdave"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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
	playAlbumByNameCommand            = "albumname"
	playAlbumByNameCommandDescription = "play an album by name"

	optionalSubsonicId  = "subsonicid"
	optionalAlbumId     = "albumid"
	optionalTrackName   = "trackname"
	optionalTrackNumber = "tracknumber"
	optionalAlbumName   = "albumname"

	dontSwitchChannels = 0

	unknownVoiceState = "10065: Unknown Voice State"
)

var (
	commandTimeout = 3 * time.Second
	commandMap     = map[string]string{
		playCommand:            playCommandDescription,
		playAlbumTrackCommand:  playAlbumTrackCommandDescription,
		playTrackByNameCommand: playTrackByNameCommandDescription,
		playAlbumByNameCommand: playAlbumByNameCommandDescription,
		clearCommand:           clearCommandDescription,
		disconnectCommand:      disconnectCommandDescription,
		skipCommand:            skipCommandDescription,
		helpCommand:            helpCommandDescription,
	}
	commandOrder = []string{
		playCommand,
		playAlbumTrackCommand,
		playTrackByNameCommand,
		playAlbumByNameCommand,
		skipCommand,
		clearCommand,
		disconnectCommand,
		helpCommand,
	}
	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        playCommand,
			Description: playCommandDescription,
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        optionalSubsonicId,
					Description: "ID of the subsonic album, playlist or track",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        clearCommand,
			Description: clearCommandDescription,
		},
		discord.SlashCommandCreate{
			Name:        disconnectCommand,
			Description: disconnectCommandDescription,
		},
		discord.SlashCommandCreate{
			Name:        skipCommand,
			Description: skipCommandDescription,
		},
		discord.SlashCommandCreate{
			Name:        helpCommand,
			Description: helpCommandDescription,
		},
		discord.SlashCommandCreate{
			Name:        playAlbumTrackCommand,
			Description: playAlbumTrackCommandDescription,
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        optionalAlbumId,
					Description: "ID of the subsonic album",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        optionalTrackNumber,
					Description: "Number of the track from the album",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        playTrackByNameCommand,
			Description: playTrackByNameCommandDescription,
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        optionalTrackName,
					Description: "Name of the track",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        playAlbumByNameCommand,
			Description: playAlbumByNameCommandDescription,
			Options: []discord.ApplicationCommandOption{
				&discord.ApplicationCommandOptionString{
					Name:        optionalAlbumName,
					Description: "Name of the album",
					Required:    true,
				},
			},
		},
	}
)

type DiscordClient struct {
	selfDisconnect bool
	mu             sync.Mutex
	bot            *bot.Client
	guildId        snowflake.ID
	commandHandler interfaces.ICommandHandler
	LastChannelId  snowflake.ID
}

func (dc *DiscordClient) Init(commandHandler interfaces.ICommandHandler) error {
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("DISCORD_BOT_TOKEN must be set")
	}

	router := handler.New()
	router.Group(func(router handler.Router) {
		router.Use(handleError())

		router.Group(func(router handler.Router) {
			router.Use(dc.handleLastChannelId)
			router.Command(fmtCmd(playCommand), dc.handlePlay)
			router.Command(fmtCmd(playAlbumTrackCommand), dc.handlePlayTrackFromAlbum)
			router.Command(fmtCmd(playTrackByNameCommand), dc.handlePlayTrackByName)
			router.Command(fmtCmd(playAlbumByNameCommand), dc.handlePlayAlbumByName)

			router.Command(fmtCmd(skipCommand), dc.handleSkip)

			router.Command(fmtCmd(clearCommand), dc.handleClear)
			router.Command(fmtCmd(disconnectCommand), dc.handleDisconnect)
		})
		router.Command(fmtCmd(helpCommand), handleHelp)
	})
	router.NotFound(handleNotFound)

	client, err := disgo.New(botToken,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildVoiceStates, gateway.IntentGuilds)),
		bot.WithEventListeners(router),
		bot.WithEventListenerFunc(dc.onGuildVoiceLeave),
		bot.WithEventListenerFunc(dc.onGuildReady),
		bot.WithVoiceManagerConfigOpts(
			voice.WithDaveSessionCreateFunc(golibdave.NewSession),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create new discord bot: %v", err)
	}

	dc.bot = client
	dc.commandHandler = commandHandler

	if err = dc.openGateway(); err != nil {
		return fmt.Errorf("failed to open gateway: %v", err)
	}

	return nil
}

func (dc *DiscordClient) SendMessage(message string) {
	_, err := dc.bot.Rest.CreateMessage(dc.LastChannelId, discord.NewMessageCreateV2(discord.NewTextDisplay(message)))
	if err != nil {
		log.Printf("ERROR: send message resulted in %v", err)
	}
}

func (dc *DiscordClient) LeaveVoiceSession() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()

	if voice := dc.bot.VoiceManager.GetConn(dc.guildId); voice != nil {
		voice.Close(ctx)
	}
}

func (dc *DiscordClient) SetFrameProvider(frameProvider voice.OpusFrameProvider) {
	dc.bot.VoiceManager.GetConn(dc.guildId).SetOpusFrameProvider(frameProvider)
}

func (dc *DiscordClient) JoinVoiceChat(switchToChannel snowflake.ID) (snowflake.ID, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()

	voice := dc.bot.VoiceManager.GetConn(dc.guildId)

	if voice == nil {
		if err := dc.enterVoice(ctx, switchToChannel); err != nil {
			return dontSwitchChannels, fmt.Errorf("failed to create new session and join: %w", err)
		}
		return dontSwitchChannels, nil
	}

	if *voice.ChannelID() == 0 {
		voice.Close(ctx)
		return dontSwitchChannels, fmt.Errorf("channel ID of bot was not valid")
	}

	alreadyInChannel := *voice.ChannelID() == switchToChannel
	if alreadyInChannel {
		return dontSwitchChannels, nil
	}

	return switchToChannel, nil
}

func (dc *DiscordClient) SwitchVoiceChannel(channelId snowflake.ID) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()

	dc.selfDisconnect = true
	dc.bot.VoiceManager.Close(ctx)

	if err := dc.enterVoice(ctx, channelId); err != nil {
		return fmt.Errorf("failed to create new session and join: %w", err)
	}
	return nil
}

func (dc *DiscordClient) Shutdown() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()
	dc.bot.Close(ctx)
}

func (dc *DiscordClient) openGateway() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()

	if err := dc.bot.OpenGateway(ctx); err != nil {
		return fmt.Errorf("failed to open gateway: %v", err)
	}

	return nil
}

func (dc *DiscordClient) onGuildVoiceLeave(event *events.GuildVoiceLeave) {
	isSelf := event.Member.User.ID == dc.bot.ID()
	if isSelf {
		if dc.selfDisconnect {
			dc.selfDisconnect = false
			return
		}

		dc.commandHandler.Reset()
		dc.LeaveVoiceSession()
	}
}

func (dc *DiscordClient) onGuildReady(e *events.GuildReady) {
	dc.guildId = e.GuildID
	if err := handler.SyncCommands(dc.bot, commands, []snowflake.ID{dc.guildId}); err != nil {
		log.Panicf("failed to sync commands: %v", err)
	}
}

func (dc *DiscordClient) handleClear(event *handler.CommandEvent) error {
	go dc.commandHandler.Reset()
	return event.CreateMessage(discord.MessageCreate{
		Content: "Clearing subcordant...",
	})
}

func (dc *DiscordClient) handleDisconnect(event *handler.CommandEvent) error {
	go dc.commandHandler.Disconnect()
	return event.CreateMessage(discord.MessageCreate{
		Content: "Disconnecting...",
	})
}

func (dc *DiscordClient) handlePlay(event *handler.CommandEvent) error {
	data := event.SlashCommandInteractionData()
	subsonicId := data.String(optionalSubsonicId)

	if err := dc.runPlayCommand(
		event,
		func(channelId snowflake.ID) error {
			_, err := dc.commandHandler.Play(subsonicId, channelId)
			return err
		},
	); err != nil {
		return err
	}
	return event.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Received %v command with subsonic ID of %v", data.CommandName(), subsonicId),
	})
}

func (dc *DiscordClient) handlePlayAlbumByName(event *handler.CommandEvent) error {
	data := event.SlashCommandInteractionData()
	albumName := data.String(optionalAlbumName)

	if err := dc.runPlayCommand(
		event,
		func(channelId snowflake.ID) error {
			_, err := dc.commandHandler.PlayAlbumByName(albumName, channelId)
			return err
		},
	); err != nil {
		return err
	}
	return event.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Received %v command with album name of %v", data.CommandName(), albumName),
	})
}

func (dc *DiscordClient) handlePlayTrackFromAlbum(event *handler.CommandEvent) error {
	data := event.SlashCommandInteractionData()

	albumId := data.String(optionalAlbumId)
	trackNumber := data.Int(optionalTrackNumber)

	if err := dc.runPlayCommand(event, func(channelId snowflake.ID) error {
		_, err := dc.commandHandler.PlayTrackFromAlbum(albumId, trackNumber, channelId)
		return err
	}); err != nil {
		return err
	}

	msg := fmt.Sprintf("Received %v command with album ID %v and track number %d", data.CommandName(), albumId, trackNumber)
	return event.CreateMessage(discord.MessageCreate{Content: msg})
}

func (dc *DiscordClient) handlePlayTrackByName(event *handler.CommandEvent) error {
	data := event.SlashCommandInteractionData()

	trackName := data.String(optionalTrackName)

	if err := dc.runPlayCommand(event, func(channelId snowflake.ID) error {
		_, err := dc.commandHandler.PlayTrackByName(trackName, channelId)
		return err
	}); err != nil {
		return err
	}

	msg := fmt.Sprintf("Received %v command with track name %v", data.CommandName(), trackName)
	return event.CreateMessage(discord.MessageCreate{Content: msg})
}

func (dc *DiscordClient) handleSkip(event *handler.CommandEvent) error {
	go dc.commandHandler.Skip()
	return event.CreateMessage(discord.MessageCreate{Content: "Skipping track..."})
}

func (dc *DiscordClient) handleLastChannelId(next handler.Handler) handler.Handler {
	return func(event *handler.InteractionEvent) error {
		dc.LastChannelId = event.Channel().ID()
		return next(event)
	}
}

func (dc *DiscordClient) runPlayCommand(
	event *handler.CommandEvent,
	playFn func(channelId snowflake.ID) error,
) error {
	vs, err := dc.bot.Rest.GetUserVoiceState(dc.guildId, event.User().ID)
	if err != nil {
		if err.Error() == unknownVoiceState {
			return fmt.Errorf("unknown voice state - the user sending the command is probably not in a voice channel")
		}
		return fmt.Errorf("failed to get user voice state: %v", err)
	}
	if *vs.ChannelID == 0 {
		return fmt.Errorf("channel ID was not valid")
	}

	go func() {
		if err := playFn(*vs.ChannelID); err != nil {
			log.Printf("ERROR: playFn resulted in %v", err)
		}
	}()

	return nil
}

func (dc *DiscordClient) enterVoice(ctx context.Context, channelId snowflake.ID) error {
	v := dc.bot.VoiceManager.CreateConn(dc.guildId)

	err := v.Open(ctx, channelId, false, true)

	if err != nil {
		return fmt.Errorf("cannot open voice session: %w", err)
	}

	if err := v.SetSpeaking(ctx, voice.SpeakingFlagMicrophone); err != nil {
		v.Close(ctx)
		return fmt.Errorf("error setting speaking flag: %v", err)
	}
	return nil
}

func fmtCmd(cmd string) string {
	return fmt.Sprintf("/%v", cmd)
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

func handleNotFound(event *handler.InteractionEvent) error {
	return event.CreateMessage(discord.MessageCreate{Content: "command not found"})
}

func handleHelp(e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Content: prettyPrintCommands(),
	})
}

func handleError() func(next handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return func(event *handler.InteractionEvent) error {
			if err := next(event); err != nil {
				log.Printf("command error: %v on event", err)

				if createMsgErr := event.CreateMessage(discord.MessageCreate{
					Content: fmt.Sprintf("❌ Something went wrong while running command: %v", err),
				}); createMsgErr != nil {
					log.Printf("failed to send error message: %v - original error was: %v", createMsgErr, err)
				}

				return nil
			}
			return nil
		}
	}
}
