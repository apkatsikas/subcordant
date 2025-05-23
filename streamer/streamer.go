package streamer

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strconv"

	"github.com/apkatsikas/subcordant/constants"
	"github.com/diamondburned/oggreader"
)

type Streamer struct {
	stdout io.ReadCloser
	cmd    *exec.Cmd
}

func (s *Streamer) PrepStream(inputUrl *url.URL) error {
	s.cmd = exec.CommandContext(context.Background(),
		"ffmpeg",
		"-hide_banner",
		"-loglevel", "warning",
		"-reconnect", "1", // These flags keep the stream running
		"-reconnect_streamed", "1", // by reconnecting after being disconnected
		"-reconnect_delay_max", "5", // from subsonic
		"-threads", "1",
		"-i", inputUrl.String(),
		"-c:a", "libopus",
		"-b:a", "128k",
		"-frame_duration", strconv.Itoa(constants.FrameDuration),
		"-vbr", "off",
		"-f", "opus",
		"-", // Output to stdout
	)

	// TODO - remove when we are publishing
	s.cmd.Stderr = os.Stderr

	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		stdout.Close()
		s.stdout.Close()
		s.cmd.Cancel()
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	s.stdout = stdout

	if err := s.cmd.Start(); err != nil {
		stdout.Close()
		s.stdout.Close()
		s.cmd.Cancel()
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	return nil
}

func (s *Streamer) Stream(ctx context.Context, voice io.Writer) error {
	defer s.stdout.Close()

	decodingDone := make(chan error, 1)
	go func() {
		if err := oggreader.DecodeBuffered(voice, s.stdout); err != nil {
			decodingDone <- fmt.Errorf("failed to decode ogg: %w", err)
			return
		}
		decodingDone <- nil
	}()

	select {
	case <-ctx.Done():
		s.stdout.Close()
		s.cmd.Cancel()
		return nil
	case err := <-decodingDone:
		if err != nil {
			s.stdout.Close()
			s.cmd.Cancel()
			return err
		}
	}

	if err := s.cmd.Wait(); err != nil {
		s.stdout.Close()
		s.cmd.Cancel()
		return fmt.Errorf("failed to finish ffmpeg: %w", err)
	}

	return nil
}
