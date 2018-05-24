package draw

import (
	"image"
	"time"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

// EventLoop is the event loop for a single window.
func (d *Display) eventLoop(errch chan<- error) {
	w := d.window
	b := d.buffer
	var err error
	for {
		switch e := w.NextEvent().(type) {
		case lifecycle.Event:
			// fmt.Println(e)
			// TODO: Closing a single windows works, but closing
			// a child window leaves the window hanging.
			if e.To == lifecycle.StageDead {
				return
			}

		case paint.Event:
			d.ScreenImage.Lock()
			w.Upload(image.Point{}, b, b.Bounds())
			w.Publish()
			d.ScreenImage.Unlock()

		case size.Event:
			d.ScreenImage.Lock()
			if b != nil {
				b.Release()
			}
			b, err = mainScreen.NewBuffer(e.Size())
			if err != nil {
				errch <- err
				d.ScreenImage.Lock()
				return
			}
			d.buffer = b
			d.ScreenImage.m = b.RGBA()
			d.ScreenImage.R = b.Bounds()
			d.ScreenImage.Unlock()
			d.mouse.Resize <- true

		case mouse.Event:
			// Mouse.Buttons stores a bitmask for each button state.
			// On the other side a mouse.Event arrives, if anything changes.
			if e.Button > 0 { // TODO: wheel is < 0
				if e.Direction == mouse.DirPress {
					// Uncomment for cursorOffset calibration:
					// fmt.Printf("shiny: mouse click: %f %f\n", e.X, e.Y)
					d.mouse.Buttons ^= 1 << uint(e.Button-1)
				} else if e.Direction == mouse.DirRelease {
					d.mouse.Buttons &= ^(1 << uint(e.Button-1))
				}
			}
			d.mouse.Point.X = int(e.X)
			d.mouse.Point.Y = int(e.Y)
			t := time.Now()
			d.mouse.Msec = uint32(t.Sub(d.mouse.last) * time.Millisecond)
			d.mouse.last = t
			d.mouse.C <- d.mouse.Mouse

		case key.Event:
			// We forward the event for key presses and subsequent events
			// if the key remains down, but not for releases.
			var sendKey rune = -1
			if r := e.Rune; e.Direction != key.DirRelease {
				if r != -1 {
					if e.Modifiers&key.ModControl != 0 {
						//r += KeyCmd
						// Why changes the 'a' key from 0x61 to 0x1 when Cntrl
						// is pressed? This happens on windows.
						// Is this a bug in shiny or expected?
						r += KeyCmd + 'a' - 1
					}
					sendKey = r
				} else {
					if r, ok := keymap[e.Code]; ok {
						sendKey = r
					}
				}

			}
			if sendKey != -1 {
				// Shiny sends \r on Enter, duit expects \n.
				if sendKey == '\r' {
					sendKey = '\n'
				}
				// fmt.Printf("shiny: key: %x %v\n", sendKey, e)
				d.keyboard.C <- sendKey
			}

			// TODO: what about Shift-KeyLeft/Right
			// to mark text? This seems to be unsupported in duit right now.

		case error:
			errch <- e

		}
	}
}
