//go:build windows
// +build windows

package xutil

import errors

// TermSize get terminal size
func TermSize() (int, int, error) {
	return 0, 0, errors.New("not support")
}
