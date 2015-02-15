package srt

import (
	"io"
	"log"
	"os"
	"testing"
)

import "time"

var tRef = time.Date(0, 1, 1, 12, 34, 56, int(789*time.Millisecond), time.UTC)
var cRef = "Hello\nWorld"

func testParseBlock(t *testing.T, buf []byte, ref string) {
	b, err := ParseBlock(buf)
	if err != nil {
		t.Fatal(err)
	}
	if !(b.Start.Equal(tRef)) {
		t.Fatalf("Stime: %s | ref: %s", b.Start, tRef)
	}
	if !(b.End.Equal(tRef)) {
		t.Fatalf("Etime: %s | ref: %s", b.End, tRef)
	}
	if string(b.Content) != ref {
		t.Fatalf("Content: %q ref: %q", b.Content, ref)
	}

	tmp := &Block{}
	if _, err := io.Copy(tmp, b); err != nil && err != io.EOF {
		log.Fatal(err)
	}
	if _, err := io.Copy(os.Stdout, tmp); err != nil && err != io.EOF {
		log.Fatal(err)
	}
}

var tBl = []byte("00\n12:34:56,789 -> 12:34:56,789\nHello\nWorld\n\n")

func TestParseBlock(t *testing.T) {
	testParseBlock(t, tBl, cRef)
}

var tBlCR = []byte("00\r\n12:34:56,789 -> 12:34:56,789\r\nHello\r\nWorld\r\n\r\n")

func TestParseBlockCR(t *testing.T) {
	testParseBlock(t, tBlCR, cRef)
}

func BenchmarkParseBloxk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseBlock(tBl)
	}
}
