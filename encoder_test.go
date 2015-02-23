package srt

import (
	"io"
	"os"
	"strings"
	"testing"
)

const input = `1
00:00:00,000 --> 00:00:01,000
one line

2
00:00:01,123 --> 00:00:02,042
one
two lines

3
12:34:56,789 --> 23:59:59,999
one
two
three lines`

func TestDecode(t *testing.T) {
	dec := NewDecoder(strings.NewReader(input))
	f, err := os.Create(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	enc := NewEncoder(f, 0)
	defer enc.Flush()
	for b, err := dec.Next(); err != io.EOF; b, err = dec.Next() {
		if err != nil {
			t.Fatal(err)
		}

		if err := enc.Block(b); err != nil {
			t.Fatal(err)
		}
	}
}
