package draw

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Image struct {
	sync.Mutex                 // Locking is used only internally.
	R          image.Rectangle // The extent of the image.
	m          image.Image
}

// Draw copies the source image with upper left corner p1 to the destination
// rectangle r, through the specified mask using operation SoverD. The
// coordinates are aligned so p1 in src and mask both correspond to r.min in
// the destination.
func (dst *Image) Draw(r image.Rectangle, src, mask *Image, p1 image.Point) {
	dst.Lock()
	defer dst.Unlock()

	if src == nil {
		fmt.Println("shiny: Draw: src is nil")
		return
	}
	if mask == nil {
		// fmt.Printf("shiny: Draw to %v\n", r)
		draw.Draw(dst.m.(*image.RGBA), r, src.m, p1, draw.Src)
		return
	}

	// It is assumed duit always calls Draw with a nil mask.
	fmt.Println("shiny: Draw: mask is not nil")
}

// Border draws a retangular border of size r and width n, with n positive
// meaning the border is inside r. It uses SoverD.
func (dst *Image) Border(r image.Rectangle, n int, src *Image, sp image.Point) {
	dst.Lock()
	defer dst.Unlock()

	for _, r := range imageutil.Border(r, n) {
		draw.Draw(dst.m.(*image.RGBA), r, src.m, sp, draw.Src)
	}
}

// Free is currently ignored.
// TODO: do we need anything about this?
func (i *Image) Free() error {
	return nil
}

// Load copies the pixel data from the buffer to the specified rectangle of the image.
// The buffer must be big enough to fill the rectangle.
//
// Duit calls Load with Load(rgba.Bounds(), rgba.Pix), so we assume image.RGBA Pix data.
func (dst *Image) Load(r image.Rectangle, data []byte) (int, error) {
	w, h := r.Dx(), r.Dy()
	if len(data) != 4*w*h {
		return 0, fmt.Errorf("image Load: wrong data size")
	}
	m := &image.RGBA{data, 4 * w, r}

	dst.R = r
	dst.m = m

	// Is len(data) ok? Duit does not read the first argument anyway.
	return len(data), nil
}

// String draws the string in the specified font using SoverD on the image,
// placing the upper left corner at p.
func (dst *Image) String(pt image.Point, src *Image, sp image.Point, f *Font, s string) image.Point {
	dst.Lock()
	defer dst.Unlock()

	m := dst.m.(*image.RGBA)
	ascent := f.face.Metrics().Ascent
	dot := fixed.P(pt.X, pt.Y).Add(fixed.Point26_6{Y: ascent})

	drawer := font.Drawer{
		Dst:  m,
		Src:  src.m,
		Face: f.face,
		Dot:  dot,
	}
	drawer.DrawString(s)
	dx := int(drawer.Dot.Sub(dot).X / 64)
	ret := pt.Add(image.Point{dx, 0})

	// fmt.Printf("shiny: String(%s) to %p at %d,%d => %d %d\n", s, m, pt.X, pt.Y, ret.X, ret.Y) // TODO remove
	return ret
}

// Line draws a line in the source color from p0 to p1, of thickness
// 1+2*radius, with the specified ends, using SoverD. The source is aligned so
// sp corresponds to p0. See the Plan 9 documentation for more information.
func (dst *Image) Line(p0, p1 image.Point, end0, end1, radius int, src *Image, sp image.Point) {
	dst.Lock()
	defer dst.Unlock()

	line(dst.m.(*image.RGBA), p0.X, p0.Y, p1.X, p1.Y, src.m.At(0, 0))
}

// Line draws a line with Besenham's algorithm.
// It only uses integer pixels.
func line(m *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	var dx, dy, sx, sy, e, e2 int

	dx = abs(x1 - x0)
	dy = -abs(y1 - y0)
	if sx = -1; x0 < x1 {
		sx = 1
	}
	if sy = -1; y0 < y1 {
		sy = 1
	}
	e = dx + dy
	for {
		m.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		if e2 = 2 * e; e2 >= dy {
			e += dy
			x0 += sx
		} else if e2 <= dx {
			e += dx
			y0 += sy
		}
	}
}
