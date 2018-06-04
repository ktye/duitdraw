package draw

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/atotto/clipboard"

	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
)

// Refresh is ignored in the implementation, but duit needs it.
const Refmesg = 2

// DefaultDPI is the initial DPI setting for a new display.
// TODO: should we get DPI settings from the screen?
// Currently there is no interface in shiny.
const DefaultDPI = 100

const DefaultFontSize = 10

// Display stores the information for a single window, that is returned to duit.
// Duit requestes a Display by calling Init for each window.
type Display struct {
	DPI         int
	ScreenImage *Image
	DefaultFont *Font
	Black       *Image // Pre-allocated color.
	White       *Image
	mouse       Mousectl
	keyboard    Keyboardctl
	window      screen.Window
	buffer      screen.Buffer
}

// AllocImage allocates a new Image on display d. The arguments are:
// - the rectangle representing the size
// - the pixel descriptor: RGBA32 etc.
// - whether the image is to be replicated (tiled)
// - the starting background color for the image
//
// Duit calls AllocImage to allocate colors for a single pixel rectange with repl = true.
// We return a uniform image instead.
func (d *Display) AllocImage(r image.Rectangle, pix Pix, repl bool, val Color) (*Image, error) {
	c := color.RGBA{
		R: uint8(val >> 24),
		G: uint8(val >> 16),
		B: uint8(val >> 8),
		A: uint8(val),
	}

	if repl {
		return &Image{
			R: r,
			m: image.NewUniform(c),
		}, nil
	} else {
		m := image.NewRGBA(r)
		draw.Draw(m, m.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
		return &Image{
			R: r,
			m: m,
		}, nil

	}
}

// Attach (re-)attaches to a display, typically after a resize, updating the
// display's associated image, screen, and screen image data structures.
func (d *Display) Attach(ref int) error {
	return nil // TODO: do we need this?
}

// Close closes the window.
func (d *Display) Close() error {
	e := lifecycle.Event{
		To: lifecycle.StageDead,
	}
	d.window.Send(e)
	return nil
}

// Flush flushes pending I/O to the server, making any drawing changes visible.
func (d *Display) Flush() error {
	d.ScreenImage.Lock()
	defer d.ScreenImage.Unlock()

	d.window.Upload(image.Point{}, d.buffer, d.buffer.Bounds())
	d.window.Publish()
	return nil
}

// InitMouse connects to the mouse and returns a Mousectl to interact with it.
func (d *Display) InitMouse() *Mousectl {
	return &d.mouse
}

// Moveto moves the mouse cursor to the specified location.
func (d *Display) MoveTo(pt image.Point) error {
	// Uncomment for cursor calibration:
	// fmt.Printf("shiny: MoveTo %v\n", pt)
	return moveTo(pt)
}

// SetDebug enables debugging for the remote devdraw server.
func (d *Display) SetDebug(debug bool) {
}

// ReadSnarf reads the snarf buffer into buf, returning the number of bytes read,
// the total size of the snarf buffer (useful if buf is too short), and any
// error. No error is returned if there is no problem except for buf being too
// short.
func (d *Display) ReadSnarf(buf []byte) (int, int, error) {
	s, err := clipboard.ReadAll()
	if err != nil {
		return 0, 0, err
	}
	src := []byte(s)
	if len(src) <= len(buf) {
		copy(buf, src)
		return len(src), len(src), nil
	} else {
		copy(buf, src[:len(buf)])
		return len(buf), len(src), fmt.Errorf("ReadSnarf: buffer is too short")
	}
}

// WriteSnarf writes the data to the snarf buffer.
func (d *Display) WriteSnarf(data []byte) error {
	return clipboard.WriteAll(string(data))
}
