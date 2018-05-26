package draw

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"
	"time"

	"github.com/golang/freetype/truetype"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font/gofont/goregular"
)

// mainScreen stores the screen which is initialized for the first window.
var mainScreen screen.Screen

var defaultFont *Font

// Init is called to create a new window.
// There is no special mechanism to create the first window.
func Init(errch chan<- error, fontname, label, winsize string) (*Display, error) {
	if mainScreen == nil {
		dpy, opt := newWindow(label, winsize)
		go driver.Main(func(s screen.Screen) {
			mainScreen = s
			createWindow(dpy, opt, errch)
		})
		return dpy, nil
	} else {
		dpy, opt := newWindow(label, winsize)
		go createWindow(dpy, opt, errch)
		return dpy, nil
	}
}

// NewWindow creates a Display with it's mouse and keyboard channels.
// It registers the window in mainScreen but does not call any shiny functions.
func newWindow(label, winsize string) (*Display, screen.NewWindowOptions) {
	opt := screen.NewWindowOptions{
		Width:  800,
		Height: 800,
		Title:  label,
	}
	if wh := strings.Split(winsize, "x"); len(wh) == 2 {
		if w, err := strconv.Atoi(wh[0]); err == nil {
			if h, err := strconv.Atoi(wh[1]); err == nil {
				opt.Width = w
				opt.Height = h
			}
		}
	}

	dpy := Display{
		Black: &Image{
			R: image.Rect(0, 0, 1, 1),
			m: image.NewUniform(color.Black),
		},
		White: &Image{
			R: image.Rect(0, 0, 1, 1),
			m: image.NewUniform(color.White),
		},
		ScreenImage: &Image{
			R: image.Rect(0, 0, opt.Width, opt.Height),
			// m will be backed by screen.Buffer on size event.
		},
		DefaultFont: defaultFont,
		/* &Font{
			Height: int(basicfont.Face7x13.Metrics().Ascent / 64),
			face:   basicfont.Face7x13,
		}, */
	}
	dpy.mouse.C = make(chan Mouse, 0)
	dpy.mouse.Resize = make(chan bool, 2) // Why 2? (copied from InitMouse).
	dpy.mouse.last = time.Now()
	dpy.keyboard.C = make(chan rune, 20)

	return &dpy, opt
}

// CreateWindow creates a new client window and runs it.
// The function is called inside a go routine and is alive as long as the window is present.
func createWindow(d *Display, opt screen.NewWindowOptions, errch chan<- error) {
	w, err := mainScreen.NewWindow(&opt)
	if err != nil {
		fmt.Printf("shiny: NewWindow error: %s\n", err)
		errch <- err
		return
	}
	defer w.Release()

	var b screen.Buffer
	defer func() {
		if b != nil {
			b.Release()
		}
	}()

	d.window = w
	d.buffer = b
	d.eventLoop(errch)
}

func init() {
	// We use Go Regular as a default font.
	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	opt := truetype.Options{
		Size: 10,
		DPI:  DefaultDPI,
	}
	face := truetype.NewFace(f, &opt)
	defaultFont = &Font{
		Height: int(face.Metrics().Ascent / 64),
		face:   face,
	}

	/* TODO: This uses sfnt/opentype, which is not working yet.
	f, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}

	opt := opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingNone,
	}

	face, err := opentype.NewFace(f, &opt)
	if err != nil {
		panic(err)
	}

	defaultFont = &Font{
		Height: int(face.Metrics().Ascent / 64),
		face:   face,
	}
	*/
}
