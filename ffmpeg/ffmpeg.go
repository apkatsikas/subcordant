package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/apkatsikas/subcordant/constants"
	"github.com/diamondburned/oggreader"
)

var flushFrequency = 250 * time.Millisecond

type FfmpegCommander struct {
	stdout io.ReadCloser
	cmd    *exec.Cmd
}

func (fc *FfmpegCommander) Start(ctx context.Context, input io.ReadCloser, inputDestination string) error {
	file, err := os.Create(inputDestination)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	fc.cmd = exec.CommandContext(ctx,
		"ffmpeg", "-hide_banner", "-loglevel", "warning",
		"-threads", "1", // Single thread
		"-i", inputDestination,
		"-c:a", "libopus", // Codec
		"-b:a", "128k", // Bitrate
		"-frame_duration", strconv.Itoa(constants.FrameDuration),
		"-vbr", "off", // Disable variable bitrate
		"-f", "opus", // Output format
		"-", // Output to stdout
	)

	fc.cmd.Stderr = os.Stderr

	stdout, err := fc.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	fc.stdout = stdout

	if err := fc.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	go func() {
		defer file.Close()

		writer := bufio.NewWriter(file)

		ticker := time.NewTicker(flushFrequency * time.Millisecond)
		defer ticker.Stop()
		go func() {
			for range ticker.C {
				err := writer.Flush()
				if err != nil {
					log.Printf("\nERROR: Flush resulted in: %v", err)
				}
			}
		}()

		n, err := io.Copy(writer, input)
		if err != nil {
			log.Printf("\nERROR: writing to file after %v bytes resulted in: %v", n, err)
		}

		if err := writer.Flush(); err != nil {
			log.Printf("\nERROR: final flush resulted in: %v", err)
		}
	}()

	return nil
}

func (fc *FfmpegCommander) Stream(voice io.Writer) error {
	// TODO - close stdout?
	if err := oggreader.DecodeBuffered(voice, fc.stdout); err != nil {
		return fmt.Errorf("failed to decode ogg: %w", err)
	}

	if err := fc.cmd.Wait(); err != nil {
		return fmt.Errorf("failed to finish ffmpeg: %w", err)
	}
	return nil
}
