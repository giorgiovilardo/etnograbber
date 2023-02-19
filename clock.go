package main

import "time"

type FunctionalClock struct {
	now func() time.Time
}

func (f FunctionalClock) Now() time.Time {
	return f.now()
}

// NewRealClock produces an instance of FunctionalClock, implementing Clock,
// with time.Now() preloaded as a time producing fn.
func NewRealClock() Clock {
	return FunctionalClock{now: func() time.Time {
		return time.Now()
	}}
}

// NewBrokenClock produces an instance of FunctionalClock, implementing Clock,
// with the passed time.Time as the returning value of Clock.Now()
func NewBrokenClock(t time.Time) Clock {
	return FunctionalClock{now: func() time.Time {
		return t
	}}
}
