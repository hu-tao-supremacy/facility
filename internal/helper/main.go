package helper

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// DayDifferenceFunc is type for DayDifference function
type DayDifferenceFunc func(start time.Time, end time.Time) int

// DayDifference is a function to find day difference in time
func DayDifference(start time.Time, end time.Time) int {
	var isStartAfterEnd bool
	if start.After(end) {
		start, end = end, start
		isStartAfterEnd = true
	}

	days := -start.YearDay()
	for year := start.Year(); year < end.Year(); year++ {
		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	}
	days += end.YearDay()

	if isStartAfterEnd {
		return -days
	}
	return days
}

// TimeStampToText is a function to convert timestamp to text layout
func TimeStampToText(time *timestamp.Timestamp, layout string) string {
	timeDate, _ := ptypes.Timestamp(time)
	return timeDate.Format(layout)
}
