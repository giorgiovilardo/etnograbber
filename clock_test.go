package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewBrokenClock(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
	}{
		{
			name: "should instance a clock that returns passed param 1",
			time: time.Date(2021, 8, 25, 8, 30, 0, 0, time.UTC),
		},
		{
			name: "should instance a clock that returns passed param 2",
			time: time.Date(1983, 10, 9, 8, 30, 0, 0, time.UTC),
		},
		{
			name: "should instance a clock that returns passed param 3",
			time: time.Date(1983, 10, 9, 8, 30, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clock := NewBrokenClock(tt.time)
			assert.Equal(t, tt.time, clock.Now())
		})
	}
}

func TestNewRealClock(t *testing.T) {
	t.Run("should produce a real time", func(t *testing.T) {
		clock := NewRealClock()
		assert.Implements(t, (*Clock)(nil), clock)
		assert.IsType(t, time.Time{}, clock.Now())
	})
}
