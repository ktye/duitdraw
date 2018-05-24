package draw

// Pix represents a pixel format described simple notation: r8g8b8 for RGB24, m8
// for color-mapped 8 bits, etc. The representation is 8 bits per channel,
// starting at the low end, with each byte represnted as a channel specifier
// (CRed etc.) in the high 4 bits and the number of pixels in the low 4 bits.
type Pix uint32

const (
	CRed = iota
	CGreen
	CBlue
	CGrey
	CAlpha
	CMap
	CIgnore
	NChan
)

var ARGB32 = MakePix(CAlpha, 8, CRed, 8, CGreen, 8, CBlue, 8) // stupid VGAs
var ABGR32 = MakePix(CAlpha, 8, CBlue, 8, CGreen, 8, CRed, 8)

// MakePix returns a Pix by placing the successive integers into 4-bit nibbles, low bits first.
func MakePix(list ...int) Pix {
	var p Pix
	for _, x := range list {
		p <<= 4
		p |= Pix(x)
	}
	return p
}
