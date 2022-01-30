package tick

import "time"

type Time interface {
	Now() time.Time
}

type FnTime struct {
	FnNow func() time.Time
}

func (fn FnTime) Now() time.Time {
	if fn.FnNow != nil {
		return fn.FnNow()
	}
	return time.Now()
}
