package abc

import (
	"bytes"
	"testing"
)

var cases = []struct {
	Input  []byte
	Output uint32
	Failed bool
}{
	{[]byte{0x01}, 0x01, false},
	{[]byte{0x02}, 0x02, false},
	{[]byte{0x7f}, 0x7f, false},
	{[]byte{0x00}, 0x00, false},
	{[]byte{0x80}, 0x00, true},
	{[]byte{0x80, 0x00}, 0x00, false},
	{[]byte{0x80, 0x7f}, 0x3f80, false},
	{[]byte{0xca, 0x7f}, 0x3fca, false},
	{[]byte{0x80, 0x80, 0x80, 0x80, 0x80}, 0x00, true},
	{[]byte{0x80, 0x80, 0x80, 0x80, 0x7f}, 0x00, true},
	{[]byte{0x80, 0x80, 0x80, 0x80, 0x04}, 0x00, true},
	{[]byte{0x80, 0x80, 0x80, 0x80, 0x03}, 0x30000000, false},
	{[]byte{0x81, 0x81, 0x81, 0x81, 0x03}, 0x30204081, false},
}

func TestU30(t *testing.T) {
	for _, c := range cases {
		v, e := readU30(bytes.NewReader(c.Input))
		if e != nil && !c.Failed {
			t.Errorf("readU30(% X): %s", c.Input, e)
		} else if e == nil && c.Failed {
			t.Errorf("readU30(% X): Should have failed!", c.Input)
		} else if v != c.Output {
			t.Errorf("readU30(% X) != %d, but = %d", c.Input, c.Output, v)
		}
	}
}
