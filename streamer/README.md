# streamer package

The `streamer` package converts an audio source (file or stream) into **Opus frames**
and exposes them via `voice.OpusFrameProvider`, which Discord pulls from during playback.

Playback is **pull-based**: Discord controls timing and requests frames as needed.
The streamer is responsible for producing frames.

---

## High-level flow

ffmpeg
└─ stdout (OGG/Opus)
└─ oggdemuxer
└─ frameQueue (io.Writer + OpusFrameProvider)
└─ Discord voice connection

---

## What’s happening

1. **ffmpeg** transcodes audio into Opus and writes an OGG stream to stdout.
2. **oggdemuxer** parses the OGG container and emits raw Opus packets.
3. **frameQueue** bridges transcoding and playback:
   - receives frames from the demuxer (push)
   - buffers them
   - supplies frames to Discord on demand (pull)

A single `frameQueue` instance sits between transcoding and playback.
It implements both `io.Writer` and `voice.OpusFrameProvider`, adapting
push-based transcoding to pull-based playback.

---

## `oggdemuxer.go`

Responsible for container parsing.

- Reads OGG pages from an `io.Reader`
- Reconstructs Opus packets from OGG segments
- Writes complete Opus packets to an `io.Writer`

This file exists because:

- ffmpeg outputs OGG/Opus, not raw frames
- Discord requires raw Opus frames, not OGG pages

The demuxer implements only what is required to extract Opus packets for streaming.

---

## `framequeue.go`

Owns frame buffering and playback coordination.

The frame queue plays three roles:

### 1. `io.Writer` (transcoding side)

Used by `oggdemuxer`.

- Receives complete Opus frames
- Copies each frame
- Pushes frames into a buffered channel

This side is push-based.

---

### 2. `voice.OpusFrameProvider` (playback side)

Used by DisGo’s audio sender.

- Implements `ProvideOpusFrame()`
- Blocks until a frame is available
- Returns `io.EOF` when the stream ends or is closed

This side is pull-based.

---

### 3. Buffer & lifecycle coordinator

- Buffered channel provides jitter tolerance between transcode and playback
- `closed` channel coordinates shutdown
- Safe to close from:
  - context cancellation
  - transcode completion
  - error paths

The queue exists solely to decouple transcode speed from playback timing.

---

## `streamer.go`

Owns the stream lifecycle.

Responsibilities:

- Spawn and manage the `ffmpeg` process
- Wire `stdout → oggdemuxer → frameQueue`
- Install the frame provider via a callback
- Handle cancellation, cleanup, and waiting for `ffmpeg` to exit

The streamer deliberately does not know about the Discord voice session.
The caller decides where and when the `OpusFrameProvider` is installed.

---

## Design characteristics

- Pull-based playback (Discord controls frame timing)
- No tickers or manual scheduling
- Single queue adapts push → pull
- Separation between:
  - process lifecycle (`Streamer`)
  - container parsing (`oggdemuxer`)
  - buffering & playback contract (`frameQueue`)

---

## Model in summary

> ffmpeg produces packets → demuxer extracts frames → queue buffers → Discord pulls
