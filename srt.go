package srt

import "time"

type Block struct {
	Start, End Time
	Content    []string
}

type Time struct {
	time.Time
}

func (t Time) String() string {
	return t.Format(timeFormat)
}

const (
	timeFormat   = "15:04:05.000"
	timeLen      = len(timeFormat)
	etimeOff     = timeLen + len(timeDelim)
	timeLineLen  = etimeOff + timeLen
	timeCommaOff = 8 // needed for time.Parse workaround
	timeDelim    = " --> "
)
