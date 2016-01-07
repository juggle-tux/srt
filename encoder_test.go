package srt

import (
	"io"
	"os"
	"strings"
	"testing"
)

const testData = `1
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

func BenchmarkDecoder(b *testing.B) {
	b.SetBytes(int64(len(testData)))
	buf := strings.NewReader(testData)
	dec := NewDecoder(buf)
	f, err := os.Create(os.DevNull)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	enc := NewEncoder(f, 0)
	defer enc.Flush()

	for i := 0; i < b.N; i++ {
		for bl, err := dec.Next(); err != io.EOF; bl, err = dec.Next() {
			if err != nil {
				b.Fatal(err)
			}

			if err := enc.Block(bl); err != nil {
				b.Fatal(err)
			}
		}
		buf.Seek(0, 0)
	}
}

func TestDecode(t *testing.T) {
	dec := NewDecoder(strings.NewReader(testData))
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
