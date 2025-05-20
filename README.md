# subcordant

## TODOs
* Don't emit that ffmpeg was cancelled
* Test for when play command is recieved while play is already running (enqueue but dont play) via go sr.Play() then normal Play - expect only x calls to deps
* Check all TODOs and make tickets
* Testing errors
* Kick bot
* Say what the album name is
* Support playlist
* Skip track
* Allow one signal kill to exit cleanly
* -stream-from flag, defaults to stream, but you can also stream from disk (if the bot has access to same file system as the subsonic instance)
* Auto determine audio channel ID

## Scenarios to test
* Stream from subsonic is slower than reading the file from ffmpeg

### Pre-requisites
* Install ffmpeg
* Create Discord bot
* Set environment variables
* Run or build and run
