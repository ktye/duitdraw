// +build windows darwin

package duitdraw

import (
	"github.com/atotto/clipboard"
)

func (d *Display) readSnarf(buf []byte) (int, int, error) {
	s, err := clipboard.ReadAll()
	if err != nil {
		return 0, 0, err
	}
	n := copy(buf, s)
	if n < len(s) {
		return n, len(s), errShortSnarfBuffer
	}
	return n, n, nil
}

func (d *Display) writeSnarf(data []byte) error {
	return clipboard.WriteAll(string(data))
}
