package duitdraw

import (
	"image"
	"testing"
)

func TestStringWidth(t *testing.T) {
	tt := []string{
		"Hello world!",
		"I can eat glass and it doesn't hurt me.",
		"私はガラスを食べられます。それは私を傷つけません。",
		"আমি কাঁচ খেতে পারি, তাতে আমার কোনো ক্ষতি হয় না।",
	}
	for _, tc := range tt {
		sum := 0
		for _, c := range tc {
			sum += defaultFont.StringWidth(string(c))
		}
		dx := defaultFont.StringWidth(tc)
		if dx != sum {
			t.Errorf("StringWidth(%q) is %v; expected %v", tc, dx, sum)
		}
	}
}

func TestImageString(t *testing.T) {
	d, err := Init(nil, "", "Image String test", "")
	if err != nil {
		t.Fatalf("can't open display: %v", err)
	}
	defer d.Close()

	dst := d.MakeImage(image.NewRGBA(image.Rect(0, 0, 100, 100)))
	p := dst.String(image.ZP, d.Black, image.ZP, defaultFont, "Hello, 世界")
	if p.X <= 0 || p.Y != 0 {
		t.Errorf("String returned bad point %v", p)
	}
}
