package timeparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"fmt"
)

func TestParseFormat(t *testing.T) {
	//utc, err := time.LoadLocation("UTC")
	//assert.Nil(t, err)

	// tests for time format
	expected := time.Date(2021, time.December, 29, 18, 24, 36, 0, time.Local)
	testcases := map[string]string {
		// general format
		" l jS \\of F Y  h:i:s A": "Wednesday 29th of December  2021 06:24:36 PM ",
		"  Y-m-d    H:i:s ": " 2021-12-29  18:24:36  ",

		// with wild cards
		"Y#n#j H:i:s" : "2021,12,29 18:24:36",
		"U H?i#s" : "1640769840 18_24:36",
	}
	for format, s := range testcases {
		tm, err := ParseFormat(format, s)

		assert.Nil(t, err)
		assert.Equal(t, expected.Year(), tm.Year())
		assert.Equal(t, expected.Month(), tm.Month())
		assert.Equal(t, expected.Day(), tm.Day())
		assert.Equal(t, expected.Hour(), tm.Hour())
		assert.Equal(t, expected.Minute(), tm.Minute())
		assert.Equal(t, expected.Second(), tm.Second())
		assert.Equal(t, expected.Location(), tm.Location())
	}

	// test form |+
	testcases = map[string]string {
		"Y?n*j|" : "2021a12___29 18:24:12",
		"Y#n#j|+j" : "2021/12/29 ___",
	}
	for format, s := range testcases {
		tm, err := ParseFormat(format, s)

		assert.Nil(t, err)
		assert.Equal(t, expected.Year(), tm.Year())
		assert.Equal(t, expected.Month(), tm.Month())
		assert.Equal(t, expected.Day(), tm.Day())
		assert.Equal(t, 0, tm.Hour())
		assert.Equal(t, 0, tm.Minute())
		assert.Equal(t, 0, tm.Second())
		assert.Equal(t, expected.Location(), tm.Location())
	}

	// Timezone (whether exists or not depends on your env)
	testcases_tz_s := []string {
		"UTC",
		"GMT",
		"MST",
		"America/New_York", 
		// "Asia/Tokyo",
	}
	for _, tz_ := range testcases_tz_s {
		tm, err := ParseFormat("Y-m-d H:i:s T", " 2021-12-29  18:24:36 " + tz_)

		assert.Nil(t, err)
		assert.Equal(t, expected.Year(), tm.Year())
		assert.Equal(t, expected.Month(), tm.Month())
		assert.Equal(t, expected.Day(), tm.Day())
		assert.Equal(t, expected.Hour(), tm.Hour())
		assert.Equal(t, expected.Minute(), tm.Minute())
		assert.Equal(t, expected.Second(), tm.Second())
		assert.Equal(t, tz_, tm.Location().String())
	}
	testcases_tz_i := []int{
		0,
		-1,
		-11,
		+1,
		+11,
	}
	expected = time.Date(2021, time.December, 29, 12, 24, 36, 0, time.Local)
	for _, i := range testcases_tz_i {
		tm, err := ParseFormat("Y-m-d H:i:s T", fmt.Sprintf(" 2021-12-29  12:24:36 %+02d:00", i) )

		assert.Nil(t, err)
		assert.Equal(t, expected.Year(), tm.Year())
		assert.Equal(t, expected.Month(), tm.Month())
		assert.Equal(t, expected.Day(), tm.Day())
		assert.Equal(t, expected.Hour() + i, tm.Hour())
		assert.Equal(t, expected.Minute(), tm.Minute())
		assert.Equal(t, expected.Second(), tm.Second())
		assert.Equal(t, "UTC", tm.Location().String())
	}
}

