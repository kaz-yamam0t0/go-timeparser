package timeparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)


func TestFormatTime(t *testing.T) {
	utc, err := time.LoadLocation("UTC")
	assert.Nil(t, err)

	tm := time.Date(2021, time.December, 29, 18, 24, 36, 123456789, utc)

	testcases := map[string]string {
		// general format
		"r" : "Wed, 29 Dec 2021 18:24:36 +0000",
		"c" : "2021-12-29T18:24:36+00:00",

		"l jS \\of F Y h:i:s A" : "Wednesday 29th of December 2021 06:24:36 PM",
		"Y-m-d H:i:s" : "2021-12-29 18:24:36",
		"n/j/y" : "12/29/21",
		"W" : "52",

		// Each format
		"D l w N" : "Wed Wednesday 3 3" ,
		"F m M n t" : "December 12 Dec 12 31" , 
		"Y y L o" : "2021 21 0 2021",
		"a A B" : "pm PM 433",
		"g G h H" : "6 18 06 18",
		"i s u v" : "24 36 123456 123",
		"Z P O" : "0 +00:00 +0000",
		"e T" : "UTC UTC",
		"U" : "1640802276" , 
	}
	for format, expected := range testcases {
		assert.Equal(t, expected, FormatTime(format, &tm))
	}
}
