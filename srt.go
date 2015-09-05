package srt

import "time"

type Block struct {
	Start, End Time
	Content    []string
}

func (b *Block) Add(d time.Duration) {
	b.Start = Time{b.Start.Add(d)}
	b.End = Time{b.End.Add(d)}
}

type Time struct {
	time.Time
}

func (t Time) String() string {
	s := t.Format(timeFormat)
	return s[:timeCommaOff] + "," + s[1+timeCommaOff:]
}

const (
	timeFormat   = "15:04:05.000"
	timeLen      = len(timeFormat)
	etimeOff     = timeLen + len(timeDelim)
	timeLineLen  = etimeOff + timeLen
	timeCommaOff = 8 // needed for time.Parse workaround
	timeDelim    = " --> "
)
