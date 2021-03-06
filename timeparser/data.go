package timeparser

import "time"

//import "fmt"

// Internal Constants
const (
	AM = 1
	PM = 2

	SET_YEAR              = 1
	SET_MONTH             = 2
	SET_DAY               = 4
	SET_HOUR              = 8
	SET_MINUTE            = 16
	SET_SECOND            = 32
	SET_NANOSECOND        = 64
	SET_TIMEZONE_OFFSET   = 128
	SET_TIMEZONE_LOCATION = 256
	SET_AP                = 512
	SKIP_ERRORS           = 1024
)

// ============================================================
// TimeAddition
// ============================================================

// Relative Additions or Subtractions
type timeAddition struct {
	n    int    // number of difference
	unit string // unit of difference

	h  int // set Hour
	i  int // set Minute
	s  int // set Second
	ns int // set Nanosecond

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

// ============================================================
// TimeData
// ============================================================

type TimeData struct {
	y         int            // Year
	m         int            // Month
	d         int            // Date
	h         int            // Hour
	i         int            // Minute
	s         int            // Second
	ns        int            // Nanosecond
	ap        int            // AM/PM flag (1=AM 2=PM)
	day       int            // Weekday (actually doesn't affect the result)
	z         int            // Timezone Offset
	loc       *time.Location // Timezone Location
	additions []timeAddition // Relative differences

	flags int // flags
}

// create a new TimeData variable of 1970/01/01
func newTimeData() *TimeData {
	d := TimeData{1970, 1, 1, 0, 0, 0, 0, 0, 0, 0, nil, make([]timeAddition, 0), 0}
	return &d
}

func (data *TimeData) appendAddition(a *timeAddition) {
	data.additions = append(data.additions, *a)
}

// resetIfUnset
func (data *TimeData) resetIfUnset() {
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
	if !data.hasFlag(SET_NANOSECOND) {
		data.ns = 0
	}
	if !data.hasFlag(SET_TIMEZONE_OFFSET) {
		data.z = 0
	}
	if !data.hasFlag(SET_TIMEZONE_LOCATION) {
		data.loc = nil
	}
}

// reset
func (data *TimeData) reset() {
	data.y = 1970
	data.m = 1
	data.d = 1
	data.h = 0
	data.i = 0
	data.s = 0
	data.ns = 0
	data.z = 0
	data.loc = nil

	data.flags &= (^SET_YEAR & ^SET_MONTH & ^SET_DAY & ^SET_HOUR & ^SET_MINUTE & ^SET_SECOND & ^SET_NANOSECOND & ^SET_TIMEZONE_OFFSET & ^SET_TIMEZONE_LOCATION & ^SET_AP)
}

// ============================================================
// public constructors
// ============================================================

// create a new TimeData variable from string format
func New(format string) (*TimeData, error) {
	data, err := parseTimeStr(format, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// create a new TimeData variable from string format
func NewAsUTC(format string) (*TimeData, error) {
	data, err := New(format)
	if err != nil {
		return nil, err
	}
	return data.AsUTC(), nil
}

// create a new TimeData variable of Now
func Now() *TimeData {
	data := newTimeData()
	data.setNow()
	return data
}

// ============================================================
// setter
// ============================================================

// private set* methods doesn't call normalize(), so you have to call normalize() if it is not sure whether argument values are normalized.
//
// public methods always call normalize().

func (data *TimeData) setYear(y int) {
	data.y = y
	data.flags |= SET_YEAR
}
func (data *TimeData) setMonth(m int) {
	data.m = m
	data.flags |= SET_MONTH
}
func (data *TimeData) setDay(d int) {
	data.d = d
	data.flags |= SET_DAY
}
func (data *TimeData) setHour(h int) {
	data.h = h
	data.flags |= SET_HOUR
}
func (data *TimeData) setMinute(i int) {
	data.i = i
	data.flags |= SET_MINUTE
}
func (data *TimeData) setSecond(s int) {
	data.s = s
	data.flags |= SET_SECOND
}
func (data *TimeData) setNanosecond(ns int) {
	data.ns = ns
	data.flags |= SET_NANOSECOND
}
func (data *TimeData) setTimezoneOffset(z int) {
	data.z = z
	data.flags |= SET_TIMEZONE_OFFSET
}
func (data *TimeData) setLocation(loc *time.Location) {
	data.loc = loc
	data.flags |= SET_TIMEZONE_LOCATION
}
func (data *TimeData) setUTC() {
	_, offset_ := data.Time().Zone()

	utc, _ := time.LoadLocation("UTC")
	data.setLocation(utc)
	if data.z != 0 {
		z := data.z

		data.setTimezoneOffset(0)
		data.AddSecond(-z)
	}
	if offset_ != 0 {
		data.AddSecond(-offset_)
	}

}

func (data *TimeData) setFromTime(t *time.Time) {
	data.setYear(t.Year())
	data.setMonth(int(t.Month()))
	data.setDay(t.Day())
	data.setHour(t.Hour())
	data.setMinute(t.Minute())
	data.setSecond(t.Second())
	data.setNanosecond(t.Nanosecond())
	data.setLocation(t.Location())

	//_, offset_ := base.Zone()
	data.setTimezoneOffset(0)
}

func (data *TimeData) setNow() {
	t := time.Now()
	if data.loc != nil {
		t = t.In(data.loc)
	}
	data.setFromTime(&t)
}

// public methods always call normalize() or normalizeYmd().
// They are more safe but have less performance.

func (data *TimeData) SetYear(y int) {
	data.setYear(y)
	data.normalizeYmd()
}
func (data *TimeData) SetMonth(m int) {
	y := data.y
	d := data.d

	// normalize year, month
	m--
	y, m = norm(data.y, m, 12)
	m++

	// normalize day
	last_d := getLastDay(y, m)
	if last_d < d {
		// data.setDay(last_d)
		d = last_d
	}
	data.setMonth(m)
	data.y = y // don't modify flags
	data.d = d // don't modify flags

	// normalize
	// data.normalizeYmd()
}
func (data *TimeData) SetDay(d int) {
	data.setDay(d)
	data.normalizeYmd()
}
func (data *TimeData) SetHour(h int) {
	data.setHour(h)
	data.normalize()
}
func (data *TimeData) SetMinute(i int) {
	data.setMinute(i)
	data.normalize()
}
func (data *TimeData) SetSecond(s int) {
	data.setSecond(s)
	data.normalize()
}
func (data *TimeData) SetMillisecond(ms int) {
	data.setNanosecond(ms * 1e6)
	data.normalize()
}
func (data *TimeData) SetMicrosecond(us int) {
	data.setNanosecond(us * 1e3)
	data.normalize()
}
func (data *TimeData) SetNanosecond(ns int) {
	data.setNanosecond(ns)
	data.normalize()
}
func (data *TimeData) SetTimezoneOffset(z int) {
	data.setUTC()
	data.setTimezoneOffset(z)
}
func (data *TimeData) SetLocation(loc *time.Location) {
	data.setLocation(loc)
}
func (data *TimeData) SetUTC() {
	data.setUTC()
}

func (data *TimeData) SetFromTime(t *time.Time) {
	data.setFromTime(t)
}
func (data *TimeData) SetNow() {
	data.setNow()
}

// this function is same as `SetLocation` but returns the pointer of the TimeData variable
func (data *TimeData) In(loc *time.Location) *TimeData {
	data.SetLocation(loc)
	return data
}

// this function is same as `SetUTC` but returns the pointer of the TimeData variable
func (data *TimeData) AsUTC() *TimeData {
	data.SetUTC()
	return data
}

// ============================================================
// getter
// ============================================================

func (data *TimeData) GetYear() int {
	return data.y
}
func (data *TimeData) GetMonth() int {
	return data.m
}
func (data *TimeData) GetDay() int {
	return data.d
}
func (data *TimeData) GetHour() int {
	return data.h
}
func (data *TimeData) GetMinute() int {
	return data.i
}
func (data *TimeData) GetSecond() int {
	return data.s
}
func (data *TimeData) GetMillisecond() int {
	return int(data.ns / 1e6)
}
func (data *TimeData) GetMicrosecond() int {
	return int(data.ns / 1e3)
}
func (data *TimeData) GetNanosecond() int {
	return data.ns
}
func (data *TimeData) GetTimezoneOffset() int {
	return data.z
}
func (data *TimeData) GetLocation() *time.Location {
	return data.loc
}

// ============================================================
// Add / Sub
// ============================================================

func (data *TimeData) AddYear(y int) {
	data.SetYear(data.y + y)
}
func (data *TimeData) AddMonth(m int) {
	data.SetMonth(data.m + m)
}
func (data *TimeData) AddDay(d int) {
	data.SetDay(data.d + d)
}
func (data *TimeData) AddHour(h int) {
	data.SetHour(data.h + h)
}
func (data *TimeData) AddMinute(i int) {
	data.SetMinute(data.i + i)
}
func (data *TimeData) AddSecond(s int) {
	data.SetSecond(data.s + s)
}
func (data *TimeData) AddMillisecond(ms int) {
	data.SetNanosecond(data.ns + ms*1e6)
}
func (data *TimeData) AddMicrosecond(us int) {
	data.SetNanosecond(data.ns + us*1e3)
}
func (data *TimeData) AddNanosecond(ns int) {
	data.SetNanosecond(data.ns + ns)
}

//func (data *TimeData) AddTimezoneOffset(z int) {
//	data.setTimezoneOffset(data.z + z)
//}
func (data *TimeData) SubYear(y int) {
	data.AddYear(-y)
}
func (data *TimeData) SubMonth(m int) {
	data.AddMonth(-m)
}
func (data *TimeData) SubDay(d int) {
	data.AddDay(-d)
}
func (data *TimeData) SubHour(h int) {
	data.AddHour(-h)
}
func (data *TimeData) SubMinute(i int) {
	data.AddMinute(-i)
}
func (data *TimeData) SubSecond(s int) {
	data.AddSecond(-s)
}
func (data *TimeData) SubMillisecond(ms int) {
	data.AddNanosecond(-ms * 1e6)
}
func (data *TimeData) SubMicrosecond(us int) {
	data.AddNanosecond(-us * 1e3)
}
func (data *TimeData) SubNanosecond(ns int) {
	data.AddNanosecond(-ns)
}

// ============================================================
// Diff
// ============================================================

func (data *TimeData) DiffYears(d *TimeData) int {
	return int(data.DiffMonths(d) / 12)
}
func (data *TimeData) DiffMonths(d *TimeData) int {
	y_ := data.y - d.y
	m_ := data.m - d.m
	d_ := data.d - d.d

	m_ += y_ * 12

	if m_ > 0 && d_ < 0 {
		m_--
	} else if m_ < 0 && d_ > 0 {
		m_++
	}
	return m_
}
func (data *TimeData) DiffDays(d *TimeData) int {
	return int(data.DiffSeconds(d) / 86400)
}
func (data *TimeData) DiffHours(d *TimeData) int {
	return int(data.DiffSeconds(d) / 3600)
}
func (data *TimeData) DiffMinutes(d *TimeData) int {
	return int(data.DiffSeconds(d) / 60)
}
func (data *TimeData) DiffSeconds(d *TimeData) int64 {
	return data.Unix() - d.Unix()
}
func (data *TimeData) DiffMilliseconds(d *TimeData) int64 {
	return data.UnixMilli() - d.UnixMilli()
}
func (data *TimeData) DiffMicroseconds(d *TimeData) int64 {
	return data.UnixMicro() - d.UnixMicro()
}
func (data *TimeData) DiffNanoseconds(d *TimeData) int64 {
	return data.UnixNano() - d.UnixNano()
}

// ============================================================
// flags
// ============================================================
func (data *TimeData) hasFlag(f int) bool {
	return (data.flags & f) == f
}

// ============================================================
// normalize
// ============================================================

// normalize Year, Month, Date, Hour, Minute, Second, Millisecond
func (data *TimeData) normalize() {
	data.s, data.ns = norm(data.s, data.ns, 1e9)
	data.i, data.s = norm(data.i, data.s, 60)
	data.h, data.i = norm(data.h, data.i, 60)
	data.d, data.h = norm(data.d, data.h, 24)

	// day, month, year
	data.normalizeYmd()
}

// normalize Year, Month, Date
func (data *TimeData) normalizeYmd() {
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
func (data *TimeData) Time() *time.Time {
	loc := time.Local
	if data.loc != nil {
		loc = data.loc
	}
	res := time.Date(data.y, time.Month(data.m), data.d, data.h, data.i, data.s, data.ns, loc)

	if data.z != 0 {
		tmp := res.Unix()
		tmp += int64(data.z)
		res = time.Unix(tmp, 0).In(loc)
	}

	return &res
}

func (data *TimeData) Unix() int64 {
	return data.Time().Unix()
}
func (data *TimeData) UnixMilli() int64 {
	return data.Time().UnixMilli()
}
func (data *TimeData) UnixMicro() int64 {
	return data.Time().UnixMicro()
}
func (data *TimeData) UnixNano() int64 {
	return data.Time().UnixNano()
}

// ============================================================
// Format
// ============================================================

func (data *TimeData) Format(s string) string {
	return FormatTime(s, data.Time())
}
func (data *TimeData) String() string {
	return data.Format("c")
}

// ============================================================
// Addition
// ============================================================

func (data *TimeData) processAdditions() {
	a_len := len(data.additions)
	for i := 0; i < a_len; i++ {
		data.move(&data.additions[i])
	}
	data.additions = []timeAddition{}
}

func (data *TimeData) move(a *timeAddition) {
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
		w_ := int(time.Date(data.y, time.Month(data.m), data.d, data.h, data.i, data.s, data.ns, time.Local).Weekday())
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
func (data *TimeData) add(a *timeAddition) {
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
		data.ns += a.n * 1e6
	case "??s":
		fallthrough
	case "??sec":
		fallthrough
	case "microsecond":
		data.ns += a.n * 1e3
	case "ns":
		fallthrough
	case "nsec":
		fallthrough
	case "nanosecond":
		data.ns += a.n
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
	if a.ns >= 0 {
		data.ns = a.ns
	}
	data.normalize()
}
