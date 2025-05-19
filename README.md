# subcordant

## TODOs
* Testability - we have to do a delay
* Not sure I really need cancel funcs and some contexts now
* Testing errors and enqueing, one song, multiple songs
* Kick bot
* Properly handle errors
* Say what the album name is
* Support playlist
* DOn't emit that ffmpeg was cancelled
* Allow one signal kill to exit cleanly
* -stream-from flag, defaults to stream, but you can also stream from disk (if the bot has access to same file system as the subsonic instance)
* -stream-buffer flag, defaults to disk, but you can choose RAM, for named pipe option (will need other stuff set, make it linux/mac only at first, then add windows support)
* option to run from disk (current) or named pipe (need windows and linux versions, and set upper limit of pipe size via envionment variable, 100mb?)

## Scenarios to test
* Stream from subsonic is slower than reading the file from ffmpeg

### Pre-requisites
* Install ffmpeg
* Create Discord bot
