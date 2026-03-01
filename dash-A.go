package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	reset = "\033[0m"

	// Foreground colors
	blue    = "\033[34m" // Directories
	green   = "\033[32m" // Executables
	cyan    = "\033[36m" // Symlinks
	yellow  = "\033[33m" // Device files
	magenta = "\033[35m" // Images/media
	red     = "\033[31m" // Archives

	// Background colors for special files
	blackOnRed    = "\033[30;41m" // Broken symlink
	blackOnYellow = "\033[30;43m" // Setuid files
	whiteOnRed    = "\033[37;41m" // Missing file
	blackOnGreen  = "\033[30;42m" // Setgid files
)

// Simple bubble sort (since sort package not allowed)
func sortNames(names []string) {
	n := len(names)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if names[j] > names[j+1] {
				names[j], names[j+1] = names[j+1], names[j]
			}
		}
	}
}

func needsQuotes(name string) bool {
	return strings.ContainsAny(name, " \t\n")
}

func formatName(name string) string {
	if needsQuotes(name) {
		return "'" + name + "'"
	}
	return name
}

func getFileColor(path string, name string) string {
	info, err := os.Lstat(path)
	if err != nil {
		return reset
	}

	mode := info.Mode()

	// Check for symlink
	if mode&os.ModeSymlink != 0 {
		// Check if symlink is broken
		if _, err := os.Stat(path); err != nil {
			return blackOnRed // Broken symlink
		}
		return cyan // Valid symlink
	}

	// Check for directory
	if mode.IsDir() {
		// Check for sticky bit and other writable
		if mode&os.ModeSticky != 0 && mode&0o002 != 0 {
			return "\033[30;42m" // Black text on green (sticky + other-writable)
		}
		// Check for other writable without sticky
		if mode&0o002 != 0 {
			return "\033[34;42m" // Blue text on green (other-writable)
		}
		// Check for sticky bit
		if mode&os.ModeSticky != 0 {
			return "\033[37;44m" // White text on blue (sticky)
		}
		return blue // Regular directory
	}

	// Check for setuid
	if mode&os.ModeSetuid != 0 {
		return blackOnYellow
	}

	// Check for setgid
	if mode&os.ModeSetgid != 0 {
		return blackOnGreen
	}

	// Check for executable
	if mode&0o111 != 0 {
		return green
	}

	// Check for special files by extension

	// Device files
	if mode&os.ModeDevice != 0 {
		return yellow
	}

	// Named pipes
	if mode&os.ModeNamedPipe != 0 {
		return yellow
	}

	// Socket
	if mode&os.ModeSocket != 0 {
		return magenta
	}

	// Regular file
	return reset
}

func A(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	var names []string

	for _, file := range files {
		names = append(names, file.Name())
	}

	// Add . and ..
	names = append(names, ".", "..")

	// Sort like original ls
	sortNames(names)

	for _, name := range names {
		fullPath := path + "/" + name
		color := getFileColor(fullPath, name)

		// Add leading space like real ls does
		fmt.Println(" " + color + formatName(name) + reset)
	}
}
