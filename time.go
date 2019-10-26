package echor

import (
	"time"
)

// NowString function
func NowString() string {
	dt := time.Now()
	return dt.String()
}

// NowMillis function
func NowMillis() int64 {
	return int64(time.Now().UnixNano() / 1000 / 1000)
}

// NxMinutes func
func NxMinutes(n int64) int64 {
	return n * 60 * 1000
}

// LastMinutesFromNow func
func LastMinutesFromNow(n int64) int64 {
	return NowMillis() - NxMinutes(n)
}

// LastMinutes func
func LastMinutes(ts int64, n int64) int64 {
	return ts - NxMinutes(n)
}
