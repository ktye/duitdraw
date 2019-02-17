// +build dragonfly freebsd linux netbsd openbsd solaris

package duitdraw

import (
	"fmt"
	"image"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

var xdpy *x11Display

func moveTo(pt image.Point) error {
	d, err := getX11Display()
	if err != nil {
		return err
	}
	return xproto.WarpPointerChecked(d.conn,
		xproto.WindowNone, // src
		d.focus,           // dst
		0, 0, 0, 0,        // src X, Y, width, height
		int16(pt.X), int16(pt.Y), // dst X, Y
	).Check()
}

func setCursor(c *Cursor) error {
	d, err := getX11Display()
	if err != nil {
		return err
	}
	if c != nil {
		return d.setCursor(c)
	}
	return d.unsetCursor()
}

type x11Display struct {
	conn   *xgb.Conn
	focus  xproto.Window
	cursor xproto.Cursor // previously set cursor
}

func getX11Display() (*x11Display, error) {
	if xdpy == nil {
		conn, err := xgb.NewConn()
		if err != nil {
			return nil, err
		}
		xdpy = &x11Display{
			conn:   conn,
			cursor: xproto.CursorNone,
		}
	}
	// TODO(fhs): This is wrong if the window is not in focus.
	gif, err := xproto.GetInputFocus(xdpy.conn).Reply()
	if err != nil {
		return nil, err
	}
	xdpy.focus = gif.Focus
	return xdpy, nil
}

func (d *x11Display) freeCursor() {
	if d.cursor != xproto.CursorNone {
		xproto.FreeCursor(d.conn, d.cursor)
		d.cursor = xproto.CursorNone
	}
}

func (d *x11Display) unsetCursor() error {
	err := xproto.ChangeWindowAttributesChecked(d.conn, d.focus,
		xproto.CwCursor, []uint32{xproto.CursorNone}).Check()
	if err != nil {
		return fmt.Errorf("ChangeWindowAttributesChecked: %v", err)
	}
	d.freeCursor()
	return nil
}

func (d *x11Display) setCursor(c *Cursor) error {
	var src, mask [2 * 16]byte

	for i := 0; i < 2*16; i++ {
		src[i] = reverseByte(c.Set[i])
		mask[i] = reverseByte(c.Set[i] | c.Clr[i])
	}

	xsrc, err := createPixmapFromData(d.conn, xproto.Drawable(d.focus), src[:], 16, 16)
	if err != nil {
		return fmt.Errorf("createPixmapFromData xsrc: %v", err)
	}
	defer xproto.FreePixmap(d.conn, xsrc)
	xmask, err := createPixmapFromData(d.conn, xproto.Drawable(d.focus), mask[:], 16, 16)
	if err != nil {
		return fmt.Errorf("createPixmapFromData xmask: %v", err)
	}
	defer xproto.FreePixmap(d.conn, xmask)
	xc, err := xproto.NewCursorId(d.conn)
	if err != nil {
		return fmt.Errorf("NewCursorId: %v", err)
	}
	err = xproto.CreateCursorChecked(d.conn, xc, xsrc, xmask,
		0, 0, 0, // foreRed, foreGreen, foreblue
		0xffff, 0xffff, 0xffff, // backRed, backGreen, backBlue
		uint16(-c.Point.X), uint16(-c.Point.Y),
	).Check()
	if err != nil {
		return fmt.Errorf("CreateCursorChecked: %v", err)
	}
	err = xproto.ChangeWindowAttributesChecked(d.conn, d.focus,
		xproto.CwCursor, []uint32{uint32(xc)}).Check()
	if err != nil {
		return fmt.Errorf("ChangeWindowAttributesChecked: %v", err)
	}
	d.freeCursor()
	d.cursor = xc
	return nil
}

func createPixmapFromData(conn *xgb.Conn, drawable xproto.Drawable, data []byte, width, height uint16) (xproto.Pixmap, error) {
	pm, err := xproto.NewPixmapId(conn)
	if err != nil {
		return 0, fmt.Errorf("NewPixmapId: %v", err)
	}
	err = xproto.CreatePixmapChecked(conn, 1, pm, drawable, width, height).Check()
	if err != nil {
		return 0, fmt.Errorf("CreatePixmapChecked: %v", err)
	}
	gc, err := xproto.NewGcontextId(conn)
	if err != nil {
		return 0, fmt.Errorf("NewGcontextId: %v", err)
	}
	err = xproto.CreateGCChecked(conn, gc, xproto.Drawable(pm),
		xproto.GcForeground|xproto.GcBackground,
		[]uint32{1, 0}).Check()
	if err != nil {
		return 0, fmt.Errorf("CreateGCChecked: %v", err)
	}
	// TODO(fhs): Just guessing the pixmap binary layout,
	// but this seems to make the Edwood boxcursor work.
	data2 := make([]byte, 0, len(data)*2)
	for i := 0; i < len(data); i += 2 {
		data2 = append(data2, data[i:i+2]...)
		data2 = append(data2, data[i:i+2]...)
	}
	err = xproto.PutImageChecked(conn, xproto.ImageFormatXYPixmap, xproto.Drawable(pm), gc,
		width, height,
		0, 0, // DstX, DstY
		0, 1, // LeftPad, Depth
		data2,
	).Check()
	if err != nil {
		return 0, fmt.Errorf("PutImageChecked: %v", err)
	}
	return pm, nil
}

func reverseByte(b byte) byte {
	var r byte

	r = 0
	r |= (b & 0x01) << 7
	r |= (b & 0x02) << 5
	r |= (b & 0x04) << 3
	r |= (b & 0x08) << 1
	r |= (b & 0x10) >> 1
	r |= (b & 0x20) >> 3
	r |= (b & 0x40) >> 5
	r |= (b & 0x80) >> 7
	return r
}
