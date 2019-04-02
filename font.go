package duitdraw

import (
	"fmt"
	"image"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/ktye/duitdraw/internal/opentype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/plan9font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	FaceID // This field is not present in 9fans draw package.
	Height int
	face   font.Face
}

type FaceID struct {
	Name string
	Size int
	DPI  int
}

// FaceCache stores font.Faces.
type FaceCache struct {
	sync.Mutex
	m map[FaceID]font.Face
}

var faceCache FaceCache

// OpenFont opens a font with a given name and an optional size.
// Currently Plan 9 bitmap fonts and truetype fonts are supported.
// Plan 9 font must have filename or extension "font" and truetype font
// have the syntax "/path/to/font.ttf@12pt".
func (d *Display) OpenFont(name string) (*Font, error) {
	if s := strings.ToLower(name); filepath.Base(s) == "font" || filepath.Ext(s) == ".font" {
		return openPlan9Font(FaceID{Name: name})
	}
	size := DefaultFontSize
	if idx := strings.LastIndex(name, "@"); idx != -1 {
		ext := name[idx+1:]
		ext = strings.TrimSuffix(ext, "pt")
		if n, err := strconv.Atoi(ext); err != nil {
			return nil, fmt.Errorf("OpenFont: cannot parse font size: %s", name)
		} else {
			size = n
		}
		name = name[:idx]
	}
	return openFont(FaceID{Name: name, Size: size, DPI: d.DPI})
}

// RegisterFont adds a font face to the font cache.
func RegisterFont(id FaceID, face font.Face) {
	faceCache.Lock()
	defer faceCache.Unlock()
	faceCache.m[id] = face
}

// OpenFont loads a font from fontCache, from Disk or returns GoRegular
// if the font name is empty.
func openFont(id FaceID) (*Font, error) {
	faceCache.Lock()
	defer faceCache.Unlock()
	if f, ok := faceCache.m[id]; ok {
		return &Font{
			FaceID: id,
			Height: f.Metrics().Height.Round(),
			face:   f,
		}, nil
	}

	var ttf []byte
	if id.Name == "" {
		ttf = goregular.TTF
	} else {
		if b, err := ioutil.ReadFile(id.Name); err != nil {
			return nil, err
		} else {
			ttf = b
		}
	}

	if f, err := sfnt.Parse(ttf); err != nil {
		return nil, fmt.Errorf("%s: %s", id.Name, err)
	} else {
		ff, err := opentype.NewFace(f, &opentype.FaceOptions{
			Size:    float64(id.Size),
			DPI:     float64(id.DPI),
			Hinting: font.HintingNone,
		})
		if err != nil {
			return nil, err
		}
		face := pixFace{Face: ff}
		faceCache.m[id] = face

		return &Font{
			FaceID: id,
			Height: face.Metrics().Height.Round(),
			face:   face,
		}, nil
	}
}

func openPlan9Font(id FaceID) (*Font, error) {
	faceCache.Lock()
	defer faceCache.Unlock()
	if f, ok := faceCache.m[id]; ok {
		return &Font{
			FaceID: id,
			Height: f.Metrics().Height.Round(),
			face:   f,
		}, nil
	}

	fontData, err := ioutil.ReadFile(id.Name)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(id.Name)
	face, err := plan9font.ParseFont(fontData, func(name string) ([]byte, error) {
		return ioutil.ReadFile(filepath.Join(dir, name))
	})
	if err != nil {
		return nil, err
	}
	faceCache.m[id] = face

	return &Font{
		FaceID: id,
		Height: face.Metrics().Height.Round(),
		face:   face,
	}, nil
}

func (f *Font) SetDPI(dpi int) *Font {
	id := f.FaceID
	id.DPI = dpi
	if font, err := openFont(id); err != nil {
		return f
	} else {
		return font
	}
}

func (f Font) StringSize(s string) image.Point {
	dx := f.StringWidth(s)
	dy := f.Height
	return image.Point{dx, dy}
}

// StringWidth returns the number of horizontal pixels that would be occupied
// by the string if it were drawn using the font.
func (f *Font) StringWidth(s string) int {
	dx := 0
	for _, c := range s {
		a, ok := f.face.GlyphAdvance(c)
		if ok {
			dx += a.Round()
		}
	}
	return dx
}

// ByteWidth returns the number of horizontal pixels that would be occupied by
// the byte slice if it were drawn using the font.
func (f *Font) BytesWidth(b []byte) int {
	return f.StringWidth(string(b))
}

// RuneWidth returns the number of horizontal pixels that would be occupied by
// the rune slice if it were drawn using the font.
func (f *Font) RunesWidth(r []rune) int {
	return f.StringWidth(string(r))
}

// pixFace wraps a font.Face which ignores Kern and advances only by full pixels.
// Duit calls StringWidth on each rune to calculate coordinates and uses only ints.
type pixFace struct {
	font.Face
}

func (f pixFace) Kern(r0, r1 rune) fixed.Int26_6 {
	return 0
}

func (f pixFace) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	dr, mask, maskp, advance, ok = f.Face.Glyph(dot, r)
	advance = 64 * fixed.Int26_6(int(advance+32)/64)
	return
}

func (f pixFace) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	bounds, advance, ok = f.Face.GlyphBounds(r)
	advance = 64 * fixed.Int26_6(int(advance+32)/64)
	return
}

func (f pixFace) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	advance, ok = f.Face.GlyphAdvance(r)
	advance = 64 * fixed.Int26_6(int(advance+32)/64)
	return
}

// defaultFont is used for new Displays.
// It is GoRegular at DefaultSize for DefaultDPI.
var defaultFont *Font

func init() {
	// DefaultFont is GoRegular which is built-in.
	faceCache.m = make(map[FaceID]font.Face)
	id := FaceID{
		Name: "",
		Size: DefaultFontSize,
		DPI:  DefaultDPI,
	}
	var err error
	defaultFont, err = openFont(id)
	if err != nil {
		panic(err)
	}
	faceCache.m[id] = defaultFont.face
}
