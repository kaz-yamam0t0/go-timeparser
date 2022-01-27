package timeparser

import "time"

const (
	AM = 1
	PM = 2
)
const (
	SET_YEAR              = 1
	SET_MONTH             = 2
	SET_DAY               = 4
	SET_HOUR              = 8
	SET_MINUTE            = 16
	SET_SECOND            = 32
	SET_MICROSECOND       = 64
	SET_TIMEZONE_OFFSET   = 128
	SET_TIMEZONE_LOCATION = 256

	SET_AP = 512

	SKIP_ERRORS = 1024
)

// Relative Additions or Subtractions
type timeAddition struct {
	n    int    // number of difference
	unit string // unit of difference

	h  int // set Hour
	i  int // set Minute
	s  int // set Second
	us int // set Microsecond

	pos     string // first | next
	month   int    // number of month
	weekday int    // number of weekday
	word    string // year | month|day
	day_flg string // empty | first | last (day of)
}

func newTimeAddition(n int, unit string) *timeAddition {
	a := timeAddition{n, unit, -1, -1, -1, -1, "", -1, -1, "", ""}
	return &a
}
func newTimeAdditionWithTime(n int, unit string, h int, i int, s int, us int) *timeAddition {
	a := timeAddition{n, unit, h, i, s, us, "", -1, -1, "", ""}
	return &a
}

// Time Data
type timeData struct {
	y         int            // Year
	m         int            // Month
	d         int            // Date
	h         int            // Hour
	i         int            // Minute
	s         int            // Second
	us        int            // Microsecond
	ap        int            // AM/PM flag (1=AM 2=PM)
	day       int            // Weekday (actually doesn't affect the result)
	z         int            // Timezone Offset
	loc       *time.Location // Timezone Location
	additions []timeAddition // Relative differences

	flags int // flags
}

// create new TimeData
func newTimeData() *timeData {
	d := timeData{1970, 1, 1, 0, 0, 0, 0, 0, 0, 0, nil, make([]timeAddition, 0), 0}
	return &d
}

func (data *timeData) appendAddition(a *timeAddition) {
	data.additions = append(data.additions, *a)
}

// resetIfUnset
func (data *timeData) resetIfUnset() {
	if !data.hasFlag(SET_YEAR) {
		data.y = 1970
	}
	if !data.hasFlag(SET_MONTH) {
		data.m = 1
	}
	if !data.hasFlag(SET_DAY) {
		data.d = 1
	}
	if !data.hasFlag(SET_HOUR) {
		data.h = 0
	}
	if !data.hasFlag(SET_MINUTE) {
		data.i = 0
	}
	if !data.hasFlag(SET_SECOND) {
		data.s = 0
	}
	if !data.hasFlag(SET_MICROSECOND) {
		data.us = 0
	}
	if !data.hasFlag(SET_TIMEZONE_OFFSET) {
		data.z = 0
	}
	if !data.hasFlag(SET_TIMEZONE_LOCATION) {
		data.loc = nil
	}
}

// reset
func (data *timeData) reset() {
	data.y = 1970
	data.m = 1
	data.d = 1
	data.h = 0
	data.i = 0
	data.s = 0
	data.us = 0
	data.z = 0
	data.loc = nil

	data.flags &= (^SET_YEAR & ^SET_MONTH & ^SET_DAY & ^SET_HOUR & ^SET_MINUTE & ^SET_SECOND & ^SET_MICROSECOND & ^SET_TIMEZONE_OFFSET & ^SET_TIMEZONE_LOCATION)
}

// ============================================================
// setter
// ============================================================

func (data *timeData) setYear(y int) {
	data.y = y
	data.flags |= SET_YEAR
}
func (data *timeData) setMonth(m int) {
	data.m = m
	data.flags |= SET_MONTH
}
func (data *timeData) setDay(d int) {
	data.d = d
	data.flags |= SET_DAY
}
func (data *timeData) setHour(h int) {
	data.h = h
	data.flags |= SET_HOUR
}
func (data *timeData) setMinute(i int) {
	data.i = i
	data.flags |= SET_MINUTE
}
func (data *timeData) setSecond(s int) {
	data.s = s
	data.flags |= SET_SECOND
}
func (data *timeData) setMicrosecond(s int) {
	data.us = s
	data.flags |= SET_MICROSECOND
}
func (data *timeData) setTimezoneOffset(z int) {
	data.z = z
	data.flags |= SET_TIMEZONE_OFFSET
}
func (data *timeData) setLocation(loc *time.Location) {
	data.loc = loc
	data.flags |= SET_TIMEZONE_LOCATION
}
func (data *timeData) setFromTime(t *time.Time) {
	data.setYear(t.Year())
	data.setMonth(int(t.Month()))
	data.setDay(t.Day())
	data.setHour(t.Hour())
	data.setMinute(t.Minute())
	data.setSecond(t.Second())
	data.setMicrosecond(int(t.Nanosecond() / 1e3))
	data.setLocation(t.Location())

	//_, offset_ := base.Zone()
	data.setTimezoneOffset(0)
}

func (data *timeData) setNow() {
	t := time.Now()
	if data.loc != nil {
		t = t.In(data.loc)
	}
	data.setFromTime(&t)
}

// ============================================================
// flags
// ============================================================
func (data *timeData) hasFlag(f int) bool {
	return (data.flags & f) == f
}

// ============================================================
// normalize
// ============================================================

// normalize Year, Month, Date, Hour, Minute, Second, Millisecond
func (data *timeData) normalize() {
	data.s, data.us = norm(data.s, data.us, 1000000)
	data.i, data.s = norm(data.i, data.s, 60)
	data.h, data.i = norm(data.h, data.i, 60)
	data.d, data.h = norm(data.d, data.h, 24)

	// day, month, year
	data.normalizeYmd()
}

// normalize Year, Month, Date
func (data *timeData) normalizeYmd() {
	// year, month
	m := data.m - 1
	data.y, m = norm(data.y, m, 12)
	data.m = m + 1

	// days
	for data.d <= 0 {
		// last month
		data.m--
		if data.m <= 0 {
			data.m = 12
			data.y--
		}
		// days of the last month
		data.d += getLastDay(data.y, data.m)
	}

	last_d := getLastDay(data.y, data.m)
	for data.d > last_d {
		data.m++
		if data.m > 12 {
			data.m = 1
			data.y++
		}
		data.d -= last_d

		last_d = getLastDay(data.y, data.m)
	}
}

// copied from time/time.go
// @see https://cs.opensource.google/go/go/+/refs/tags/go1.17.5:src/time/time.go;l=1352;drc=refs%2Ftags%2Fgo1.17.5
func norm(hi, lo, base int) (nhi, nlo int) {
	if lo < 0 {
		n := (-lo-1)/base + 1
		hi -= n
		lo += n * base
	}
	if lo >= base {
		n := lo / base
		hi += n
		lo -= n * base
	}
	return hi, lo
}

// ============================================================
// time
// ============================================================

// convert to a time.Time variable
func (data *timeData) Time() *time.Time {
	loc := time.Local
	if data.loc != nil {
		loc = data.loc
	}
	res := time.Date(data.y, time.Month(data.m), data.d, data.h, data.i, data.s, data.us*1e3, loc)

	if data.z != 0 {
		tmp := res.Unix()
		tmp += int64(data.z)
		res = time.Unix(tmp, 0).In(loc)
	}

	return &res
}

// ============================================================
// Addition
// ============================================================

func (data *timeData) processAdditions() {
	a_len := len(data.additions)
	for i := 0; i < a_len; i++ {
		data.move(&data.additions[i])
	}
	data.additions = []timeAddition{}
}

func (data *timeData) move(a *timeAddition) {
	if a.unit != "" {
		data.add(a)
		return
	}
	switch {
	case a.month > 0:
		if a.pos == "next" && data.m >= a.month {
			data.y++
		} else if a.pos == "last" && data.m <= a.month {
			data.y--
		}
		data.m = a.month
	case a.weekday >= 0:
		w_ := int(time.Date(data.y, time.Month(data.m), data.d, data.h, data.i, data.s, data.us*1e3, time.Local).Weekday())
		if a.pos == "next" {
			n := (a.weekday+7-1-w_)%7 + 1
			a.n, a.unit = n, "day"
			data.add(a)
		} else if a.pos == "last" {
			n := -((w_+7-1-a.weekday)%7 + 1)
			a.n, a.unit = n, "day"
			data.add(a)
		}
	case a.word != "":
		if a.word != "year" && a.word != "month" && a.word != "day" {
			panic("unknown word") // never be occurred
		}
		n := -1
		if a.pos == "next" {
			n = 1
		}
		a.n, a.unit = n, a.word
		data.add(a)
	}

}
func (data *timeData) add(a *timeAddition) {
	// might be outside of the range
	// it will be normalized when data is instantiated.
	switch a.unit {
	case "year":
		data.y += a.n
	case "month":
		data.m += a.n
	case "day":
		data.d += a.n
	case "hour":
		data.h += a.n
	case "min":
		fallthrough
	case "minute":
		data.i += a.n
	case "sec":
		fallthrough
	case "second":
		data.s += a.n
	case "ms":
		fallthrough
	case "msec":
		fallthrough
	case "millisecond":
		data.us += a.n * 1e3
	case "µs":
		fallthrough
	case "µsec":
		fallthrough
	case "microsecond":
		data.us += a.n
	case "week":
		data.d += a.n * 7
	case "forthnight":
		fallthrough
	case "fortnight":
		data.d += a.n * 14
	default:
		panic("unsupported units")
	}
	data.normalize()

	/*
		if a.y > 0 {
			data.y = a.y
		}
		if a.m > 0 {
			data.m = a.m
		}
		if a.d > 0 {
			data.d = a.d
		} */
	if a.h >= 0 {
		data.h = a.h
	}
	if a.i >= 0 {
		data.i = a.i
	}
	if a.s >= 0 {
		data.s = a.s
	}
	if a.us >= 0 {
		data.us = a.us
	}
	data.normalize()
}
