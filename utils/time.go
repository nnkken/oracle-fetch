package utils

import (
	"testing"
	"time"
)

var TimeNow = time.Now

func MockTimeNow(t *testing.T, timePoint time.Time) {
	TimeNow = func() time.Time {
		return timePoint
	}
	t.Cleanup(UnmockTimeNow)
}

func UnmockTimeNow() {
	TimeNow = time.Now
}
