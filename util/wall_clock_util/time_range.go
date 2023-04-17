package wall_clock_util

import (
	"time"
)

type TimeRange [2]Clock
type MaybeTimeRange []Clock

func (r MaybeTimeRange) NowInRange() bool {
	if len(r) == 0 {
		return true
	}
	rng := [2]Clock{r[0], r[1]}
	return TimeRange(rng).NowInRange()
}

func (r TimeRange) ClockInRange(cur Clock) (ret bool) {
	if r[0].LessThan(r[1]) {
		ret = r[0].Seconds() <= cur.Seconds() && cur.Seconds() <= r[1].Seconds()
	} else {
		ret = r[0].Seconds() <= cur.Seconds() || cur.Seconds() <= r[1].Seconds()
	}
	return
}
func (r TimeRange) TimeInRange(cur time.Time) bool {
	c := NewClockFromTime(cur)
	return r.ClockInRange(*c)
}
func (r TimeRange) NowInRange() bool {
	return r.TimeInRange(time.Now())
}
