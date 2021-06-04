package util

import (
	"time"

	"github.com/kroppt/voxels/log"
)

// StopWatch is a time.Time with a stopping methods
type StopWatch struct {
	t time.Time
}

// Start returns a newly started stopwatch
func Start() StopWatch {
	return StopWatch{time.Now()}
}

// Stop prints the time duration since the stopwatch start
func (sw StopWatch) Stop(str string) {
	log.Perff("%v=%v\n", str, time.Since(sw.t))
}

// StopGetNano returns the nanoseconds from the stopwatch start
func (sw StopWatch) StopGetNano() int64 {
	return time.Since(sw.t).Nanoseconds()
}

func (sw StopWatch) StopRecordAverage(key string) {
	RecordAverageTime(key, time.Since(sw.t).Nanoseconds())
}
