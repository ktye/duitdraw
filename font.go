package draw

import (
	"image"

	"golang.org/x/image/font"
)

type Font struct {
	Height int
	face   font.Face
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
