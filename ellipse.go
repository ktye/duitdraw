package duitdraw

import (
	"image"
	"image/color"
	"image/draw"
)

// How to draw ellipses?
// As I understand, low level rasterizers don't have a function to draw arcs directly, you need to create the paths before.
//
// This is how srwiley does it in oksvg:
// https://github.com/srwiley/oksvg/blob/bb18c2355556cc1db9c11f026bfe350deab1a482/svgp.go#L321
//
// This is what I did before to draw full circles:
/*
type circle struct {
        x, y, r fixed.Int26_6
}

// getPath approximates a circle by 8 quadrativ curve segments.
func (c circle) getPath() raster.Path {
        d := fixed.Point26_6{c.x, c.y}
        r := c.r
        s := fixed.Int26_6(float64(c.r) * math.Sqrt(2.0) / 2.0)
        t := fixed.Int26_6(float64(c.r) * math.Tan(math.Pi/8))
        P := func(x, y fixed.Int26_6) fixed.Point26_6 {
                return fixed.Point26_6{x, y}
        }
        var path raster.Path
        path.Start(d.Add(P(r, 0)))
        path.Add2(d.Add(P(r, t)), d.Add(P(s, s)))
        path.Add2(d.Add(P(t, r)), d.Add(P(0, r)))
        path.Add2(d.Add(P(-t, r)), d.Add(P(-s, s)))
        path.Add2(d.Add(P(-r, t)), d.Add(P(-r, 0)))
        path.Add2(d.Add(P(-r, -t)), d.Add(P(-s, -s)))
        path.Add2(d.Add(P(-t, -r)), d.Add(P(0, -r)))
        path.Add2(d.Add(P(t, -r)), d.Add(P(s, -s)))
        path.Add2(d.Add(P(r, -t)), d.Add(P(r, 0)))
        return path
}
*/
// How accurate does it have to be anyway?
// Is it used only for small segments (like 3 pixel edges to input boxes)?

// Arc draws, using SoverD, the arc centered at c, with thickness 1+2*thick,
// using the specified source color. The arc starts at angle alpha and extends
// counterclockwise by phi; angles are measured in degrees from the x axis.
func (dst *Image) Arc(c image.Point, a, b, thick int, src *Image, sp image.Point, alpha, phi int) {
	// doellipse('e', dst, c, a, b, thick, src, sp, uint32(alpha)|1<<31, phi, SoverD)
	//
	// Plan9 draw(3):
	// ellipse(dst, c, a, b, thick, src, sp)
	// Ellipse draws in dst an ellipse centered on c with horizontal
	// and vertical semiaxes a and b.
	// The source is aligned so sp in src corresponds to c in dst.
	// The ellipse is drawn with thickness 1+2*thick.
	//
	// arc(dst, c, a, b, thick, src, sp, alpha, phi)
	// Arc is like ellipse, but draws only that portion of the ellipse
	// starting at angle alpha and extending through an angle of phi.
	// The angles are measured in degrees counterclockwise from the positive x axis.

	// For full circles we assume that a==b and ignore thick.
	if alpha == 0 && phi == 360 {
		dst.Lock()
		defer dst.Unlock()
		drawCircle(dst.m.(*image.RGBA), c.X, c.Y, a, src.m.At(0, 0))
		return
	}

	// We assume duit only calls with phi=90 and alpha: 0, 90, 180, 270.
	var p0, p1 image.Point
	switch alpha {
	case 0:
		p0 = image.Point{a, 0}
		p1 = image.Point{0, -b}
	case 90:
		p0 = image.Point{0, -b}
		p1 = image.Point{-a, 0}
	case 180:
		p0 = image.Point{-a, 0}
		p1 = image.Point{0, b}
	case 270:
		p0 = image.Point{0, b}
		p1 = image.Point{a, 0}
	}

	// We just draw a line.
	dst.Line(p0.Add(c), p1.Add(c), 0, 0, 0, src, sp)

	// fmt.Printf("Arc: a=%d b=%d thick=%d alpha=%d phi=%d\n", a, b, thick, alpha, phi)
}

// FillArc draws and fills, using SoverD, the arc centered at c, with thickness
// 1+2*thick, using the specified source color. The arc starts at angle alpha
// and extends counterclockwise by phi; angles are measured in degrees from the
// x axis.
func (dst *Image) FillArc(c image.Point, a, b, thick int, src *Image, sp image.Point, alpha, phi int) {
	//doellipse('E', dst, c, a, b, thick, src, sp, uint32(alpha)|1<<31, phi, SoverD)

	// For full arcs, we assume a==b and ignore thick.
	if alpha == 0 && phi == 360 {
		dst.Lock()
		defer dst.Unlock()
		fillCircle(dst.m.(*image.RGBA), c.X, c.Y, a, src.m)
	}
}

// drawCircle is a simple rasterizer for a circle with integer pixel coordinates and a thin border.
func drawCircle(im draw.Image, xm, ym, r int, c color.Color) {
	var x, y, e int
	x = -r
	e = 2 - 2*r
	for x < 0 {
		im.Set(xm-x, ym+y, c)
		im.Set(xm-y, ym-x, c)
		im.Set(xm+x, ym-y, c)
		im.Set(xm+y, ym+x, c)
		r = e
		if r <= y {
			y++
			e += 2*y + 1
		}
		if r > x || e > y {
			x++
			e += 2*x + 1
		}
	}
}

// fillCircle fills a circle using a mask.
func fillCircle(im draw.Image, xm, ym, r int, src image.Image) {
	draw.DrawMask(im, im.Bounds(), src, image.ZP, &circle{image.Point{xm, ym}, r}, image.ZP, draw.Over)
}

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}
func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r-1, c.p.Y-c.r-1, c.p.X+c.r+1, c.p.Y+c.r+1)
}
func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X), float64(y-c.p.Y), float64(c.r)+0.5
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
