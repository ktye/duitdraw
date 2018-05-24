package draw

import (
	"image"
	"time"
)

// Mouse is the structure describing the current state of the mouse.
type Mouse struct {
	image.Point        // Location.
	Buttons     int    // Buttons; bit 0 is button 1, bit 1 is button 2, etc.
	Msec        uint32 // Time stamp in milliseconds.

}

type Mousectl struct {
	Mouse             // Store Mouse events here.
	C      chan Mouse // Channel of Mouse events.
	Resize chan bool  // Each received value signals a window resize (see the display.Attach method).
	last   time.Time  // Time of last update.
}
