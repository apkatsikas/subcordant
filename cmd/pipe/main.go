package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func main() {
	const pipePath = "/tmp/music_pipe"
	const bufferSize = 100 * 1024 * 1024 // 100 MB

	// Create a named pipe if it doesn't exist
	if _, err := os.Stat(pipePath); os.IsNotExist(err) {
		if err := syscall.Mkfifo(pipePath, 0666); err != nil {
			log.Fatalf("Failed to create named pipe: %v", err)
		}
		fmt.Println("Named pipe created.")
	} else {
		fmt.Println("Named pipe already exists.")
	}

	// Open the named pipe
	file, err := os.OpenFile(pipePath, os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Failed to open named pipe: %v", err)
	}
	defer file.Close()

	fmt.Println("Named pipe opened.")

	// Retrieve the file descriptor for the pipe
	fd := file.Fd()

	// Adjust the buffer size
	if _, _, errno := syscall.Syscall(syscall.SYS_FCNTL, fd, syscall.F_SETPIPE_SZ, uintptr(bufferSize)); errno != 0 {
		log.Fatalf("Failed to set pipe buffer size: %v", errno)
	}
	fmt.Printf("Pipe buffer size adjusted to %d bytes.\n", bufferSize)

	// Hold the pipe open for testing
	fmt.Println("Named pipe is ready. Press Enter to exit.")
	fmt.Scanln()
}
