package draw

import "golang.org/x/mobile/event/key"

// Uncommented Key constants are defined for plan9 but not used by duit.

const (
	KeyFn = '\uF000'

	KeyHome   = KeyFn | 0x0D
	KeyUp     = KeyFn | 0x0E
	KeyPageUp = KeyFn | 0xF
	//KeyPrint     = KeyFn | 0x10
	KeyLeft  = KeyFn | 0x11
	KeyRight = KeyFn | 0x12
	KeyDown  = 0x80
	//KeyView      = 0x80
	KeyPageDown = KeyFn | 0x13
	KeyInsert   = KeyFn | 0x14
	KeyEnd      = KeyFn | 0x18
	//KeyAlt       = KeyFn | 0x15
	//KeyShift     = KeyFn | 0x16
	//KeyCtl       = KeyFn | 0x17
	//KeyBackspace = 0x08
	KeyDelete = 0x7F
	KeyEscape = 0x1b
	//KeyEOF       = 0x04
	KeyCmd = 0xF100
)

// Keymap maps from key event codes to runes, that duit expects.
var keymap = map[key.Code]rune{
	key.CodeHome:       KeyHome,
	key.CodeUpArrow:    KeyUp,
	key.CodePageUp:     KeyPageUp,
	key.CodeLeftArrow:  KeyLeft,
	key.CodeRightArrow: KeyRight,
	key.CodeDownArrow:  KeyDown,
	key.CodePageDown:   KeyPageDown,
	key.CodeEnd:        KeyEnd,
	//key.CodeDelete:     KeyDelete,
	key.CodeEscape: KeyEscape,
	//key.CodeCmd:        KeyCmd,
}

// Keyboardctl is the source of keyboard events.
type Keyboardctl struct {
	C chan rune // Channel on which keyboard characters are delivered.
}

// InitKeyboard connects to the keyboard and returns a Keyboardctl to listen to it.
func (d *Display) InitKeyboard() *Keyboardctl {
	return &d.keyboard
}

// KeyTranslator translates a key.Event to a rune.
// If present, it overwrites the default mechanism.
// If TranslateKey returns -1, the key event is ignored.
type KeyTranslator interface {
	TranslateKey(key.Event) rune
}
