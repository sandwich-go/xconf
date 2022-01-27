//go:build windows
// +build windows

package xutil

// TermSize get terminal size
func TermSize() (int, int, error) {
	return 0, 0, errors.New("not support")
}
