package timeparser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseTimeStr(t *testing.T) {
	utc, err := time.LoadLocation("UTC")
	assert.Nil(t, err)

	base := time.Date(2000, time.September, 10, 0, 0, 0, 0, time.Local)

	testcases := map[string]time.Time{
		// General Format
		"10 September 2000":   time.Date(2000, time.September, 10, 0, 0, 0, 0, time.Local),
		"2000/09/10":          time.Date(2000, time.September, 10, 0, 0, 0, 0, time.Local),
		"2000-09-10":          time.Date(2000, time.September, 10, 0, 0, 0, 0, time.Local),
		"2000-09-10 12:34:56": time.Date(2000, time.September, 10, 12, 34, 56, 0, time.Local),

		"Wed, 29 Dec 2021 18:24:00 +0000":          time.Date(2021, time.December, 29, 18, 24, 0, 0, utc),
		"Wed, 29 Dec 2021 18:24:00 +0900":          time.Date(2021, time.December, 30, 3, 24, 0, 0, utc),
		"Wed, 29 Dec 2021 18:24:00 +08:00":         time.Date(2021, time.December, 30, 2, 24, 0, 0, utc),
		"2021-12-29T18:24:00+09:00":                time.Date(2021, time.December, 30, 3, 24, 0, 0, utc),
		"Wednesday 29th December 2021 06:24:00 PM": time.Date(2021, time.December, 29, 18, 24, 0, 0, time.Local),

		"2021-12-29T18:24:00Z":      time.Date(2021, time.December, 29, 18, 24, 0, 0, utc),
		"2021-12-29T18:24:00Z09:00": time.Date(2021, time.December, 30, 3, 24, 0, 0, utc),

		// American or European format (sometimes this may be Ambiguious)
		"9/10/2000": time.Date(2000, time.September, 10, 0, 0, 0, 0, time.Local),
		"10.9.2000": time.Date(2000, time.September, 10, 0, 0, 0, 0, time.Local),
		"10-9-2000": time.Date(2000, time.September, 10, 0, 0, 0, 0, time.Local),

		// Relative format
		"+0 day":        time.Date(2000, time.September, 10, 0, 0, 0, 0, time.Local),
		"+1 day":        time.Date(2000, time.September, 11, 0, 0, 0, 0, time.Local),
		"+1 day 1 week": time.Date(2000, time.September, 18, 0, 0, 0, 0, time.Local),
		"+1 day 2 week 3 months -4 years 5 hours -6 minutes 7 seconds": time.Date(1996, time.December, 25, 4, 54, 7, 0, time.Local),
		"-1year -13months -8weeks":                                     time.Date(1998, time.June, 15, 0, 0, 0, 0, time.Local),
		"1 year ago 13months ago 8 weeks ago":                          time.Date(1998, time.June, 15, 0, 0, 0, 0, time.Local),

		// ISO8601 Interval format
		"P1D":              time.Date(2000, time.September, 11, 0, 0, 0, 0, time.Local),
		"P1Y2M3D":          time.Date(2001, time.November, 13, 0, 0, 0, 0, time.Local),
		"P1Y2M3DT4H5M6S":   time.Date(2001, time.November, 13, 4, 5, 6, 0, time.Local),
		"P1Y2M3DT4H5M6.7S": time.Date(2001, time.November, 13, 4, 5, 6, 7e8, time.Local),

		// Other relative format
		"next Thursday": time.Date(2000, time.September, 14, 0, 0, 0, 0, time.Local),
		"last Monday":   time.Date(2000, time.September, 4, 0, 0, 0, 0, time.Local),
	}
	for format, expected := range testcases {
		tm, err := ParseTimeStr(format, &base)

		if err != nil {
			panic(err)
		}

		assert.Nil(t, err)
		assert.Equal(t, expected.Year(), tm.Year())
		assert.Equal(t, expected.Month(), tm.Month())
		assert.Equal(t, expected.Day(), tm.Day())
		assert.Equal(t, expected.Hour(), tm.Hour())
		assert.Equal(t, expected.Minute(), tm.Minute())
		assert.Equal(t, expected.Second(), tm.Second())
		assert.Equal(t, expected.Location(), tm.Location())
	}
}
func ExampleParseTimeStr() {
	// Strtotime(format string) returns int64
	// or -1 when an error has occurred.

	tm := time.Now()

	// general format
	fmt.Println(ParseTimeStr("10 September 2000", nil))
	fmt.Println(ParseTimeStr("2000/09/10", nil))
	fmt.Println(ParseTimeStr("2000-09-10", nil))
	fmt.Println(ParseTimeStr("2000-09-10 12:34:56", nil))

	// American or European format (Ambiguous)
	fmt.Println(ParseTimeStr("9/10/2000", nil))
	fmt.Println(ParseTimeStr("10.9.2000", nil))
	fmt.Println(ParseTimeStr("10-9-2000", nil))

	// various formats
	fmt.Println(ParseTimeStr("Wed, 29 Dec 2021 18:24:00 +0900", nil))
	fmt.Println(ParseTimeStr("2021-12-29T18:24:00+09:00", nil))
	fmt.Println(ParseTimeStr("Wednesday 29th December 2021 06:24:00 PM", nil))

	// relative format
	fmt.Println(ParseTimeStr("+0 day", &tm))
	fmt.Println(ParseTimeStr("+1 day 2 week 3 months -4 years 5 hours -6 minutes 7 seconds", &tm))
	fmt.Println(ParseTimeStr("-1year -13months -8week", &tm))
	fmt.Println(ParseTimeStr("1 year ago", &tm))
	fmt.Println(ParseTimeStr("P1Y2M3DT4H5M6.123456S", &tm)) // ISO 8601 Interval Format

	fmt.Println(ParseTimeStr("next Thursday", &tm))
	fmt.Println(ParseTimeStr("last Monday", &tm))
}
