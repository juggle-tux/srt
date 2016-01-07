package srt

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

//
type Decoder struct {
	s *bufio.Scanner
	c int // count of scanned lines
}

//
func NewDecoder(r io.Reader) Decoder {
	return Decoder{
		s: bufio.NewScanner(r),
	}
}

// Next returns the next Subtitle line
func (d *Decoder) Next() (b Block, err error) {
	var s string

	// idx
	if s, err = d.scan(); err != nil {
		return b, d.error(err)
	}
	// we check that the Block starts with a int but ignore the value
	// since the Encoder keeps his own idx
	if _, err = strconv.Atoi(s); err != nil {
		return b, d.error(err)
	}

	// Time
	if s, err = d.scan(); err == io.EOF {
		return b, d.error(io.ErrUnexpectedEOF)
	} else if err != nil {
		return b, d.error(err)
	}
	if b.Start, b.End, err = parseTime(s); err != nil {
		return b, d.error(err)
	}

	// Content
	for s, err = d.scan(); err == nil && s != ""; s, err = d.scan() {
		b.Content = append(b.Content, s)
	}
	return b, d.error(err)
}

func (d *Decoder) scan() (string, error) {
	if !d.s.Scan() {
		if err := d.s.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	d.c++
	return d.s.Text(), nil
}

func (d Decoder) error(err error) error {
	switch err {
	case nil:
		return nil
	case io.EOF:
		return io.EOF
	default:
		return fmt.Errorf("line %d: %s", d.c, err)
	}
}

func parseTime(s string) (start, end Time, err error) {
	if len(s) < timeLineLen {
		return start, end, fmt.Errorf("TimeLine too short: %q", s)
	}

	start, err = Parse(s)
	if err != nil {
		return start, end, err
	}

	end, err = Parse(s[etimeOff:])
	return start, end, err
}

//
type Encoder struct {
	w   *bufio.Writer
	idx int
}

//
func NewEncoder(w io.Writer, idx int) Encoder {
	return Encoder{
		w:   bufio.NewWriter(w),
		idx: idx,
	}
}

//
func (e *Encoder) Block(b Block) error {
	str := strconv.Itoa(e.idx) + "\n"
	str += b.Start.String() + timeDelim + b.End.String() + "\n"
	for _, s := range b.Content {
		str += s + "\n"
	}
	str += "\n"
	e.idx++
	_, err := e.w.WriteString(str)
	return err
}

//
func (e *Encoder) Flush() error {
	return e.w.Flush()
}
