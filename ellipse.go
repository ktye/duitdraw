package draw

import "image"

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
	// TODO
}
