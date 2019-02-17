// +build !windows
// +build !dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris

package duitdraw

import (
	"errors"
	"fmt"
	"image"
)

func moveTo(p image.Point) error {
	return fmt.Errorf("moveTo: TODO for this os")
}

func setCursor(c *Cursor) error {
	if c != nil {
		return errors.New("duitdraw: SetCursor is not implemented")
	}
	return nil
}
