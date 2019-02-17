package duitdraw

import (
	"errors"
	"fmt"
	"image"
	"os"

	"github.com/as/cursor"
	"github.com/as/ms/win"
)

// This sets the mouse cursor relative to the current window.
//
// TODO: There seems to be an offset, maybe it's related to the window borders.
// You can calculate your offset by pressing the tab key and clicking with the mouse.
// Log coordinates of MouseMove and mouse.Event (in display.go) and substract the results.
//
// Source: github.com/as/a/mouse_windows.go

var (
	winfd        win.Window
	winfderr     error
	cursorOffset = image.Point{4, 4}
)

func tryWindow() {
	if winfd != 0 && winfderr == nil {
		return
	}
	winfd, winfderr = win.Open(os.Getpid())
}

func moveTo(pt image.Point) error {
	tryWindow()
	abs, err := winfd.Client()
	if err != nil {
		winfderr = err
		return err
	}
	pt = pt.Add(abs.Min).Sub(cursorOffset)
	if cursor.MoveTo(pt) == false {
		return fmt.Errorf("move cursor failed")
	}
	return nil
}

func setCursor(c *Cursor) error {
	if c != nil {
		return errors.New("duitdraw: SetCursor is not implemented")
	}
	return nil
}
