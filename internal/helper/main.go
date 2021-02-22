package helper

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// DayDifference is a function to find day difference in time
func DayDifference(start time.Time, end time.Time) int {
	days := -start.YearDay()
	for year := start.Year(); year < end.Year(); year++ {
		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	}
	days += end.YearDay()

	return days
}

// TimeStampToText is a function to convert timestamp to text layout
func TimeStampToText(time *timestamp.Timestamp, layout string) string {
	timeDate, _ := ptypes.Timestamp(time)
	return timeDate.Format(layout)
}
