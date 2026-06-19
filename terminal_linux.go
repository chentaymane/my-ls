//go:build linux

package main

import (
	"syscall"
	"unsafe"
)

type WinSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getTerminalWidth() uint {
	ws := &WinSize{}
	retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))
	if int(retCode) == -1 || ws.Col == 0 {
		return 80
	}
	return uint(ws.Col)
}
