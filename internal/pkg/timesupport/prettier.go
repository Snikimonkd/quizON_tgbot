package timesupport

import "time"

func Pretty(t time.Time) time.Time  {
	return t.In(LocMsk).Truncate(time.Millisecond)
}
