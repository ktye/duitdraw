package duitdraw

import (
	"image"
)

// Mouse is the structure describing the current state of the mouse.
type Mouse struct {
	image.Point        // Location.
	Buttons     int    // Buttons; bit 0 is button 1, bit 1 is button 2, etc.
	Msec        uint32 // Time stamp in milliseconds.

}

type Mousectl struct {
	Mouse              // Store Mouse events here.
	C       chan Mouse // Channel of Mouse events.
	Resize  chan bool  // Each received value signals a window resize (see the display.Attach method).
	Display *Display
}

// Read returns the next mouse event.
func (mc *Mousectl) Read() Mouse {
	mc.Display.Flush()
	m := <-mc.C
	mc.Mouse = m
	return m
}
