package duitdraw

import (
	"bytes"
	"testing"
)

func TestSnarf(t *testing.T) {
	tt := []struct {
		input  []byte
		output []byte
		nbuf   int
		err    error
	}{
		{nil, nil, 10, nil},
		{[]byte("Hello"), []byte("Hello"), 5, nil},
		{[]byte("Hello, 世界"), []byte("Hello, 世界"), 100, nil},
		{[]byte("one\ntwo\three\n"), []byte("one\ntwo\three\n"), 100, nil},
		{[]byte("0123456789"), []byte("0123456"), 7, errShortSnarfBuffer},
	}
	for _, tc := range tt {
		var d Display
		err := d.WriteSnarf(tc.input)
		if err != nil {
			t.Errorf("writing snarf buffer %q failed: %v\n", tc.input, err)
		}
		b := make([]byte, tc.nbuf)
		n, size, err := d.ReadSnarf(b)
		if err != tc.err {
			t.Errorf("reading snarf buffer %q failed: %v\n", tc.input, err)
		}
		if size != len(tc.input) {
			t.Errorf("snarf buffer size is %v after writing %v bytes\n", size, len(tc.input))
		}
		if !bytes.Equal(b[:n], tc.output) {
			t.Errorf("wrote %q to snarf buffer but read %q\n", tc.output, b[:n])
		}
	}
}
