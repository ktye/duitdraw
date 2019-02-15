package duitdraw

import "image"

// AllocImageMix blends the two colors to create a tiled image representing
// their combination. For pixel formats of 8 bits or less, it creates a 2x2
// pixel texture whose average value is the mix. Otherwise it creates a 1-pixel
// solid color blended using 50% alpha for each.
func (d *Display) AllocImageMix(color1, color3 Color) *Image {
	//	if d.ScreenImage.Depth <= 8 { // create a 2x2 texture
	//		t, _ := d.allocImage(image.Rect(0, 0, 1, 1), d.ScreenImage.Pix, false, color1)
	//		b, _ := d.allocImage(image.Rect(0, 0, 2, 2), d.ScreenImage.Pix, true, color3)
	//		b.draw(image.Rect(0, 0, 1, 1), t, nil, image.ZP)
	//		t.free()
	//		return b
	//	}

	// use a solid color, blended using alpha
	const (
		q1 = 0x3f
		q3 = 0xff - q1
	)
	c1 := color1.rgba()
	c3 := color3.rgba()
	r := (uint32(c1.R)*q1 + uint32(c3.R)*q3) / 0xff
	g := (uint32(c1.G)*q1 + uint32(c3.G)*q3) / 0xff
	b := (uint32(c1.B)*q1 + uint32(c3.B)*q3) / 0xff
	a := uint32(c1.A)
	c := Color(r<<24 | g<<16 | b<<8 | a)
	img, _ := d.AllocImage(image.Rect(0, 0, 1, 1), d.ScreenImage.Pix, true, c)
	return img
}
