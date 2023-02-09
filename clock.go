package main

import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (r RealClock) Now() time.Time {
	return time.Now()
}

type FixedClock struct {
	now time.Time
}

func (f FixedClock) Now() time.Time {
	return f.now
}
