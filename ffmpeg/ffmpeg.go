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

func (fc *FfmpegCommander) Start(ctx context.Context, input io.ReadCloser,
	inputDestination string, cancelFunc context.CancelFunc) error {

	file, err := os.Create(inputDestination)
	if err != nil {
		input.Close()
		file.Close()
		cancelFunc()
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
		input.Close()
		file.Close()
		stdout.Close()
		fc.stdout.Close()
		fc.cmd.Cancel()
		cancelFunc()
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	fc.stdout = stdout

	if err := fc.cmd.Start(); err != nil {
		fc.stdout.Close()
		fc.cmd.Cancel()
		file.Close()
		input.Close()
		cancelFunc()
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	go func() {
		<-ctx.Done()
		fc.stdout.Close()
		fc.cmd.Cancel()
		file.Close()
		input.Close()
		log.Println("cancelling ffmpeg as context was cancelled")
	}()

	// TODO - could this go in a separate module, or would it be too annoying?
	// Would want to pass context in so we could cancel and close stuff
	go func() {
		defer file.Close()
		defer input.Close()

		writer := bufio.NewWriter(file)

		ticker := time.NewTicker(flushFrequency * time.Millisecond)
		defer ticker.Stop()
		go func() {
			for range ticker.C {
				err := writer.Flush()
				if err != nil {
					file.Close()
					input.Close()
					fc.stdout.Close()
					fc.cmd.Cancel()
					cancelFunc()
					log.Printf("\nERROR: Flush resulted in: %v", err)
					return
				}
			}
		}()

		n, err := io.Copy(writer, input)
		if err != nil {
			file.Close()
			input.Close()
			fc.stdout.Close()
			fc.cmd.Cancel()
			cancelFunc()
			log.Printf("\nERROR: writing to file after %v bytes resulted in: %v", n, err)
			return
		}

		if err := writer.Flush(); err != nil {
			file.Close()
			input.Close()
			fc.stdout.Close()
			fc.cmd.Cancel()
			cancelFunc()
			log.Printf("\nERROR: final flush resulted in: %v", err)
			return
		}
	}()

	return nil
}

func (fc *FfmpegCommander) Stream(voice io.Writer, cancelFunc context.CancelFunc) error {
	defer fc.stdout.Close()
	if err := oggreader.DecodeBuffered(voice, fc.stdout); err != nil {
		cancelFunc()
		fc.stdout.Close()
		return fmt.Errorf("failed to decode ogg: %w", err)
	}

	if err := fc.cmd.Wait(); err != nil {
		cancelFunc()
		fc.stdout.Close()
		return fmt.Errorf("failed to finish ffmpeg: %w", err)
	}
	return nil
}
