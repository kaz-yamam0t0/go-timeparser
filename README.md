# go-timeparser
Flexible Time Parser for Golang

## Installation

Download `timeparser` with the following command:

```shell
$ go get -u github.com/kaz-yamam0t0/go-timeparser/timeparser
```

Then import this library in your source code:

```go
import "github.com/kaz-yamam0t0/go-timeparser/timeparser"
```

## Usage

### Simply Parse or Format time.Time variable

```go
// 
// time.Time vars to formatted strings
// 
tm := time.Now()
fmt.Println(timeparser.FormatTime("r", &tm))                     // Wed, 29 Dec 2021 18:24:00 +0900
fmt.Println(timeparser.FormatTime("l jS \\of F Y h:i:s A", &tm)) // Wednesday 29th of December 2021 06:24:00 PM
fmt.Println(timeparser.FormatTime("Y-m-d H:i:s", &tm))           // 2021-12-29 18:24:00

// 
// parse strings into time.Time vars
// 
tm, err := timeparser.ParseFormat(" l jS \\of F Y  h:i:s A", "Wednesday 29th of December  2021 06:24:12 PM ")
if err != nil {
	panic(err)
}
fmt.Println(tm)

// 
// parse strings more flexibly
// 
tm, err := timeparser.ParseTimeStr("Wed, 29 Dec 2021 18:24:00 +0900", nil)
tm, err := timeparser.ParseTimeStr("2021-12-29T18:24:00+09:00", nil)
tm, err := timeparser.ParseTimeStr("2021-12-29T18:24:00Z", nil)
tm, err := timeparser.ParseTimeStr("Wednesday 29th December 2021 06:24:00 PM", nil)

// 
// relative format
// 
tm := time.Now()
tm, err := timeparser.ParseTimeStr("1 year ago", &tm)
tm, err := timeparser.ParseTimeStr("+1 day 2 week 3 months -4 years 5 hours -6 minutes 7 seconds", &tm)
tm, err := timeparser.ParseTimeStr("P1Y2M3DT4H5M6.123456S", &tm) // ISO 8601 Interval Format

tm, err := timeparser.ParseTimeStr("next Thursday", &tm)
tm, err := timeparser.ParseTimeStr("last year", &tm)

```

### TimeData

TimeData is a simple and flexible struct used in `timeparser`.

```go
// tdata, _ := timeparser.Now()
// tdata, _ := timeparser.New("2022-01-31 11:22:33.123456789")
tdata, _ := timeparser.NewAsUTC("2022-01-31 11:22:33.123456789")

// 
// If you want to adjust the timezone offset difference from your local env, 
// first create a variable with `timeparser.New()` and then convert it to UTC with `SetUTC()`
// 
// tdata, _ := timeparser.New("2022-01-31 11:22:33.123456789")
// tdata.SetUTC()

fmt.Println(tdata.GetYear())        // 2022
fmt.Println(tdata.GetMonth())       // 1
fmt.Println(tdata.GetDay())         // 31
fmt.Println(tdata.GetHour())        // 11
fmt.Println(tdata.GetMinute())      // 22
fmt.Println(tdata.GetSecond())      // 33
fmt.Println(tdata.GetMillisecond()) // 123
fmt.Println(tdata.GetMicrosecond()) // 123456
fmt.Println(tdata.GetNanosecond())  // 123456789

// 
// Format
// 
fmt.Println(tdata.String()) // 2022-01-31T11:22:33+00:00
fmt.Println(tdata.Format("l jS \\of F Y h:i:s A")) // Monday 31st of January 2022 06:22:33 PM

// 
// Setter methods
// 
tdata.SetYear(2022)
tdata.SetMonth(1)
tdata.SetDay(31)
tdata.SetHour(11)
tdata.SetMinute(22)
tdata.SetSecond(33)

// 
// Difference 
// 
tm, _ := timeparser.NewAsUTC("2020-12-31 11:22:33.123456789")
fmt.Println(tdata.DiffYears(tm))   // 1
fmt.Println(tdata.DiffMonths(tm))  // 13
fmt.Println(tdata.DiffDays(tm))    // 396
fmt.Println(tdata.DiffHours(tm))   // 9504
fmt.Println(tdata.DiffMinutes(tm)) // 570240
fmt.Println(tdata.DiffSeconds(tm)) // 34214400

```


## Documentation

`go doc` style documentation:

https://pkg.go.dev/github.com/kaz-yamam0t0/go-timeparser/timeparser

Most format characters are based on PHP datetime related functions.

https://www.php.net/manual/en/datetime.createfromformat.php
