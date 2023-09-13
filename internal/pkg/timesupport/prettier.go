package timesupport

import "time"

// Pretty - сделать время красивым
func Pretty(t time.Time) time.Time {
	return t.In(LocMsk).Truncate(time.Millisecond)
}
