package ecmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/apkatsikas/subcordant/constants"
	"github.com/diamondburned/oggreader"
)

// TODO - we could make this class take dependencies to the files system and os/exec
// and rename it to Streamer. then we can run tests over it specifically

type ExecCommander struct {
	stdout io.ReadCloser
	cmd    *exec.Cmd
}

func (ecmd *ExecCommander) Start(ctx context.Context, input string,
	cancelFunc context.CancelFunc) error {

	ecmd.cmd = exec.CommandContext(ctx,
		"ffmpeg", "-hide_banner", "-loglevel", "warning",
		"-reconnect", "1", "-reconnect_streamed", "1", "-reconnect_delay_max", "5",
		"-threads", "1", // Single thread
		"-i", input,
		"-c:a", "libopus", // Codec
		"-b:a", "128k", // Bitrate
		"-frame_duration", strconv.Itoa(constants.FrameDuration),
		"-vbr", "off", // Disable variable bitrate
		"-f", "opus", // Output format
		"-", // Output to stdout
	)

	ecmd.cmd.Stderr = os.Stderr

	stdout, err := ecmd.cmd.StdoutPipe()
	if err != nil {
		stdout.Close()
		ecmd.stdout.Close()
		ecmd.cmd.Cancel()
		cancelFunc()
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	ecmd.stdout = stdout

	if err := ecmd.cmd.Start(); err != nil {
		ecmd.stdout.Close()
		ecmd.cmd.Cancel()
		cancelFunc()
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	go func() {
		<-ctx.Done()
		ecmd.stdout.Close()
		ecmd.cmd.Cancel()
		log.Println("cancelling ffmpeg as context was cancelled")
	}()

	return nil
}

func (ecmd *ExecCommander) Stream(voice io.Writer, cancelFunc context.CancelFunc) error {
	defer ecmd.stdout.Close()
	if err := oggreader.DecodeBuffered(voice, ecmd.stdout); err != nil {
		cancelFunc()
		ecmd.stdout.Close()
		return fmt.Errorf("failed to decode ogg: %w", err)
	}

	if err := ecmd.cmd.Wait(); err != nil {
		cancelFunc()
		ecmd.stdout.Close()
		return fmt.Errorf("failed to finish ffmpeg: %w", err)
	}
	return nil
}
