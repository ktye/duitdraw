package draw

// A Color represents an RGBA value, 8 bits per element. Red is the high 8
// bits, green the next 8 and so on.
type Color uint32

const (
	Transparent Color = 0x00000000 /* only useful for allocimage memfillcolor */
	White       Color = 0xFFFFFFFF

/*
	Opaque        Color = 0xFFFFFFFF
	Black         Color = 0x000000FF
	White         Color = 0xFFFFFFFF
	Red           Color = 0xFF0000FF
	Green         Color = 0x00FF00FF
	Blue          Color = 0x0000FFFF
	Cyan          Color = 0x00FFFFFF
	Magenta       Color = 0xFF00FFFF
	Yellow        Color = 0xFFFF00FF
	Paleyellow    Color = 0xFFFFAAFF
	Darkyellow    Color = 0xEEEE9EFF
	Darkgreen     Color = 0x448844FF
	Palegreen     Color = 0xAAFFAAFF
	Medgreen      Color = 0x88CC88FF
	Darkblue      Color = 0x000055FF
	Palebluegreen Color = 0xAAFFFFFF
	Paleblue      Color = 0x0000BBFF
	Bluegreen     Color = 0x008888FF
	Greygreen     Color = 0x55AAAAFF
	Palegreygreen Color = 0x9EEEEEFF
	Yellowgreen   Color = 0x99994CFF
	Medblue       Color = 0x000099FF
	Greyblue      Color = 0x005DBBFF
	Palegreyblue  Color = 0x4993DDFF
	Purpleblue    Color = 0x8888CCFF

	Notacolor Color = 0xFFFFFF00
	Nofill    Color = Notacolor
*/
)
