# subcordant

## TODOs
* Testability - we have to do a delay
* Testing errors and enqueing, one song, multiple songs
* Kick bot
* Properly handle errors
* Say what the album name is
* Support playlist
* DOn't emit that ffmpeg was cancelled
* Allow one signal kill to exit cleanly
* option to run from disk (current) or named pipe (need windows and linux versions, and set upper limit of pipe size via envionment variable, 100mb?)

## Scenarios to test
* Stream from subsonic is slower than reading the file from ffmpeg

### Pre-requisites
* Install ffmpeg
* Create Discord bot
