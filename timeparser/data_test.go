package timeparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"fmt"
)


func TestTimeData(t *testing.T) {
	tdata, err := New("2022-01-31 11:22:33.123456789")

	assert.Nil(t, err)

	// GET
	assert.Equal(t, 2022, tdata.GetYear())
	assert.Equal(t, 1, tdata.GetMonth())
	assert.Equal(t, 31, tdata.GetDay())
	assert.Equal(t, 11, tdata.GetHour())
	assert.Equal(t, 22, tdata.GetMinute())
	assert.Equal(t, 33, tdata.GetSecond())
	assert.Equal(t, 123, tdata.GetMillisecond())
	assert.Equal(t, 123456, tdata.GetMicrosecond())
	assert.Equal(t, 123456789, tdata.GetNanosecond())

	// SET
	tdata.SetYear(2030)
	tdata.SetMonth(2)
	tdata.SetDay(3)
	tdata.SetHour(4)
	tdata.SetMinute(5)
	tdata.SetSecond(6)

	assert.Equal(t, 2030, tdata.GetYear())
	assert.Equal(t, 2, tdata.GetMonth())
	assert.Equal(t, 3, tdata.GetDay())
	assert.Equal(t, 4, tdata.GetHour())
	assert.Equal(t, 5, tdata.GetMinute())
	assert.Equal(t, 6, tdata.GetSecond())

	tdata.SetMillisecond(123)
	assert.Equal(t, 123, tdata.GetMillisecond())
	assert.Equal(t, 123000, tdata.GetMicrosecond())
	assert.Equal(t, 123000000, tdata.GetNanosecond())

	tdata.SetMicrosecond(2456)
	assert.Equal(t, 2, tdata.GetMillisecond())
	assert.Equal(t, 2456, tdata.GetMicrosecond())
	assert.Equal(t, 2456000, tdata.GetNanosecond())

	tdata.SetNanosecond(1789)
	assert.Equal(t, 0, tdata.GetMillisecond())
	assert.Equal(t, 1, tdata.GetMicrosecond())
	assert.Equal(t, 1789, tdata.GetNanosecond())
}
func TestAdd(t *testing.T) {
	// Format
	tdata, _ := New("2022-01-31 18:22:33.123456789 +0000")

	// Add
	tdata.AddYear(10)
	assert.Equal(t, 2032, tdata.GetYear())

	tdata.AddMonth(10)
	assert.Equal(t, 11, tdata.GetMonth())

	tdata.AddDay(31)
	assert.Equal(t, 12, tdata.GetMonth())
	assert.Equal(t, 31, tdata.GetDay())

	tdata.AddHour(12)
	assert.Equal(t, 1, tdata.GetMonth())
	assert.Equal(t, 1, tdata.GetDay())
	assert.Equal(t, 6, tdata.GetHour())

	tdata.AddMinute(50)
	assert.Equal(t, 7, tdata.GetHour())
	assert.Equal(t, 12, tdata.GetMinute())

	tdata.AddSecond(50)
	assert.Equal(t, 7, tdata.GetHour())
	assert.Equal(t, 13, tdata.GetMinute())
	assert.Equal(t, 23, tdata.GetSecond())

	tdata.AddMillisecond(999)
	assert.Equal(t, 24, tdata.GetSecond())
	assert.Equal(t, 122456789, tdata.GetNanosecond())

	tdata.AddMicrosecond(999999)
	assert.Equal(t, 25, tdata.GetSecond())
	assert.Equal(t, 122455789, tdata.GetNanosecond())

	tdata.AddNanosecond(999999999)
	assert.Equal(t, 26, tdata.GetSecond())
	assert.Equal(t, 122455788, tdata.GetNanosecond())
}
func TestSub(t *testing.T) {
	// Format
	tdata, _ := New("2022-01-31 18:22:33.123456789 +0000")

	// Add
	tdata.SubYear(10)
	assert.Equal(t, 2012, tdata.GetYear())

	//tdata.SubMonth(11)
	tdata.SetMonth(2)
	assert.Equal(t, 2, tdata.GetMonth())
	assert.Equal(t, 29, tdata.GetDay()) // last day

	tdata.SubDay(32)
	assert.Equal(t, 1, tdata.GetMonth())
	assert.Equal(t, 28, tdata.GetDay())

	tdata.SubHour(20)
	assert.Equal(t, 1, tdata.GetMonth())
	assert.Equal(t, 27, tdata.GetDay())
	assert.Equal(t, 22, tdata.GetHour())

	tdata.SubMinute(50)
	assert.Equal(t, 21, tdata.GetHour())
	assert.Equal(t, 32, tdata.GetMinute())

	tdata.SubSecond(50)
	assert.Equal(t, 21, tdata.GetHour())
	assert.Equal(t, 31, tdata.GetMinute())
	assert.Equal(t, 43, tdata.GetSecond())

	tdata.SubMillisecond(999)
	assert.Equal(t, 42, tdata.GetSecond())
	assert.Equal(t, 124456789, tdata.GetNanosecond())

	tdata.SubMicrosecond(999999)
	assert.Equal(t, 41, tdata.GetSecond())
	assert.Equal(t, 124457789, tdata.GetNanosecond())

	tdata.SubNanosecond(999999999)
	assert.Equal(t, 40, tdata.GetSecond())
	assert.Equal(t, 124457790, tdata.GetNanosecond())
}


//func TestDiff(t *testing.T) {
//	tdata, _ := New("2022-08-15 18:22:33.123456789 +0000")
//	tdata2, _ := New("2023-02-12 12:24:36.123456789 +0000")
//}
func TestDiffYears(t *testing.T) {
	// Format
	tdata, _ := New("2022-01-15 18:22:33.123456789 +0000")

	testcases := map[string]int {
		// 2020-12-*
		"2020-12-14 18:22:33.123456789 +0000": 1,
		"2020-12-15 18:22:33.123456789 +0000": 1,
		"2020-12-16 18:22:33.123456789 +0000": 1,

		// 2021-01-*
		"2021-01-14 18:22:33.123456789 +0000": 1,
		"2021-01-15 18:22:33.123456789 +0000": 1,
		"2021-01-16 18:22:33.123456789 +0000": 0,

		// 2022-01-*
		"2022-01-14 18:22:33.123456789 +0000": 0,
		"2022-01-15 18:22:33.123456789 +0000": 0,
		"2022-01-16 18:22:33.123456789 +0000": 0,

		// 2022-02-*
		"2022-02-14 18:22:33.123456789 +0000": 0,
		"2022-02-15 18:22:33.123456789 +0000": 0,
		"2022-02-16 18:22:33.123456789 +0000": 0,

		// 2023-01-*
		"2023-01-14 18:22:33.123456789 +0000": 0,
		"2023-01-15 18:22:33.123456789 +0000": -1,
		"2023-01-16 18:22:33.123456789 +0000": -1,
	}
	for format, v := range testcases {
		td, _ := New(format)
		assert.Equal(t, v, tdata.DiffYears(td))
	}
}

func TestDiffMonths(t *testing.T) {
	// Format
	tdata, _ := New("2022-01-15 18:22:33.123456789 +0000")

	testcases := map[string]int {
		// 2020-12-*
		"2020-12-14 18:22:33.123456789 +0000": 13,
		"2020-12-15 18:22:33.123456789 +0000": 13,
		"2020-12-16 18:22:33.123456789 +0000": 12,

		// 2021-01-*
		"2021-01-14 18:22:33.123456789 +0000": 12,
		"2021-01-15 18:22:33.123456789 +0000": 12,
		"2021-01-16 18:22:33.123456789 +0000": 11,

		// 2022-01-*
		"2022-01-14 18:22:33.123456789 +0000": 0,
		"2022-01-15 18:22:33.123456789 +0000": 0,
		"2022-01-16 18:22:33.123456789 +0000": 0,

		// 2022-02-*
		"2022-02-14 18:22:33.123456789 +0000": 0,
		"2022-02-15 18:22:33.123456789 +0000": -1,
		"2022-02-16 18:22:33.123456789 +0000": -1,

		// 2023-01-*
		"2023-01-14 18:22:33.123456789 +0000": -11,
		"2023-01-15 18:22:33.123456789 +0000": -12,
		"2023-01-16 18:22:33.123456789 +0000": -12,
	}
	for format, v := range testcases {
		td, _ := New(format)
		assert.Equal(t, v, tdata.DiffMonths(td))
	}
}
func TestFormat(t *testing.T) {
	// Format
	tdata, _ := New("2022-01-31 18:22:33.123456789 +0000")
	testcases := map[string]string {
		// general format
		"r" : "Mon, 31 Jan 2022 18:22:33 +0000",
		"c" : "2022-01-31T18:22:33+00:00",

		"l jS \\of F Y h:i:s A" : "Monday 31st of January 2022 06:22:33 PM",
		"Y-m-d H:i:s" : "2022-01-31 18:22:33",
		"n/j/y" : "1/31/22",
		"W" : "5",

		// Each format
		"D l w N" : "Mon Monday 1 1" ,
		"F m M n t" : "January 01 Jan 1 31" , 
		"Y y L o" : "2022 22 0 2022",
		"a A B" : "pm PM 432",
		"g G h H" : "6 18 06 18",
		"i s u v" : "22 33 123456 123",
		"Z P O" : "0 +00:00 +0000",
		"e T" : "UTC UTC",
		"U" : "1643653353" , 
	}
	for format, expected := range testcases {
		assert.Equal(t, expected, tdata.Format(format))
	}
}


func TestTimezoneOffset(t *testing.T) {
	// TimeZone Offset always modifys Location to UTC
	tdata, _ := New("2022-01-31 11:22:33.123456789 +01:00")
	utc, _ := time.LoadLocation("UTC")

	assert.Equal(t, utc, tdata.GetLocation())
	assert.Equal(t, 3600, tdata.GetTimezoneOffset())

	tdata, _ = New("2022-01-31 11:22:33.123456789 -0030")
	assert.Equal(t, -1800, tdata.GetTimezoneOffset())

	// When location is set, GetTimeZoneOffset() is always 0.
	// Note that GetTimeZoneOffset() doesn't mean the offset of the location 
	// and it should be always used with UTC Timezone.
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	
	var err error
	tdata, err = New("2022-01-31 11:22:33.123456789 Asia/Tokyo")
	assert.Nil(t, err)
	assert.Equal(t, tokyo, tdata.GetLocation())
	assert.Equal(t, 0, tdata.GetTimezoneOffset())

	// SetTimezoneOffset() will set not only timezone but also Location to "UTC".
	tdata.SetTimezoneOffset(-5 * 60 * 60)
	assert.Equal(t, utc, tdata.GetLocation())
	assert.Equal(t, -5 * 60 * 60, tdata.GetTimezoneOffset())

	// you can also define UTC TimeData like this
	tdata, err = NewAsUTC("2022-01-31 11:22:33.123456789 Asia/Tokyo")
	//tdata, err = Now().AsUTC()
	assert.Equal(t, utc, tdata.GetLocation())
	assert.Equal(t, 0, tdata.GetTimezoneOffset())
	assert.Equal(t, 11 - 9, tdata.GetHour()) // Asia/Tokyo offset is +9:00
}

// ============================================================
// examples
// ============================================================

func ExampleNew() {
	tdata, _ := New("2022-01-31 11:22:33.123456789")

	fmt.Println(tdata.GetYear())        // 2022
	fmt.Println(tdata.GetMonth())       // 1
	fmt.Println(tdata.GetDay())         // 31
	fmt.Println(tdata.GetHour())        // 11
	fmt.Println(tdata.GetMinute())      // 22
	fmt.Println(tdata.GetSecond())      // 333
	fmt.Println(tdata.GetMillisecond()) // 123
	fmt.Println(tdata.GetMicrosecond()) // 123456
	fmt.Println(tdata.GetNanosecond())  // 123456789
}
func ExampleNow() {
	tdata := Now()

	fmt.Println(tdata.GetYear())        // 2022
	fmt.Println(tdata.GetMonth())       // 1
	fmt.Println(tdata.GetDay())         // 31
	fmt.Println(tdata.GetHour())        // 11
	fmt.Println(tdata.GetMinute())      // 22
	fmt.Println(tdata.GetSecond())      // 333
	fmt.Println(tdata.GetMillisecond()) // 123
	fmt.Println(tdata.GetMicrosecond()) // 123456
	fmt.Println(tdata.GetNanosecond())  // 123456789
}
