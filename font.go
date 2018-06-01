package draw

import (
	"image"

	"golang.org/x/image/font"
)

type Font struct {
	Height int
	face   font.Face
}

// OpenFont reads the named file and returns the font it defines. The name may
// be an absolute path, or identify a file in a standard font directory:
// /lib/font/bit, /usr/local/plan9, /mnt/font, etc.
func (d *Display) OpenFont(name string) (*Font, error) {
	// TODO
	return defaultFont, nil
}

func (f Font) StringSize(s string) image.Point {
	dx := f.StringWidth(s)
	dy := f.Height
	// fmt.Printf("shiny: StringSize(%s) %d %d\n", s, dx, dy)
	return image.Point{dx, dy}
}

// StringWidth returns the number of horizontal pixels that would be occupied
// by the string if it were drawn using the font.
func (f *Font) StringWidth(s string) int {
	dx := int(font.MeasureString(f.face, s) / 64)
	// fmt.Printf("shiny: StringWidth(%s) %d\n", s, dx)
	return dx
}
