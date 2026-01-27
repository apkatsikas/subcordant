package streamer

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strconv"

	"github.com/apkatsikas/subcordant/constants"
	"github.com/disgoorg/disgo/voice"
)

const frameBuffer = 100

type Streamer struct {
	stdout io.ReadCloser
	cmd    *exec.Cmd
}

func (s *Streamer) PrepStreamFromStream(inputUrl *url.URL) error {
	return s.prepStream(true, inputUrl.String())
}

func (s *Streamer) PrepStreamFromFile(inputPath string) error {
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf(
			"failed to prepare stream with file %v, error was %w", inputPath, err)
	}
	return s.prepStream(false, inputPath)
}

func (s *Streamer) prepStream(streamFromStream bool, inputString string) error {
	args := getArgs(streamFromStream, inputString)
	s.cmd = exec.CommandContext(context.Background(),
		"ffmpeg", args...,
	)

	// Enable this for debugging ffmpeg issues
	//s.cmd.Stderr = os.Stderr

	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		stdout.Close()
		s.stdout.Close()
		safeCancel(s.cmd.Cancel)
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	s.stdout = stdout

	if err := s.cmd.Start(); err != nil {
		stdout.Close()
		s.stdout.Close()
		safeCancel(s.cmd.Cancel)
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	return nil
}

func (s *Streamer) Stream(
	ctx context.Context,
	setFrameProvider func(voice.OpusFrameProvider),
) error {
	frameQueue := newOpusFrameQueue(frameBuffer)
	setFrameProvider(frameQueue)

	decodingDone := make(chan error, 1)

	go func() {
		defer close(frameQueue.frames)

		if err := demuxOpusFromOGGBuffered(frameQueue, s.stdout); err != nil {
			decodingDone <- fmt.Errorf("failed to decode ogg: %w", err)
			return
		}
		decodingDone <- nil
	}()

	cleanup := func() {
		s.stdout.Close()
		frameQueue.Close()
		if s.cmd != nil {
			safeCancel(s.cmd.Cancel)
		}
	}

	select {
	case <-ctx.Done():
		cleanup()
		return nil

	case err := <-decodingDone:
		if err != nil {
			cleanup()
			return err
		}
	}

	return nil
}

func safeCancel(cancel func() error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic during cmd.Cancel(): %v", r)
		}
	}()
	cancel()
}

func getArgs(streamFromStream bool, inputString string) []string {
	args := preInputArgs()
	if streamFromStream {
		args = append(args, reconnectArgs()...)
	}
	args = append(args, inputAndPostArgs(inputString)...)
	return args
}

func preInputArgs() []string {
	return []string{
		"-hide_banner",
		"-loglevel", "warning",
		"-threads", "1",
	}
}

func reconnectArgs() []string {
	return []string{
		"-reconnect", "1", // These flags keep the stream running
		"-reconnect_streamed", "1", // by reconnecting after being disconnected
		"-reconnect_delay_max", "5", // from subsonic
	}
}

func inputAndPostArgs(inputString string) []string {
	return []string{
		"-i", inputString,
		"-c:a", "libopus",
		"-b:a", "128k",
		"-frame_duration", strconv.Itoa(constants.FrameDuration),
		"-vbr", "off",
		"-f", "opus",
		"-", // Output to stdout
	}
}
