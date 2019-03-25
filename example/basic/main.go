// +build example
//
// This build tag means that "go install github.com/ktye/duitdraw/..." doesn't
// install this example program. Use "go run main.go" to run it or "go install
// -tags=example" to install it.

// Basic is an example that demonstrates how to use the Main function to create
// one or more windows.
package main

import (
	"fmt"
	"image"
	"log"
	"strconv"

	draw "github.com/ktye/duitdraw"
)

func main() {
	draw.Main(func(dd *draw.Device) {
		for i := 0; i < 3; i++ {
			label := fmt.Sprintf("duitdraw (%v)", i+1)
			display, err := dd.NewDisplay(nil, "@16pt", label, "500x500")
			if err != nil {
				log.Fatalf("can't open display: %v\n", err)
			}
			if err := display.Attach(draw.Refnone); err != nil {
				log.Fatalf("failed to attach to window: %v", err)
			}
			mousectl := display.InitMouse()
			keyboardctl := display.InitKeyboard()

			text := strconv.Itoa(i + 1)
			redraw(display, text)

			go func() {
				for {
					select {
					case mousectl.Mouse = <-mousectl.C:
					case <-mousectl.Resize:
						if err := display.Attach(draw.Refnone); err != nil {
							log.Fatalf("failed to attach to window: %v", err)
						}
						redraw(display, text)
					case <-keyboardctl.C:
						display.Close()
						return
					}
				}
			}()
		}
	})
}

func redraw(display *draw.Display, text string) {
	r := display.ScreenImage.R
	display.ScreenImage.Draw(r, display.White, nil, image.ZP)

	p0 := image.Pt(r.Dx()/2, r.Dy()/2)
	display.ScreenImage.String(p0, display.Black, p0, display.DefaultFont, text)

	display.Flush()
}
