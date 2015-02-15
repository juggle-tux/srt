package srt

import (
	"bytes"
	"errors"
	"io"
	"time"
)

const (
	TimeLayout  = "15:04:05.000"
	lenTime     = 12
	etOff       = 16
	lenTimeLine = 28
)

// errors
var (
	ErrShortString = errors.New("string too short")
	ErrParseIdx    = errors.New("string is not a number")
)

type Block struct {
	Start, End Time
	Content    []byte
}

func ParseBlock(buf []byte) (Block, error) {
	b, _, err := parseBlock(buf)
	return b, err
}

func (b Block) Bytes() []byte {
	return []byte(
		b.Start.String() + " -> " + b.End.String() + "\n" +
			string(b.Content) + "\n\n")
}

func (b *Block) Write(buf []byte) (int, error) {
	buf = append([]byte("00\n"), buf...)
	t, n, err := parseBlock(buf)
	if err != nil {
		return 0, err
	}
	*b = t
	return n, nil
}

func (b Block) Read(buf []byte) (int, error) {
	return copy(buf, b.Bytes()), io.EOF
}

type Time struct {
	time.Time
}

func (t Time) String() string {
	return t.Format(TimeLayout)
}

func parseBlock(buf []byte) (Block, int, error) {
	var err error
	b := Block{}
	bs := bytes.SplitN(buf, []byte{'\n'}, 3)
	if !isIndex(trimCR(bs[0])) {
		return b, 0, ErrParseIdx
	}
	b.Start, b.End, err = parseTimeLine(trimCR(bs[1]))
	if err != nil {
		return b, 0, err
	}
	var n int
	b.Content, n, err = parseContent(bs[2])
	if err != nil {
		return b, 0, err
	}
	return b, len(bs[0]) + len(bs[1]) + n - 1, nil
}

func parseContent(buf []byte) ([]byte, int, error) {
	l := len(buf)
	for i := 0; i < l; i++ {
		if l-i > 1 && buf[i] == byte('\n') && buf[i+1] == byte('\n') {
			buf = buf[:i]
			break
		} else if l-i > 3 && buf[i] == byte('\r') && buf[i+1] == byte('\n') && buf[i+2] == byte('\r') && buf[i+3] == byte('\n') {
			buf = buf[:i]
			buf = bytes.Replace(buf, []byte("\r\n"), []byte("\n"), -1)
			break
		}
	}
	return buf, l, nil
}

func parseTime(buf []byte) (Time, error) {
	s := string(buf)
	if len(s) < lenTime {
		return Time{}, ErrShortString
	}
	s = s[0:8] + "." + s[9:12]
	t, err := time.Parse(TimeLayout, s)
	if err != nil {
		return Time{}, err
	}
	return Time{t}, err
}

func parseTimeLine(buf []byte) (stime, etime Time, err error) {
	if len(buf) < lenTimeLine {
		return Time{}, Time{}, ErrShortString
	}

	stime, err = parseTime(buf)
	if err != nil {
		return
	}
	etime, err = parseTime(buf[etOff:])
	return
}

func trimCR(buf []byte) []byte {
	l := len(buf)
	if buf[l-1] == byte('\r') {
		buf = buf[:l-1]
	}
	return buf
}

func isIndex(buf []byte) bool {
	for _, c := range buf {
		if !(c >= byte('0') && c <= byte('9')) {
			return false
		}
	}
	return true
}
