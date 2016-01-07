package srt

import (
	"fmt"
	"time"
)

// Block is a single subtitle Block without a index
type Block struct {
	Start, End Time
	Content    []string
}

// Add d to Start and End time
func (b *Block) Add(d time.Duration) {
	b.Start = Time{b.Start.Add(d)}
	b.End = Time{b.End.Add(d)}
}

//Time is a time.Time with a Stringer and custom Parser
type Time struct {
	time.Time
}

func (t Time) String() string {
	s := t.Format(timeFormat)
	return s[:timeCommaOff] + "," + s[1+timeCommaOff:]
}

//
func Parse(s string) (Time, error) {
	if len(s) < timeLen {
		return Time{}, fmt.Errorf("string too short to be a time: %q", s)
	}
	s = s[:timeCommaOff] + "." + s[1+timeCommaOff:timeLen]
	t, err := time.Parse(timeFormat, s)
	return Time{t}, err
}

const (
	timeDelim    = " --> "
	timeFormat   = "15:04:05.000"
	timeLen      = len(timeFormat)
	etimeOff     = timeLen + len(timeDelim)
	timeLineLen  = etimeOff + timeLen
	timeCommaOff = 8 // needed for time.Parse workaround
)
