# subcordant

TODO
ffmpeg will be a dependency, as we need it to create opus stream:
https://github.com/Gimzie/submeister/blob/015218a906599f9abe208f7cd6685b8209147f4d/player.py#L69

where we provide the subsonic stream to the audio source for ffmpeg to stream (streaming the stream)
compare to golang stream:
https://github.com/diamondburned/arikawa/blob/8a78eb04430cfd0f4997c8bf206cf36c0c2e604d/0-examples/voice/main.go#L75
