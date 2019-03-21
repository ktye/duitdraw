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

// TODO: Mouse field is racy but okay.

// Mousectl holds the interface to receive mouse events.
// The Mousectl's Mouse is updated after send so it doesn't
// have the wrong value if the sending goroutine blocks during send.
// This means that programs should receive into Mousectl.Mouse
// if they want full synchrony.
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
	// Mouse field is racy. See Mousectl documentation.
	mc.Mouse = m
	return m
}
