package srt

import (
	"io"
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
	enc := NewDecoder(strings.NewReader(input))
	for _, err := enc.Next(); err != io.EOF; _, err = enc.Next() {
		if err != nil {
			t.Fatal(err)
		}
	}
}
