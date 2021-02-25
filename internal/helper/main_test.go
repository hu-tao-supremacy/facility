package helper

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestDayDifference(t *testing.T) {
	assert := assert.New(t)

	mockStart := time.Now()
	mockEnd := time.Now()
	diff1 := DayDifference(mockStart, mockEnd)
	assert.Equal(0, diff1, "DayDifference should be zero")

	const layoutISO = "2006-01-02"

	var tests = []struct {
		start    string
		end      string
		expected int
	}{
		{"1999-12-31", "2000-01-01", 1},
		{"2000-01-01", "1999-12-31", -1},
		{"1999-12-31", "1999-12-25", -6},
		{"2000-12-31", "2003-10-25", 1028},
		{"2003-10-25", "2000-12-31", -1028},
		{"2000-12-30", "2020-12-31", 7306},
		{"2020-12-31", "2000-12-30", -7306},
	}

	for _, test := range tests {
		mockStart, _ = time.Parse(layoutISO, test.start)
		mockEnd, _ = time.Parse(layoutISO, test.end)
		diff1 = DayDifference(mockStart, mockEnd)
		assert.Equal(test.expected, diff1)
	}

	layout2 := "Mon Jan 02 15:04:05 2006"

	var tests2 = []struct {
		start    string
		end      string
		expected int
	}{
		{"Mon Jan 02 23:54:05 2006", "Mon Jan 02 05:04:05 2006", 0},
		{"Mon Jan 02 23:59:59 2006", "Mon Jan 03 00:00:00 2006", 1},
		{"Mon Jan 03 00:00:00 2006", "Mon Jan 02 23:59:59 2006", -1},
	}

	for _, test := range tests2 {
		mockStart, _ = time.Parse(layout2, test.start)
		mockEnd, _ = time.Parse(layout2, test.end)
		diff1 = DayDifference(mockStart, mockEnd)
		assert.Equal(test.expected, diff1)
	}
}

func TestTimeStampToText(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		input    int64
		layout   string
		expected string
	}{
		{1614177990, "2006-01-02", "2021-02-24"},
		{16141779904, "2006-01-02", "2481-07-06"},
		{16141779904, "2006-01-02 15:04:05", "2481-07-06 03:45:04"},
		{0, "2006-01-02 15:04:05", "1970-01-01 00:00:00"},
	}

	for _, test := range tests {
		timeStamp := timestamp.Timestamp{Seconds: test.input}
		text := TimeStampToText(&timeStamp, test.layout)
		assert.Equal(test.expected, text, "timestamp should be correct")
	}
}
