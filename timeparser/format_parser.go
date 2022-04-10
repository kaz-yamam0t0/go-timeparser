package timeparser

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func detectLocation(zone_name string) (*time.Location, error) {
	loc, err := time.LoadLocation(zone_name)
	if err == nil {
		return loc, nil
	}

	// @TODO detect unknown location
	data := [][]string{
		{"Z", "UTC"},
		{"JST", "Asia/Tokyo"},
	}
	for _, items := range data {
		hit := -1
		for j, item := range items {
			if item == zone_name {
				hit = j
				break
			}
		}
		if hit >= 0 {
			for j, item := range items {
				if j != hit {
					loc, err = time.LoadLocation(item)
					if err == nil {
						return loc, nil
					}
				}
			}
			break
		}
	}
	if err == nil {
		err = errors.New("failed to load Location: " + zone_name)
	}
	return nil, err
}

func parseLocation(s *string, pos_s *int) (*time.Location, bool, error) {
	s_len := len(*s)
	start_pos := *pos_s
	pos := *pos_s

	for pos < s_len {
		if (*s)[pos] == '/' || (*s)[pos] == '_' || (*s)[pos] == '+' || (*s)[pos] == '-' || ('a' <= (*s)[pos] && (*s)[pos] <= 'z') || ('A' <= (*s)[pos] && (*s)[pos] <= 'Z') || ('0' <= (*s)[pos] && (*s)[pos] <= '9') {
			pos++
		} else {
			break
		}
	}
	if pos <= start_pos {
		return nil, false, nil
	}
	if (*s)[start_pos] == '/' || (*s)[pos-1] == '/' {
		return nil, false, errors.New(fmt.Sprintf("invalid timezone name : %s", (*s)[start_pos:pos]))
	}
	(*pos_s) = pos

	// https://github.com/golang/go/blob/eef0140137bae3bc059f598843b8777f9223fac8/src/time/zoneinfo_unix.go#L19-L26
	zone_name := (*s)[start_pos:pos]
	//loc, err := time.LoadLocation(zone_name)
	loc, err := detectLocation(zone_name)
	if err != nil {
		return nil, false, errors.New(fmt.Sprintf("invalid timezone name : %s", (*s)[start_pos:pos]))
	}
	return loc, true, nil
}

func parseZoneOffset(s *string, pos_s *int) (int, bool) {
	s_len := len(*s) - *pos_s

	h_ := 0
	m_ := 0
	n_ := 0
	ok := false

	sign_ := 0

	if s_len > 0 {
		if (*s)[*pos_s] == '+' {
			sign_ = 1
		}
		if (*s)[*pos_s] == '-' {
			sign_ = -1
		}
		if sign_ != 0 {
			(*pos_s)++
			s_len--
		}
	}
	if sign_ == 0 {
		sign_ = 1
	}

	// 00:00
	if s_len >= 5 {
		if isNumeric((*s)[*pos_s]) && isNumeric((*s)[(*pos_s)+1]) && (*s)[(*pos_s)+2] == ':' && isNumeric((*s)[(*pos_s)+3]) && isNumeric((*s)[(*pos_s)+4]) {
			h_, _ = parseInt(s, pos_s, 2, 2)
			(*pos_s)++
			m_, _ = parseInt(s, pos_s, 2, 2)

			tz := (h_*60 + m_) * 60 * sign_
			return tz, true
		}
	}
	// 0:00
	// 0000
	if s_len >= 4 {
		if isNumeric((*s)[*pos_s]) && (*s)[(*pos_s)+1] == ':' && isNumeric((*s)[(*pos_s)+2]) && isNumeric((*s)[(*pos_s)+3]) {
			h_, _ = parseInt(s, pos_s, 1, 1)
			(*pos_s)++
			m_, _ = parseInt(s, pos_s, 2, 2)

			tz := (h_*60 + m_) * 60 * sign_
			return tz, true
		} else if n_, ok = parseInt(s, pos_s, 4, 4); ok {
			h_ = (n_ / 100)
			m_ = (n_ % 100)
			if h_ > 60 || m_ > 60 {
				return -1, false
			}
			tz := (h_*60 + m_) * 60 * sign_
			return tz, true
		}
	}

	return -1, false
}
func parseAMPM(s *string, pos_s *int) (int, bool) {
	_s := strings.ToLower((*s)[*pos_s:])
	switch {
	// AM
	case strings.HasPrefix(_s, "a.m."):
		*pos_s += 4
		return AM, true
	case strings.HasPrefix(_s, "am."):
		*pos_s += 3
		return AM, true
	case strings.HasPrefix(_s, "am"):
		*pos_s += 2
		return AM, true

	// PM
	case strings.HasPrefix(_s, "p.m."):
		*pos_s += 4
		return PM, true
	case strings.HasPrefix(_s, "pm."):
		*pos_s += 3
		return PM, true
	case strings.HasPrefix(_s, "pm"):
		*pos_s += 2
		return PM, true
	}
	return -1, false
}
func parseMonth(s *string, pos_s *int) (int, bool) {
	month_num, month_name := startsWithMonthName(strings.ToLower((*s)[*pos_s:]))
	if month_num > 0 {
		*pos_s += len(month_name)
		return month_num, true
	}
	return -1, false
}
func parseWeekday(s *string, pos_s *int) (int, bool) {
	weekday_num, weekday_name := startsWithWeekdayName(strings.ToLower((*s)[*pos_s:]))
	if weekday_num > 0 {
		*pos_s += len(weekday_name)
		return weekday_num, true
	}
	return -1, false
}
func parseInt(s *string, pos_s *int, min_length int, max_length int) (int, bool) {
	s_len := len(*s)
	max_ := max_length
	if s_len-*pos_s <= max_ {
		max_ = s_len - *pos_s
	}
	if max_ < min_length {
		return -1, false
	}

	res := 0
	for i := 0; i < max_; i++ {
		if (*s)[*pos_s] < '0' || '9' < (*s)[*pos_s] {
			if i < min_length {
				return -1, false
			}
			break
		}
		res = res*10 + int((*s)[*pos_s]-'0')
		(*pos_s)++
	}
	return res, true
}
func parseDecimal(s *string, pos_s *int, min_length int, max_length int) (float64, bool) {
	s_len := len(*s)
	max_ := max_length
	if s_len-*pos_s < max_ {
		max_ = s_len - *pos_s
	}
	if max_ < min_length {
		return -1, false
	}

	res := float64(0)
	div := float64(10)
	for i := 0; i < max_; i++ {
		if (*s)[*pos_s] < '0' || '9' < (*s)[*pos_s] {
			if i < min_length {
				return -1, false
			}
			break
		}
		res += float64((*s)[*pos_s]-'0') / div
		div = div * 10
		(*pos_s)++
	}
	return res, true
}

func parseSuffix(s *string, pos_s *int) bool {
	if len(*s) <= *pos_s+1 {
		return false
	}
	s_ := strings.ToLower((*s)[*pos_s : *pos_s+2])
	if s_ == "st" || s_ == "nd" || s_ == "rd" || s_ == "th" {
		*pos_s += 2
		return true
	}
	return false
}

func parseFormatChar(format *string, pos *int, s *string, pos_s *int, d *TimeData) (int, error) {
	if (*format)[*pos] == '\\' {
		(*pos)++
		if (*format)[*pos] == (*s)[*pos_s] {
			(*pos)++
			(*pos_s)++
		}
	}
	not_hit := false
	n := 0
	ok := false
	switch (*format)[*pos] {
	// Long Format
	case 'r':
		_format := "D, d M Y H:i:s O"
		_pos := 0
		for _pos < len(_format) {
			_ = skipSpaces(&_format, &_pos)
			_ = skipSpaces(s, pos_s)

			if _, err := parseFormatChar(&_format, &_pos, s, pos_s, d); err != nil {
				return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
			}
		}
		(*pos)++
	case 'c':
		_format := "Y-m-d\\TH:i:sP"
		_pos := 0
		for _pos < len(_format) {
			_ = skipSpaces(&_format, &_pos)
			_ = skipSpaces(s, pos_s)

			if _, err := parseFormatChar(&_format, &_pos, s, pos_s, d); err != nil {
				return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
			}
		}
		(*pos)++

		/*
						return FormatTime("Y-m-d\\TH:i:sP", d), true
			case 'r':
				return FormatTime("D, d M Y H:i:s O", d), true
		*/
	// Day
	case 'd':
		fallthrough
	case 'j':
		if n, ok = parseInt(s, pos_s, 1, 2); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		if n < 0 || 31 < n {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		d.setDay(n)
		(*pos)++
	// Weekday
	case 'D':
		fallthrough
	case 'l':
		if n, ok = parseWeekday(s, pos_s); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		d.setDay(n)
		(*pos)++
	// suffix (ignores)
	case 'S':
		if parseSuffix(s, pos_s) != true {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		(*pos)++
	// Day of Year
	case 'z':
		if n, ok = parseInt(s, pos_s, 1, 3); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		if n <= 0 || 365 < n {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		d.setMonth(1)
		d.setDay(1 + n - 1)
		d.normalizeYmd()
		(*pos)++
	// Months
	case 'm':
		fallthrough
	case 'n':
		if n, ok = parseInt(s, pos_s, 1, 2); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		if n < 0 || 12 < n {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		d.setMonth(n)
		(*pos)++
	case 'F':
		fallthrough
	case 'M':
		if n, ok = parseMonth(s, pos_s); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		d.setMonth(n)
		(*pos)++
	// Year
	case 'Y':
		if n, ok = parseInt(s, pos_s, 4, 4); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s %s %d", string((*format)[*pos]), (*s), *pos_s))
			//return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		d.setYear(n)
		(*pos)++
	case 'y':
		if n, ok = parseInt(s, pos_s, 1, 2); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		if n < 70 {
			d.setYear(n + 2000)
		} else {
			d.setYear(n + 1900)
		}
		(*pos)++
	// AM/PM
	case 'a':
		fallthrough
	case 'A':
		if n, ok = parseAMPM(s, pos_s); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		d.ap = n
		if d.ap == PM && 1 <= d.h && d.h < 12 {
			d.h += 12
			d.ap = 0
		}
		(*pos)++
	// Hour
	case 'g':
		fallthrough
	case 'h':
		if n, ok = parseInt(s, pos_s, 1, 2); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		if n <= 0 || 12 < n {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		if d.ap == PM {
			n += 12
			d.ap = 0
		}
		d.setHour(n)
		(*pos)++
	case 'G':
		fallthrough
	case 'H':
		if n, ok = parseInt(s, pos_s, 1, 2); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		if n <= 0 || 23 < n {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		d.setHour(n)
		(*pos)++
	// Minute
	case 'i':
		if n, ok = parseInt(s, pos_s, 1, 2); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		if n <= 0 || 59 < n {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		d.setMinute(n)
		(*pos)++
	// Seconds
	case 's':
		if n, ok = parseInt(s, pos_s, 1, 2); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		if n <= 0 || 59 < n {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		d.setSecond(n)
		(*pos)++
	// Milliseconds
	case 'v':
		var f float64
		if f, ok = parseDecimal(s, pos_s, 1, 6); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		d.setNanosecond(int(f * 1e9))
		(*pos)++
	// Microseconds
	case 'u':
		var f float64
		if f, ok = parseDecimal(s, pos_s, 1, 6); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		d.setNanosecond(int(f * 1e9))
		(*pos)++
	// Unixtime
	case 'U':
		if n, ok = parseInt(s, pos_s, 0, 20); !ok {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		// set Unixtime
		_t := time.Unix(int64(n), 0) // .In()
		d.setYear(_t.Year())
		d.setMonth(int(_t.Month()))
		d.setDay(_t.Day())
		d.setHour(_t.Hour())
		d.setMinute(_t.Minute())
		d.setSecond(_t.Second())
		d.setNanosecond(0)
		(*pos)++
	// Timezone
	case 'e':
		fallthrough
	case 'O':
		fallthrough
	case 'P':
		fallthrough
	case 'T':
		// [+-]00:00 format means timezone offset from UTC
		if n, ok = parseZoneOffset(s, pos_s); ok {
			d.setTimezoneOffset(n)

			loc, err_ := time.LoadLocation("UTC")
			if err_ != nil {
				return -1, err_
			}
			d.setLocation(loc)
		} else {
			loc, ok_, err_ := parseLocation(s, pos_s)
			if ok_ != true {
				if err_ != nil {
					return -1, err_
				} else {
					return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
				}
			}
			d.setLocation(loc)
		}

		(*pos)++

	// Other formats
	case '!':
		d.reset()
	case '|':
		d.resetIfUnset()
		d.ap = 0
		(*pos)++
	case '+':
		d.flags |= SKIP_ERRORS
		(*pos)++
	case '#':
		if !isSeparator((*s)[*pos_s]) {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
		(*pos_s)++
		(*pos)++
	case '?':
		(*pos_s)++
		(*pos)++
	case '*':
		next_c := -1
		if *pos+1 < len(*format) {
			next_c = int((*format)[(*pos)+1])
		}

		s_len := len(*s)
		(*pos_s)++
		for *pos_s < s_len && int((*s)[*pos_s]) != next_c && !isAlphanumeric((*s)[*pos_s]) {
			(*pos_s)++
		}
		(*pos)++

	// other same letter
	//case (*s)[*pos_s]:
	//	(*pos_s)++
	//	(*pos)++

	// Not hit
	default:
		not_hit = true
	}
	if not_hit {
		// space
		if isSpace((*format)[*pos]) {
			_ = skipSpaces(format, pos)
			_ = skipSpaces(s, pos_s)
		} else if (*format)[*pos] == (*s)[*pos_s] {
			(*pos_s)++
			(*pos)++
		} else {
			return -1, errors.New(fmt.Sprintf("failed to parse format: %s", string((*format)[*pos])))
		}
	}

	return (*pos_s), nil
}

// Convert a datetime string to a time.Time variable with format specification
func ParseFormat(format string, s string) (*time.Time, error) {
	format = strings.TrimSpace(format)
	if format == "" {
		return nil, errors.New("empty format")
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, errors.New("empty data")
	}

	pos_s := 0
	pos := 0

	s_len := len(s)
	f_len := len(format)

	tmp := 0

	data := newTimeData()
	for pos < f_len {
		// skip spaces
		if isSpace(format[pos]) {
			_ = skipSpaces(&format, &pos)
			_ = skipSpaces(&s, &pos_s)
		}

		if pos >= f_len {
			break
		}
		if pos_s >= s_len {
			if data.hasFlag(SKIP_ERRORS) {
				break
			}
			return nil, errors.New(fmt.Sprintf("failed to parse format: %s", string(format[pos])))
		}

		if _, err := parseFormatChar(&format, &pos, &s, &pos_s, data); err != nil {
			if data.hasFlag(SKIP_ERRORS) {
				break
			}
			return nil, err
		}

		tmp++
		if tmp > 30 {
			break
		}
	}

	res := data.Time()

	return res, nil
}
