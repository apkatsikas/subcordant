# subcordant

[![ subcordant Discord ][subcordant_img ]][subcordant ]

Subcordant is a Discord bot that streams music from your Subsonic-API compatible server.

[subcordant]: https://discord.gg/db4HrbjMSt
[subcordant_img]: https://img.shields.io/badge/subcordant-Discord-%237289da?style=flat-square

## Available Commands

- `/play` - takes a `subsonicid` parameter. Connects the bot to the voice channel currently occupied by the command issuer, if it is not yet connected. Enqueues all tracks from the specified album, playlist or track and initates playback, if not already playing. If the bot is already present in a different voice channel, playback will stop, the current playlist will be cleared and the bot will join the new channel
- `/albumtrack` - takes an `albumid` parameter and `tracknumber` parameter. Behaves the same as the `/play` command, but instead targets a single track from the specified album/track number combination
- `/trackname` - takes a `trackname` parameter. Behaves the same as the `albumtrack` command, but instead takes the track name
- `/albumname` - takes an `albumname` parameter. Behaves the same as the `albumtrack` command, but instead takes the album name
- `/clear` - clears the playlist and stops playback
- `/disconnect` - disconnects the subcordant bot, stopping playback and clearing the playlist
- `/skip` - skips the currently playing song
- `/help` - describes all commands

## Downloading

Download a binary of Subcordant from the releases page for your desired platform and extract.

Example on \*nix:

`tar -xzf /path/subcordant-v1.0.0-linux-amd64.tar.gz -C  /path/`

## Building

Run `make build`.

## Pre-requisites

- [FFmpeg](https://ffmpeg.org/) must be installed and available on your path
- Create a [Discord bot](docs/bot.md)

## Running

Run the executable, specifying the following environment variables:

- SUBSONIC_URL
- SUBSONIC_USER
- SUBSONIC_PASSWORD
- DISCORD_BOT_TOKEN

### Flags

The following flags are available:

- `streamFrom` - valid values are `stream` or `file`, defaults to `stream` if not specified. Subcordant will stream to the voice chat using:
  - `stream` - the Subsonic `/stream` endpoint
  - `file` - the file path returned from the `path` field on the song. This mode will only work if subcordant has access to the audio library folder, and that the API returns the actual absolute path of the song. For Navidrome, an admin can toggle to allow the subcordant player to `Report Real Path`, or use the [Subsonic.DefaultReportRealPath](https://www.navidrome.org/docs/usage/configuration-options/#:~:text=subsonic.defaultreportrealpath) configuration option. Consult your server's documentation for more information.
- `idleDisconnectTimeout` - duration in minutes after which bot will disconnect when no music is playing. Defaults to 5 if not specified.

## Installing as a systemd unit

These instructions are for installing subcordant as a [systemd unit](https://www.freedesktop.org/software/systemd/man/latest/systemd.unit.html) on Linux. This enables subcordant to run on machine startup.

0. Run `make build` or download a binary from the releases page
1. Create a directory at `/opt/subcordant`
2. Add your subcordant executable to this path
3. Run `chmod +x /opt/subcordant/subcordant`
4. Run `groupadd subcordant`
5. Run `useradd -m -g subcordant subcordant`
6. Run `chown subcordant:subcordant /opt/subcordant`
7. Run `vim /etc/systemd/system/subcordant.service` and copy the contents of [docs/subcordant.service](docs/subcordant.service) into the file, replacing all the environment variable instances of `foobar` with your values.
8. Run `systemctl daemon-reload`
9. Run `systemctl start subcordant.service`
10. Run `systemctl status subcordant.service`

## Contributing

Please see the Issues associated with this repo to help contribute.

### Testing

Run `make run-tests`.

Subcordant makes use of the following:

- [Ginkgo](https://github.com/onsi/ginkgo) - testing framework
- [Mockery](https://vektra.github.io/mockery/latest/) - for generating mocks from interfaces
- [Testify](https://github.com/stretchr/testify?tab=readme-ov-file#mock-package) - for the mock package

dont force libdave build on push check
use matrix of releases
update install docs
make release binaries work on every Linux distro automatically
Environment=LD_LIBRARY_PATH=/opt/subcordant vs Build with rpath:

-Wl,-rpath,'$ORIGIN'
