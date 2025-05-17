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

type ExecCommander struct {
	stdout io.ReadCloser
	cmd    *exec.Cmd
}

func (ecmd *ExecCommander) Start(ctx context.Context, input io.ReadCloser,
	inputDestination string, cancelFunc context.CancelFunc) error {

	file, err := os.Create(inputDestination)
	if err != nil {
		input.Close()
		file.Close()
		cancelFunc()
		return fmt.Errorf("error creating file: %v", err)
	}

	// TODO - could this go in a separate module, or would it be too annoying?
	// Would want to pass context in so we could cancel and close stuff
	go func() {
		defer file.Close()
		defer input.Close()

		writer := bufio.NewWriter(file)

		ticker := time.NewTicker(flushFrequency * time.Millisecond)
		defer ticker.Stop()

		errorOccurred := false
		go func() {
			for range ticker.C {
				err := writer.Flush()
				if err != nil {
					log.Printf("\nERROR: Flush resulted in: %v", err)
					errorOccurred = true
					break
				}
			}
		}()

		n, err := io.Copy(writer, input)
		if err != nil {
			log.Printf("\nERROR: writing to file after %v bytes resulted in: %v", n, err)
			errorOccurred = true
		}

		if finalErr := writer.Flush(); finalErr != nil {
			log.Printf("\nERROR: final flush resulted in: %v", finalErr)
			errorOccurred = true
		}

		if errorOccurred {
			ecmd.stdout.Close()
			ecmd.cmd.Cancel()
			cancelFunc()
		}
	}()

	// TODO - fix this and use context cancel to wait until file is a certain size before proceeding
	// album ID fc9bf7bb3a0c6c4218112c72fedb0a29 shows this
	// Wait until the file is ready for streaming
	minSize := int64(1024 * 100) // Example: Wait for 100KB of data
	checkInterval := 50 * time.Millisecond
	if err := waitForFileReady(inputDestination, minSize, checkInterval); err != nil {
		return fmt.Errorf("file not ready for streaming: %w", err)
	}

	ecmd.cmd = exec.CommandContext(ctx,
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

	ecmd.cmd.Stderr = os.Stderr

	stdout, err := ecmd.cmd.StdoutPipe()
	if err != nil {
		input.Close()
		file.Close()
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
		file.Close()
		input.Close()
		cancelFunc()
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	go func() {
		<-ctx.Done()
		ecmd.stdout.Close()
		ecmd.cmd.Cancel()
		file.Close()
		input.Close()
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

func waitForFileReady(filePath string, minSize int64, checkInterval time.Duration) error {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		info, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue // File not created yet
			}
			return fmt.Errorf("failed to stat file: %w", err)
		}

		if info.Size() >= minSize {
			return nil // File is ready
		}
	}
	return nil
}
