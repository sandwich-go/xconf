//go:build !windows
// +build !windows

package xutil

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

type size struct {
	rows uint16
	cols uint16
}

// TermSize get terminal size
func TermSize() (int, int, error) {

	var sz size
	var fd uintptr

	if runtime.GOOS == "windows" {
		if fh, err := syscall.Open("CONOUT$", syscall.O_RDWR, 0); err != nil {
			return int(0), int(0), err
		} else {
			fd = uintptr(fh)
		}
	} else {
		if fp, err := os.OpenFile("/dev/tty", syscall.O_WRONLY, 0); err != nil {
			return int(0), int(0), err
		} else {
			fd = fp.Fd()
		}
	}

	if _, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&sz))); err != 0 {
		return int(0), int(0), err
	}

	return int(sz.cols), int(sz.rows), nil
}
