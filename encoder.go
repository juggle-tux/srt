package srt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

type Decoder struct {
	s *bufio.Scanner
	c int // count of scanned lines
}

func NewDecoder(r io.Reader) Decoder {
	return Decoder{
		s: bufio.NewScanner(r),
	}
}

func (d *Decoder) Next() (Block, error) {
	var (
		b   Block
		err error
	)

	// idx
	if !d.scan() {
		if err = d.s.Err(); err != nil {
			return b, d.error(err)
		}
		return b, io.EOF // no scan no error we are done
	}
	// we check that the Block starts with a int but ignore the value
	// since we will build the index later out of the []Block order.
	if _, err = strconv.Atoi(d.s.Text()); err != nil {
		return b, d.error(err)
	}

	// Time
	if !d.scan() {
		if err = d.s.Err(); err != nil {
			return b, d.error(err)
		}
		return b, d.error(io.ErrUnexpectedEOF) // oops we are not done yet
	}
	if b.Start, b.End, err = parseTime(d.s.Text()); err != nil {
		return b, d.error(err)
	}

	// Content
	for d.scan() {
		l := d.s.Text()
		if l == "" {
			break
		}
		b.Content = append(b.Content, l)
	}
	if len(b.Content) < 1 {
		return b, d.error(errors.New("no content"))
	}
	return b, d.s.Err()
}

func (d *Decoder) scan() bool {
	ok := d.s.Scan()
	if ok {
		d.c++
	}
	return ok
}

func (d Decoder) error(e error) error {
	return fmt.Errorf("line %d: %s", d.c, e)
}

func parseTime(s string) (st, et Time, err error) {
	if len(s) != timeLineLen {
		return st, et, fmt.Errorf("malformated TimeLine: %q", s)
	}
	// start time
	sts := s[:timeLen]
	// workaround: time.Parse can't handle a "," as a delim betwen seconds and milli seconds
	sts = sts[0:timeCommaOff] + "." + sts[1+timeCommaOff:]
	t, err := time.Parse(timeFormat, sts)
	if err != nil {
		return st, et, err
	}
	st = Time{t}
	// end time
	ets := s[etimeOff : etimeOff+timeLen]
	ets = ets[0:timeCommaOff] + "." + ets[1+timeCommaOff:]
	t, err = time.Parse(timeFormat, ets)
	if err != nil {
		return st, et, err
	}
	et = Time{t}
	return st, et, nil
}

type Encoder struct {
	w   *bufio.Writer
	idx int
}

func NewEncoder(w io.Writer, idx int) Encoder {
	return Encoder{
		w:   bufio.NewWriter(w),
		idx: idx,
	}
}

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

func (e *Encoder) Flush() error {
	return e.w.Flush()
}
