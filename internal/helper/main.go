package helper

import (
	"math"
	"time"
)

// DayDifference is a function to find day difference in time
func DayDifference(start *time.Time, end *time.Time) int {
	return int(math.Ceil(end.Sub(*start).Hours() / 24))
}
