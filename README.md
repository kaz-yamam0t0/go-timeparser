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
tm, err := timeparser.ParseTimeStr("Wednesday 29th December 2021 06:24:00 PM", nil)

// 
// relative format
// 
tm := time.Now()
fmt.Println(timeparser.ParseTimeStr("1 year ago", &tm))
fmt.Println(timeparser.ParseTimeStr("+1 day 2 week 3 months -4 years 5 hours -6 minutes 7 seconds", &tm))
fmt.Println(timeparser.ParseTimeStr("P1Y2M3DT4H5M6.123456S", &tm)) // ISO 8601 Interval Format

fmt.Println(timeparser.ParseTimeStr("next Thursday", &tm))
fmt.Println(timeparser.ParseTimeStr("last year", &tm))

```

## Documentation

`go doc` style documentation:

https://pkg.go.dev/github.com/kaz-yamam0t0/go-timeparser/timeparser


