# subcordant

## TODOs
* Don't emit that ffmpeg was cancelled - remove TODO
* Make issues on GitHub
* Update README
* DC during song, during no song
* set idle disconnect handler - after X amount of time without playing, DC and reset
* eventual code like
func (h *handler) play(albumId string) {
	if _, err := h.commandHandler.Play(albumId); err != nil {

		id, _ := discord.ParseSnowflake("1371301075998740483")

		h.state.SendMessage(discord.ChannelID(id), fmt.Sprintf("Failed to play album with ID of %v", albumId))
		log.Printf("\nERROR: Play resulted in %v", err)
	}
}
to send a message if album is not found
* When disconnected via Discord, cleanly exit and clear playlist

func (dc *DiscordClient) setupHandler(hand *handler) {
	dc.handler = hand
	dc.handler.state.AddInteractionHandler(dc.handler)

	ch := make(chan *gateway.VoiceStateUpdateEvent)
	hand.state.AddHandler(ch)

	go func() {
		for event := range ch {
			me, err := dc.handler.state.Me()
			if err != nil {
				fmt.Printf("ERROR: could not determine ME")
				return
			}
			if !event.ChannelID.IsValid() && event.Member.User.ID == me.ID {
				fmt.Printf("")
			}
		}
	}()

	voice.AddIntents(dc.handler.state)
}

* Command to disconnect
* Auto determine audio channel ID
* Other commands like skip, track, playlist
* Say what the album name is
* Allow one signal kill to exit cleanly
* -stream-from flag, defaults to stream, but you can also stream from disk (if the bot has access to same file system as the subsonic instance)

### Pre-requisites
* Install ffmpeg
* Create Discord bot
* Set environment variables
* Run or build and run


SUBSONIC_URL=http://localhost:4533
SUBSONIC_USER=admin
SUBSONIC_PASSWORD=admin
