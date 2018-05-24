package draw

import (
	"fmt"
	"image"
)

// Arc draws, using SoverD, the arc centered at c, with thickness 1+2*thick,
// using the specified source color. The arc starts at angle alpha and extends
// counterclockwise by phi; angles are measured in degrees from the x axis.
func (dst *Image) Arc(c image.Point, a, b, thick int, src *Image, sp image.Point, alpha, phi int) {
	//doellipse('e', dst, c, a, b, thick, src, sp, uint32(alpha)|1<<31, phi, SoverD)
	// TODO
}

// FillArc draws and fills, using SoverD, the arc centered at c, with thickness
// 1+2*thick, using the specified source color. The arc starts at angle alpha
// and extends counterclockwise by phi; angles are measured in degrees from the
// x axis.
func (dst *Image) FillArc(c image.Point, a, b, thick int, src *Image, sp image.Point, alpha, phi int) {
	//doellipse('E', dst, c, a, b, thick, src, sp, uint32(alpha)|1<<31, phi, SoverD)
	// TODO
}
