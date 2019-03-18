package duitdraw

import (
	"image"
	"image/color"
	"testing"
)

func testDrawMask(t *testing.T, d *Display, c *Image, expect func(color.RGBA) color.RGBA) {
	r := image.Rect(0, 0, 2, 2)

	m := image.NewRGBA(r)
	m.SetRGBA(0, 0, color.RGBA{0x11, 0x55, 0x99, 0xdd})
	m.SetRGBA(0, 1, color.RGBA{0x22, 0x66, 0xaa, 0xee})
	m.SetRGBA(1, 0, color.RGBA{0x33, 0x77, 0xbb, 0xff})
	m.SetRGBA(1, 1, color.RGBA{0x44, 0x88, 0xcc, 0x00})
	img := d.MakeImage(m)

	dst := d.MakeImage(image.NewRGBA(r))
	dst.Draw(r, img, nil, image.ZP)

	mask := d.MakeImage(image.NewRGBA(r))
	mask.Draw(r, c, nil, image.ZP)

	dst.Draw(dst.R, d.Black, mask, image.ZP)

	q := dst.m.(*image.RGBA)
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			got := q.RGBAAt(x, y)
			want := expect(m.RGBAAt(x, y))
			if got != want {
				t.Errorf("dst value at (%v, %v) is %v; expected %v", x, y, got, want)
			}
		}
	}
}

func TestTransparent(t *testing.T) {
	d, err := Init(nil, "", "Transparen test", "")
	if err != nil {
		t.Fatalf("can't open display: %v", err)
	}
	testDrawMask(t, d, d.Transparent, func(c color.RGBA) color.RGBA {
		return c
	})
}

func TestOpaque(t *testing.T) {
	d, err := Init(nil, "", "Opaque test", "")
	if err != nil {
		t.Fatalf("can't open display: %v", err)
	}
	testDrawMask(t, d, d.Opaque, func(c color.RGBA) color.RGBA {
		return Black.rgba()
	})
}
