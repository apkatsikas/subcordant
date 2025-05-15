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

type FfmpegCommander struct {
	stdout io.ReadCloser
	cmd    *exec.Cmd
}

func (fc *FfmpegCommander) Start(ctx context.Context, input io.ReadCloser, inputDestination string) error {
	// Open the file for writing
	file, err := os.Create(inputDestination)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}

	// FFmpeg command with "-" as input (reads from stdin)
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

	go func() {
		defer file.Close()

		// Create a buffered writer
		writer := bufio.NewWriter(file)

		// Periodically flush the buffer
		ticker := time.NewTicker(250 * time.Millisecond) // Adjust flush interval as needed
		defer ticker.Stop()

		// Goroutine to flush the buffer periodically
		go func() {
			for range ticker.C {
				err := writer.Flush()
				if err != nil {
					log.Printf("Error flushing file: %v", err)
				}
			}
		}()

		// Copy data from the input to the file
		n, err := io.Copy(writer, input)
		if err != nil {
			log.Printf("Error writing to file after %v bytes: %v", n, err)
		} else {
			log.Printf("Wrote %d bytes to file", n)
		}

		// Final flush to ensure all data is written
		if err := writer.Flush(); err != nil {
			log.Printf("Final flush error: %v", err)
		}
	}()

	return nil
}

// Stream processes FFmpeg's stdout and writes to the provided io.Writer.
func (fc *FfmpegCommander) Stream(voice io.Writer) error {
	// TODO - close stdout?
	// Decode and process the FFmpeg output
	if err := oggreader.DecodeBuffered(voice, fc.stdout); err != nil {
		return fmt.Errorf("failed to decode ogg: %w", err)
	}

	// Wait for FFmpeg to finish
	if err := fc.cmd.Wait(); err != nil {
		return fmt.Errorf("failed to finish ffmpeg: %w", err)
	}
	return nil
}
