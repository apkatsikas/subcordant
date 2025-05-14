package ffmpeg

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"

	"github.com/apkatsikas/subcordant/constants"
	"github.com/diamondburned/oggreader"
)

type FfmpegCommander struct {
	stdout io.ReadCloser
	cmd    *exec.Cmd
}

func (fc *FfmpegCommander) Start(ctx context.Context, file string) error {
	fc.cmd = exec.CommandContext(ctx,
		"ffmpeg", "-hide_banner", "-loglevel", "error",
		// Streaming is slow, so a single thread is all we need.
		"-threads", "1",
		// Input file.
		"-i", file,
		// Output format; leave as "libopus".
		"-c:a", "libopus",
		// Bitrate in kilobits.
		"-b:a", "128k",
		// Frame duration should be the same as what's given into
		// udp.DialFuncWithFrequency.
		"-frame_duration", strconv.Itoa(constants.FrameDuration),
		// Disable variable bitrate to keep packet sizes consistent. This is
		// optional.
		"-vbr", "off",
		// Output format, which is opus, so we need to unwrap the opus file.
		"-f", "opus",
		"-",
	)

	fc.cmd.Stderr = os.Stderr

	stdout, err := fc.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	fc.stdout = stdout

	// FFmpeg will wait until we start consuming the stream to process further.
	if err := fc.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}
	return nil
}

func (fc *FfmpegCommander) Stream(voice io.Writer) error {
	if err := oggreader.DecodeBuffered(voice, fc.stdout); err != nil {
		return fmt.Errorf("failed to decode ogg: %w", err)
	}

	if err := fc.cmd.Wait(); err != nil {
		return fmt.Errorf("failed to finish ffmpeg: %w", err)
	}
	return nil
}
