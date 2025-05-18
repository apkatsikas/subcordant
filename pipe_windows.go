//go:build windows

package main

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func createNamedPipe(pipePath string, bufferSize int) error {
	// Convert pipe path to Windows-style
	pipePath = `\\.\pipe\` + pipePath

	// Create the named pipe
	handle, err := windows.CreateNamedPipe(
		windows.StringToUTF16Ptr(pipePath),
		windows.PIPE_ACCESS_DUPLEX,               // Read/Write access
		windows.PIPE_TYPE_BYTE|windows.PIPE_WAIT, // Byte-oriented, blocking
		1,                  // Max instances
		uint32(bufferSize), // Output buffer size
		uint32(bufferSize), // Input buffer size
		0,                  // Default timeout
		nil,                // Default security attributes
	)
	if err != nil {
		return fmt.Errorf("failed to create named pipe: %w", err)
	}

	// Ensure the handle is closed properly
	defer windows.CloseHandle(handle)

	// The pipe is created; it can now be used for reading/writing.
	return nil
}
