package ffmpeg

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

type FfmpegCommander struct {
	stdout io.ReadCloser
	cmd    *exec.Cmd
	stdin  io.WriteCloser // Optional: Keep track of stdin to allow cleanup or testing
}

func (fc *FfmpegCommander) Start(ctx context.Context, input io.ReadCloser) error {
	// FFmpeg command with "-" as input (reads from stdin)
	fc.cmd = exec.CommandContext(ctx,
		"ffmpeg", "-hide_banner", "-loglevel", "warning",
		"-threads", "1", // Single thread
		"-i", "-", // Read from standard input
		"-c:a", "libopus", // Codec
		"-b:a", "128k", // Bitrate
		"-frame_duration", strconv.Itoa(constants.FrameDuration),
		"-vbr", "off", // Disable variable bitrate
		"-f", "opus", // Output format
		"-", // Output to stdout
	)

	fc.cmd.Stderr = os.Stderr

	// Set stdin to the input stream
	stdin, err := fc.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}
	fc.stdin = stdin

	// Set stdout to capture the FFmpeg output
	stdout, err := fc.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	fc.stdout = stdout

	// Start the FFmpeg process
	if err := fc.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Pipe the input stream to FFmpeg's stdin in a goroutine
	go func() {
		defer stdin.Close()
		n, err := io.Copy(stdin, input)
		if err != nil {
			log.Printf("Error streaming input after %v bytes: %v", n, err)
		} else {
			log.Printf("Streamed %d bytes to FFmpeg", n)
		}
	}()

	return nil
}

// Stream processes FFmpeg's stdout and writes to the provided io.Writer.
func (fc *FfmpegCommander) Stream(output io.Writer) error {
	// Decode and process the FFmpeg output
	if err := oggreader.DecodeBuffered(output, fc.stdout); err != nil {
		return fmt.Errorf("failed to decode ogg: %w", err)
	}

	// Wait for FFmpeg to finish
	if err := fc.cmd.Wait(); err != nil {
		return fmt.Errorf("failed to finish ffmpeg: %w", err)
	}
	return nil
}
