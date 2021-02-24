package helper

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestDayDifference(t *testing.T) {
	mockStart := time.Now()
	mockEnd := time.Now()
	diff1 := DayDifference(mockStart, mockEnd)
	assert.Equal(t, 0, diff1, "DayDifference should be zero")

	const layoutISO = "2006-01-02"

	date := "1999-12-31"
	mockStart, _ = time.Parse(layoutISO, date)
	date = "2000-01-01"
	mockEnd, _ = time.Parse(layoutISO, date)
	diff1 = DayDifference(mockStart, mockEnd)
	assert.Equal(t, 1, diff1, "DayDifference should be one")

	date = "2000-01-01"
	mockStart, _ = time.Parse(layoutISO, date)
	date = "1999-12-31"
	mockEnd, _ = time.Parse(layoutISO, date)
	diff1 = DayDifference(mockStart, mockEnd)
	assert.Equal(t, -1, diff1, "DayDifference should be minus one")

	date = "1999-12-31"
	mockStart, _ = time.Parse(layoutISO, date)
	date = "1999-12-25"
	mockEnd, _ = time.Parse(layoutISO, date)
	diff1 = DayDifference(mockStart, mockEnd)
	assert.Equal(t, -6, diff1, "DayDifference should be minus six")

	date = "2000-12-31"
	mockStart, _ = time.Parse(layoutISO, date)
	date = "2003-10-25"
	mockEnd, _ = time.Parse(layoutISO, date)
	diff1 = DayDifference(mockStart, mockEnd)
	assert.Equal(t, 1028, diff1)

	date = "2003-10-25"
	mockStart, _ = time.Parse(layoutISO, date)
	date = "2000-12-31"
	mockEnd, _ = time.Parse(layoutISO, date)
	diff1 = DayDifference(mockStart, mockEnd)
	assert.Equal(t, -1028, diff1)

	date = "2000-12-30"
	mockStart, _ = time.Parse(layoutISO, date)
	date = "2020-12-31"
	mockEnd, _ = time.Parse(layoutISO, date)
	diff1 = DayDifference(mockStart, mockEnd)
	assert.Equal(t, 7306, diff1)

	date = "2020-12-31"
	mockStart, _ = time.Parse(layoutISO, date)
	date = "2000-12-30"
	mockEnd, _ = time.Parse(layoutISO, date)
	diff1 = DayDifference(mockStart, mockEnd)
	assert.Equal(t, -7306, diff1)

	layout2 := "Mon Jan 02 15:04:05 2006"
	date = "Mon Jan 02 23:54:05 2006"
	mockStart, _ = time.Parse(layout2, date)
	date = "Mon Jan 02 05:04:05 2006"
	mockEnd, _ = time.Parse(layout2, date)
	diff1 = DayDifference(mockStart, mockEnd)
	assert.Equal(t, 0, diff1)

	date = "Mon Jan 02 23:59:59 2006"
	mockStart, _ = time.Parse(layout2, date)
	date = "Mon Jan 03 00:00:00 2006"
	mockEnd, _ = time.Parse(layout2, date)
	diff1 = DayDifference(mockStart, mockEnd)
	assert.Equal(t, 1, diff1)

	date = "Mon Jan 03 00:00:00 2006"
	mockStart, _ = time.Parse(layout2, date)
	date = "Mon Jan 02 23:59:59 2006"
	mockEnd, _ = time.Parse(layout2, date)
	diff1 = DayDifference(mockStart, mockEnd)
	assert.Equal(t, -1, diff1)

}

func TestTimeStampToText(t *testing.T) {
	timeStamp := timestamp.Timestamp{Seconds: 1614177990}
	layout := "2006-01-02"
	text := TimeStampToText(&timeStamp, layout)
	assert.Equal(t, text, "2021-02-24", "timestamp should be correct")

	timeStamp = timestamp.Timestamp{Seconds: 16141779904}
	layout = "2006-01-02"
	text = TimeStampToText(&timeStamp, layout)
	assert.Equal(t, text, "2481-07-06", "timestamp should be correct")

	timeStamp = timestamp.Timestamp{Seconds: 16141779904}
	layout = "2006-01-02 15:04:05"
	text = TimeStampToText(&timeStamp, layout)
	assert.Equal(t, text, "2481-07-06 03:45:04", "timestamp should be correct")

	timeStamp = timestamp.Timestamp{Seconds: 0}
	layout = "2006-01-02 15:04:05"
	text = TimeStampToText(&timeStamp, layout)
	assert.Equal(t, text, "1970-01-01 00:00:00", "timestamp should be correct")
}
