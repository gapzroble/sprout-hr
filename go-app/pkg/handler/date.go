package handler

import (
	"time"

	"github.com/gapzroble/sprout-hr/pkg/sprout"
)

func isWeekend() bool {
	switch sprout.Now().Weekday() {
	case time.Saturday, time.Sunday:
		return true
	}
	return false
}
