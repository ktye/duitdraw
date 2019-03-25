package duitdraw

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
)

// mainScreen stores the screen which is initialized for the first window.
var mainScreen screen.Screen

// Main is called by the program's main function to run the graphical
// application.
//
// It calls f on the Device, possibly in a separate goroutine, as some OS-
// specific libraries require being on 'the main thread'. It returns when f
// returns.
func Main(f func(*Device)) {
	driver.Main(func(ss screen.Screen) {
		dev := newDevice(ss)
		f(dev)
		dev.wait()
	})
}

// Device represents the draw device on which multiple windows
// can be created.
type Device struct {
	wg sync.WaitGroup
}

func newDevice(ss screen.Screen) *Device {
	mainScreen = ss
	return &Device{}
}

// NewDisplay is called to create a new window.
// There is no special mechanism to create the first window.
func (dev *Device) NewDisplay(errch chan<- error, fontname, label, winsize string) (*Display, error) {
	if errch == nil {
		setDefaultErrorChan()
		errch = defaultErrorChan
	}
	dpy, opt := newDisplay(label, winsize, fontname)
	dev.wg.Add(1)
	go func() {
		createWindow(dpy, opt, errch)
		dev.wg.Done()
	}()

	// make sure ScreenImage buffer is allocated
	<-dpy.mouse.Resize

	return dpy, nil
}

func (dev *Device) wait() {
	dev.wg.Wait()
}

// Init is called to create a new window.
// There is no special mechanism to create the first window.
// This function does not work on systems (e.g. macOS) where the
// OS-specific graphics libraries require being on 'the main thread'.
// Use Main for best compatiblity.
func Init(errch chan<- error, fontname, label, winsize string) (*Display, error) {
	if errch == nil {
		setDefaultErrorChan()
		errch = defaultErrorChan
	}
	if mainScreen == nil {
		dpy, opt := newDisplay(label, winsize, fontname)
		go driver.Main(func(s screen.Screen) {
			mainScreen = s
			createWindow(dpy, opt, errch)
		})
		// make sure ScreenImage buffer is allocated
		<-dpy.mouse.Resize
		return dpy, nil
	} else {
		dpy, opt := newDisplay(label, winsize, fontname)
		go createWindow(dpy, opt, errch)
		// make sure ScreenImage buffer is allocated
		<-dpy.mouse.Resize
		return dpy, nil
	}
}

// NewDisplay creates a Display with it's mouse and keyboard channels.
// It registers the window in mainScreen but does not call any shiny functions.
func newDisplay(label, winsize, fontname string) (*Display, screen.NewWindowOptions) {
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
		DPI: DefaultDPI,
	}
	dpy.Black = &Image{
		Display: &dpy,
		R:       image.Rect(0, 0, 1, 1),
		m:       image.NewUniform(color.Black),
	}
	dpy.White = &Image{
		Display: &dpy,
		R:       image.Rect(0, 0, 1, 1),
		m:       image.NewUniform(color.White),
	}
	dpy.Opaque = &Image{
		Display: &dpy,
		R:       image.Rect(0, 0, 1, 1),
		m:       image.NewUniform(color.Opaque),
	}
	dpy.Transparent = &Image{
		Display: &dpy,
		R:       image.Rect(0, 0, 1, 1),
		m:       image.NewUniform(color.Transparent),
	}
	dpy.ScreenImage = &Image{
		Display: &dpy,
		R:       image.Rect(0, 0, opt.Width, opt.Height),
		// m will be backed by screen.Buffer on size event.
	}
	if f, err := dpy.OpenFont(fontname); err != nil {
		dpy.DefaultFont = defaultFont
		log.Print(err)
	} else {
		dpy.DefaultFont = f
	}
	dpy.mouse.C = make(chan Mouse, 0)
	dpy.mouse.Resize = make(chan bool, 2) // Why 2? (copied from InitMouse).
	dpy.mouse.Display = &dpy
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

var (
	defaultErrorChan     chan<- error
	defaultErrorChanOnce sync.Once
)

func setDefaultErrorChan() {
	defaultErrorChanOnce.Do(func() {
		ch := make(chan error)
		go func() {
			for err := range ch {
				if err != io.EOF {
					fmt.Fprintf(os.Stderr, "duitdraw: %v\n", err)
				}
			}
		}()
		defaultErrorChan = ch
	})
}
