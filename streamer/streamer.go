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
	"github.com/diamondburned/oggreader"
)

type Streamer struct {
	stdout io.ReadCloser
	cmd    *exec.Cmd
}

func (s *Streamer) PrepStream(ctx context.Context, inputUrl *url.URL,
	cancelFunc context.CancelFunc) error {

	s.cmd = exec.CommandContext(ctx,
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

	s.cmd.Stderr = os.Stderr

	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		stdout.Close()
		s.stdout.Close()
		s.cmd.Cancel()
		cancelFunc()
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	s.stdout = stdout

	if err := s.cmd.Start(); err != nil {
		s.stdout.Close()
		s.cmd.Cancel()
		cancelFunc()
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	go func() {
		<-ctx.Done()
		s.stdout.Close()
		s.cmd.Cancel()
		log.Println("cancelling ffmpeg as context was cancelled")
	}()

	return nil
}

func (s *Streamer) Stream(voice io.Writer, cancelFunc context.CancelFunc) error {
	defer s.stdout.Close()
	if err := oggreader.DecodeBuffered(voice, s.stdout); err != nil {
		cancelFunc()
		s.stdout.Close()
		return fmt.Errorf("failed to decode ogg: %w", err)
	}

	if err := s.cmd.Wait(); err != nil {
		cancelFunc()
		s.stdout.Close()
		return fmt.Errorf("failed to finish ffmpeg: %w", err)
	}
	return nil
}
