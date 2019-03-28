// This is a modified version of the X11 clipboard implementation in nucular
// (https://github.com/aarzilli/nucular/blob/7a48478aebff2ca5dadd95d52d34be4cb24af4ec/clipboard/clipboard_linux.go)
// distributed under the following license:
//
// The MIT License (MIT)
//
// Copyright (c) 2016 Alessandro Arzilli
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// +build freebsd linux,!android netbsd openbsd solaris dragonfly

package duitdraw

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func (d *Display) readSnarf(buf []byte) (int, int, error) {
	xc, err := getXClip()
	if err != nil {
		return 0, 0, err
	}
	xc.mu.Lock()
	defer xc.mu.Unlock()

	n, size, err := xc.getSelection(primaryAtom, buf)
	if size > 0 {
		return n, size, err
	}
	return xc.getSelection(clipboardAtom, buf)
}

func (d *Display) writeSnarf(text []byte) error {
	xc, err := getXClip()
	if err != nil {
		return err
	}
	xc.mu.Lock()
	defer xc.mu.Unlock()

	xc.text = text
	ssoc := xproto.SetSelectionOwnerChecked(xc.conn, xc.win, clipboardAtom, xproto.TimeCurrentTime)
	if err := ssoc.Check(); err != nil {
		return fmt.Errorf("error setting clipboard: %v", err)
	}
	ssoc = xproto.SetSelectionOwnerChecked(xc.conn, xc.win, primaryAtom, xproto.TimeCurrentTime)
	if err := ssoc.Check(); err != nil {
		return fmt.Errorf("error setting primary selection: %v", err)
	}
	return nil
}

const debugClipboardRequests = false

type xClip struct {
	conn      *xgb.Conn
	win       xproto.Window
	text      []byte
	selnotify chan bool
	err       error
	mu        sync.Mutex
}

var clipboardAtom, primaryAtom, textAtom, targetsAtom, atomAtom xproto.Atom

var (
	xclip     *xClip
	xclipOnce sync.Once
)

func getXClip() (*xClip, error) {
	xclipOnce.Do(func() {
		var xc xClip
		xclip = &xc

		xc.conn, xc.err = xgb.NewConnDisplay("")
		if xc.err != nil {
			return
		}

		xc.selnotify = make(chan bool, 1)

		xc.win, xc.err = xproto.NewWindowId(xc.conn)
		if xc.err != nil {
			return
		}

		setup := xproto.Setup(xc.conn)
		s := setup.DefaultScreen(xc.conn)
		xc.err = xproto.CreateWindowChecked(xc.conn, s.RootDepth, xc.win, s.Root, 100, 100, 1, 1, 0, xproto.WindowClassInputOutput, s.RootVisual, 0, []uint32{}).Check()
		if xc.err != nil {
			return
		}

		clipboardAtom = xc.internAtom("CLIPBOARD")
		primaryAtom = xc.internAtom("PRIMARY")
		textAtom = xc.internAtom("UTF8_STRING")
		targetsAtom = xc.internAtom("TARGETS")
		atomAtom = xc.internAtom("ATOM")

		go xc.eventLoop()
	})
	return xclip, xclip.err
}

func (xc *xClip) setError(err error) {
	if xc.err == nil && err != nil {
		xc.err = err
	}
}

func (xc *xClip) getSelection(selAtom xproto.Atom, buf []byte) (int, int, error) {
	err := xproto.ConvertSelectionChecked(xc.conn, xc.win, selAtom, textAtom, selAtom, xproto.TimeCurrentTime).Check()
	if err != nil {
		return 0, 0, err
	}

	select {
	case r := <-xc.selnotify:
		if !r {
			return 0, 0, fmt.Errorf("bad response from selection owner")
		}
		gpr, err := xproto.GetProperty(xc.conn, true, xc.win, selAtom, textAtom, 0, uint32(len(buf))).Reply()
		if err != nil {
			return 0, 0, err
		}
		n := copy(buf, gpr.Value[:gpr.ValueLen])
		if n < int(gpr.ValueLen) || gpr.BytesAfter != 0 {
			return n, int(gpr.ValueLen + gpr.BytesAfter), errShortSnarfBuffer
		}
		return n, n, nil
	case <-time.After(1 * time.Second):
		return 0, 0, fmt.Errorf("clipboard retrieval failed, timeout")
	}
}

func (xc *xClip) eventLoop() {
	targetAtoms := []xproto.Atom{targetsAtom, textAtom}

	for {
		e, err := xc.conn.WaitForEvent()
		if err != nil {
			continue
		}

		switch e := e.(type) {
		case xproto.SelectionRequestEvent: // write snarf
			if debugClipboardRequests {
				tgtname := xc.lookupAtom(e.Target)
				fmt.Fprintln(os.Stderr, "SelectionRequest", e, textAtom, tgtname, "isPrimary:", e.Selection == primaryAtom, "isClipboard:", e.Selection == clipboardAtom)
			}
			t := xc.text

			switch e.Target {
			case textAtom:
				if debugClipboardRequests {
					fmt.Fprintln(os.Stderr, "Sending as text")
				}
				err := xproto.ChangePropertyChecked(xc.conn, xproto.PropModeReplace, e.Requestor, e.Property, textAtom, 8, uint32(len(t)), []byte(t)).Check()
				if err == nil {
					xc.sendSelectionNotify(e)
				} else {
					fmt.Fprintf(os.Stderr, "duitdraw: %v\n", err)
				}

			case targetsAtom:
				if debugClipboardRequests {
					fmt.Fprintln(os.Stderr, "Sending targets")
				}
				buf := make([]byte, len(targetAtoms)*4)
				for i, atom := range targetAtoms {
					xgb.Put32(buf[i*4:], uint32(atom))
				}

				xproto.ChangePropertyChecked(xc.conn, xproto.PropModeReplace, e.Requestor, e.Property, atomAtom, 32, uint32(len(targetAtoms)), buf).Check()
				if err == nil {
					xc.sendSelectionNotify(e)
				} else {
					fmt.Fprintf(os.Stderr, "duitdraw: %v\n", err)
				}

			default:
				if debugClipboardRequests {
					fmt.Fprintln(os.Stderr, "Skipping")
				}
				e.Property = 0
				xc.sendSelectionNotify(e)
			}

		case xproto.SelectionNotifyEvent: // read snarf
			xc.selnotify <- (e.Property == clipboardAtom) || (e.Property == primaryAtom)
		}
	}
}

func (xc *xClip) sendSelectionNotify(e xproto.SelectionRequestEvent) {
	sn := xproto.SelectionNotifyEvent{
		Time:      xproto.TimeCurrentTime,
		Requestor: e.Requestor,
		Selection: e.Selection,
		Target:    e.Target,
		Property:  e.Property,
	}
	err := xproto.SendEventChecked(xc.conn, false, e.Requestor, 0, string(sn.Bytes())).Check()
	if err != nil {
		fmt.Fprintf(os.Stderr, "duitdraw: %v\n", err)
	}
}

func (xc *xClip) internAtom(n string) xproto.Atom {
	iar, err := xproto.InternAtom(xc.conn, true, uint16(len(n)), n).Reply()
	xc.setError(err)
	return iar.Atom
}

func (xc *xClip) lookupAtom(at xproto.Atom) string {
	reply, err := xproto.GetAtomName(xc.conn, at).Reply()
	if err != nil {
		panic(err)
	}
	return string(reply.Name)
}
