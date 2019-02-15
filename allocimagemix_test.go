package duitdraw

import (
	"image"
	"reflect"
	"testing"
)

func TestAllocImageMix(t *testing.T) {
	tt := []struct {
		color1, color3, mix Color
	}{
		{Palebluegreen, White, 0xEAFFFFFF},
		{Paleyellow, White, 0xFFFFEAFF},
		{0xAAAAAAAA, 0x55555555, 0x6A6A6AAA},
		{0x55555555, 0xAAAAAAAA, 0x95959555},
		{0x0A0A0A0A, 0x05050505, 0x0606060A},
		{0x05050505, 0x0A0A0A0A, 0x08080805},
	}
	d, err := Init(nil, "", "", "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	r := image.Rect(0, 0, 1, 1)
	for _, tc := range tt {
		b := d.AllocImageMix(tc.color1, tc.color3)
		if b.Display != d {
			t.Errorf("image has display %p; exptected %p", b.Display, d)
		}
		if !reflect.DeepEqual(b.R, r) {
			t.Errorf("image rect is %v; expected %v", b.R, r)
		}
		bm := image.NewUniform(tc.mix.rgba())
		if !reflect.DeepEqual(bm, b.m) {
			t.Errorf("image is %X; exptected %X", b.m, bm)
		}
	}
}
