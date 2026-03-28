package main

import (
	"syscall"
	"unsafe"
)

type Config struct {
	long, recursive, all, reverse, time bool
}

/*
Required flags:
-l, -R, -a, -r, -t
*/

func parseFlags(args []string) ([]string, Config) {
	var paths []string
	var config Config

	for _, arg := range args {
		switch arg {
		case "-l":
			config.long = true
			//
		case "-R", "--recursive":
			config.recursive = true
			//
		case "-a", "--all":
			config.all = true
			//
		case "-r", "--reverse":
			config.reverse = true
			//
		case "-t": // different from --time ?!
			config.time = true
			//
		default:
			paths = append(paths, arg)
		}
	}

	if len(paths) == 0 { // default
		paths = append(paths, ".")
	}
	return paths, config
}

type WinSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getTerminalWidth() uint {
	ws := &WinSize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin), // Stdout ?
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	// OR:
	// if err != 0 {
	// 	return 80, fmt.Errorf("unable to get terminal size")
	// }
	return uint(ws.Col)
	// return int(ws.Col), nil
}
