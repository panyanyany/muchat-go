package wall_clock_util

import (
	"errors"
	"fmt"
	"time"
)

type Clock struct {
	Hour, Minute, Second int
}

func NewClock(h, m, s int) (*Clock, error) {
	if h > 23 || m > 59 || s > 59 {
		return nil, errors.New("invalid input")
	}
	return &Clock{h, m, s}, nil
}

func (c Clock) Seconds() int {
	return c.Hour*60*60 + c.Minute*60 + c.Second
}

func (c Clock) LessThan(d Clock) bool {
	return c.Seconds() < d.Seconds()
}

var timeLayout = "15:04:05"
var TimeParseError = errors.New(`TimeParseError: should be a string formatted as "15:04:05"`)

func (t Clock) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}
func (t Clock) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
}

func (t *Clock) UnmarshalJSON(b []byte) error {
	s := string(b)
	// len(`"23:59"`) == 7
	if len(s) != 8 {
		return TimeParseError
	}
	//ret, err := time.Parse(timeLayout, s[1:6])
	ret, err := time.Parse(timeLayout, s)
	if err != nil {
		err = fmt.Errorf("time.Parse(): %w", err)
		return err
	}
	t.Hour = ret.Hour()
	t.Minute = ret.Minute()
	t.Second = ret.Second()
	return nil
}

func NewClockFromTime(tm time.Time) (r *Clock) {
	r = new(Clock)
	r.Hour = tm.Hour()
	r.Minute = tm.Minute()
	r.Second = tm.Second()
	return
}

// Implements the Unmarshaler interface of the yaml pkg.
func (r *Clock) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	strs := ""
	err = unmarshal(&strs)
	if err != nil {
		return
	}

	err = r.UnmarshalJSON([]byte(strs))
	if err != nil {
		err = fmt.Errorf("r.UnmarshalJSON(): %w", err)
		return
	}

	return
}
func (r Clock) MarshalYAML() (res interface{}, err error) {

	res = r.String()

	return
}
