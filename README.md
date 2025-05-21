# subcordant

## TODOs
* log arikawa issue for disconnect
* set idle disconnect handler - after X amount of time without playing, DC and reset
* Don't emit that ffmpeg was cancelled
* Check all TODOs and make tickets
* Testing errors
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

	ch := make(chan *gateway.VoiceStateUpdateEvent)
	hand.state.AddHandler(ch)

	go func() {
		for event := range ch {
			log.Printf("ChannelID %v UserID %v Member.User.ID %v Session ID %v", event.ChannelID, event.UserID, event.Member.User.ID, event.SessionID)
		}
	}()

	Something like this - but how do we know its a disconnect?
* Command to disconnect
* Other commands like skip, track, playlist
* Say what the album name is
* Allow one signal kill to exit cleanly
* -stream-from flag, defaults to stream, but you can also stream from disk (if the bot has access to same file system as the subsonic instance)
* Auto determine audio channel ID

### Pre-requisites
* Install ffmpeg
* Create Discord bot
* Set environment variables
* Run or build and run


SUBSONIC_URL=http://localhost:4533
SUBSONIC_USER=admin
SUBSONIC_PASSWORD=admin
