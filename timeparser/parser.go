package timeparser

import (
	"errors"
	"strings"
	"time"
	//"fmt"
)

func preprocessScannedStr(s string) string {
	s_len := len(s)
	dst := make([]byte, 0, s_len)

	head_flg := true // head of a word
	pos := 0
	for pos = 0; pos < s_len; pos++ {
		// normalize spaces
		if isSpace(s[pos]) || s[pos] == ',' {
			pos++
			for pos < s_len && (isSpace(s[pos]) || s[pos] == ',') {
				pos++
			}
			if pos >= s_len {
				break
			}
			dst = append(dst, ' ')
			head_flg = true
		}

		if head_flg {
			// ignore "the"
			if s_len-pos > 3 && strings.ToLower(s[pos:pos+3]) == "the" {
				if s_len-pos < 4 || isSpace(s[pos+4]) {
					pos += 4
				}
			}
			if pos >= s_len {
				break
			}

			// iso format
		}
		// lower
		c := s[pos]
		//if 'A' <= c && c <= 'Z' {
		//	c = c - 'A' + 'a'
		//}

		dst = append(dst, c)

		// head_flg for the next character
		head_flg = isSpace(s[pos])
	}
	return string(dst)
}
func scanWord(s string, pos_s int, word string, check_end bool) int {
	s_len := len(s)
	length := len(word)
	if pos_s+length > s_len {
		return -1
	}
	for i := 0; i < length; i++ {
		if !cmpichr(s[pos_s+i], word[i]) {
			return -1
		}
	}
	if check_end {
		if pos_s+length < s_len {
			c := s[pos_s+length]
			if !isSpace(c) && !isSeparator(c) {
				return -1
			}
		}
	}
	return length
}
func scanWords(s string, pos_s int, words []string, check_end bool) *string {
	for _, w := range words {
		if scanWord(s, pos_s, w, check_end) > 0 {
			return &w
		}
	}
	return nil
}

func scanSuffix(s string, pos_s int) (length int) {
	if scanWord(s, pos_s, "st", true) > 0 || scanWord(s, pos_s, "nd", true) > 0 || scanWord(s, pos_s, "rd", true) > 0 || scanWord(s, pos_s, "th", true) > 0 {
		return 2
	}
	return -1
}
func scanDayWithSuffix(s string, pos_s int) (d int, length int) {
	//s_len := len(s)
	pos := pos_s

	d_ := 0
	ok := false
	if d_, ok = parseInt(&s, &pos, 1, 2); !ok {
		return -1, -1
	}

	len_ := -1
	if len_ = scanSuffix(s, pos); len_ <= 0 {
		return -1, -1
	}
	pos += len_

	if d_ < 1 || 31 < d_ {
		return -1, -1
	}
	return d_, (pos - pos_s)
}

func scanMonth(s string, pos_s int) (m int, length int) {
	m = -1
	length = -1

	s_len := len(s)
	for month_name, month_num := range getMonthTable() {
		hit := true

		l := len(month_name)
		pos := pos_s
		for i := 0; i < l; i++ {
			if pos >= s_len {
				hit = false
				break
			}
			if !cmpichr(s[pos], month_name[i]) {
				hit = false
				break
			}
			pos++
			if i == 2 {
				if pos >= s_len || isSpace(s[pos]) {
					hit = true
					break
				}
			}
		}
		if hit {
			if pos < s_len && isAlphanumeric(s[pos]) {
				return -1, -1
			}
			return month_num, (pos - pos_s)
		}
	}
	return -1, -1
}
func scanWeekday(s string, pos_s int) (w int, length int) {
	w = -1
	length = -1

	s_len := len(s)
	for weekday_name, weekday_num := range getWeekdayTable() {
		hit := true

		l := len(weekday_name)
		pos := pos_s
		for i := 0; i < l; i++ {
			if pos >= s_len {
				hit = false
				break
			}
			if !cmpichr(s[pos], weekday_name[i]) {
				hit = false
				break
			}
			pos++
			if i == 2 {
				if pos >= s_len || isSpace(s[pos]) {
					hit = true
					break
				}
			}
		}
		if hit {
			if pos < s_len && isAlphanumeric(s[pos]) {
				return -1, -1
			}
			return weekday_num, (pos - pos_s)
		}
	}
	return -1, -1
}

func scanTimezoneOffset(s string, pos_s int) (s_ int, length int) {
	// ([\-\+]?)(\d{2})\:?(\d{2})
	s_len := len(s)
	pos := pos_s

	sign_ := 1
	ok := false

	h_ := -1
	m_ := -1

	if s_len-pos < 4 {
		return -1, -1
	}
	if s[pos] != 'Z' && s[pos] != 'z' && s[pos] != '+' && s[pos] != '-' {
		return -1, -1
	}

	// Z
	if s[pos] == 'Z' || s[pos] == 'z' {
		pos++
	}

	// + or -
	if s[pos] == '+' {
		pos++
	} else if s[pos] == '-' {
		pos++
		sign_ = -1
	}
	if s_len-pos < 4 {
		return -1, -1
	}
	// h
	if h_, ok = parseInt(&s, &pos, 1, 2); !ok {
		return -1, -1
	}
	// :
	if pos < s_len && s[pos] == ':' {
		pos++
	}

	// m
	if m_, ok = parseInt(&s, &pos, 2, 2); !ok {
		return -1, -1
	}
	// next character must not be a digit
	if pos < s_len && isNumeric(s[pos]) {
		return -1, -1
	}
	return ((h_*60 + m_) * 60 * sign_), (pos - pos_s)
}

func scanLocation(s string, pos_s int) (*time.Location, int) {
	// ([\-\+]?)(\d{2})\:?(\d{2})
	s_len := len(s)
	pos := pos_s

	for pos < s_len {
		if s[pos] == '/' || s[pos] == '_' || s[pos] == '+' || s[pos] == '-' || ('a' <= s[pos] && s[pos] <= 'z') || ('A' <= s[pos] && s[pos] <= 'Z') || ('0' <= s[pos] && s[pos] <= '9') {
			pos++
		} else {
			break
		}
	}
	if pos <= pos_s {
		return nil, -1
	}
	if s[pos_s] == '/' || s[pos-1] == '/' {
		return nil, -1
	}

	zone_name := s[pos_s:pos]
	loc, err := detectLocation(zone_name)
	if err != nil {
		return nil, -1
	}

	return loc, (pos - pos_s)
}

func scanTime(s string, pos_s int) (h_ int, m_ int, s_ int, ns_ int, length int) {
	// (\d{2})\:(\d{2})(\:(\d{2}))?( (a\.m\.|p\.m\.|am|pm))?
	// (\d{2})\:(\d{2})(\:(\d{2}))?(\.\d+)?( (a\.m\.|p\.m\.|am|pm))?
	// may start with "T"
	h_ = 0
	m_ = 0
	s_ = 0
	ns_ = 0
	length = -1

	ap_ := -1

	s_len := len(s)
	pos := pos_s
	ok := false

	// T
	if pos < s_len && (s[pos] == 't' || s[pos] == 'T') {
		pos++
	}

	// h
	if h_, ok = parseInt(&s, &pos, 1, 2); !ok {
		return -1, -1, -1, -1, -1
	}
	// :
	if pos >= s_len || s[pos] != ':' {
		return -1, -1, -1, -1, -1
	}
	pos++
	// m
	if m_, ok = parseInt(&s, &pos, 2, 2); !ok {
		return -1, -1, -1, -1, -1
	}
	// (:s)?
	if pos < s_len && s[pos] == ':' {
		pos++
		tmp_ := -1
		if tmp_, ok = parseInt(&s, &pos, 2, 2); ok {
			s_ = tmp_
		}

		// (.\d+)?
		if pos < s_len && s[pos] == '.' {
			pos++
			tmp2_ := float64(-1)
			if tmp2_, ok = parseDecimal(&s, &pos, 1, 11); ok {
				ns_ = int(tmp2_ * 1e9)
			} else {
				return -1, -1, -1, -1, -1
			}
		}
	}
	// next character must not be a digit
	if pos < s_len && isNumeric(s[pos]) {
		return -1, -1, -1, -1, -1
	}

	// spaces
	skipSpaces(&s, &pos)

	if (h_ < 0 || 23 < h_) || (m_ < 0 || 59 < m_) || (59 < s_) {
		return -1, -1, -1, -1, -1
	}

	// am / pm
	if ap_, ok = parseAMPM(&s, &pos); ok {
		if h_ < 1 && 12 < h_ {
			return -1, -1, -1, -1, -1
		}
		if ap_ == AM {
			if h_ == 12 {
				h_ = 0
			}
		} else if ap_ == PM {
			if h_ == 12 {
				h_ = 12
			} else {
				h_ += 12
			}
		}
	}
	return h_, m_, s_, ns_, (pos - pos_s)
}

func scanDmy(s string, pos_s int) (y int, m int, d int, length int) {
	// (\d{1,2})(st|nd|rd|th)?\s*` + _months + `(\s+(\d+))?
	// 1st january 2006
	y = -1
	m = -1
	d = -1
	length = -1

	pos := pos_s
	ok := false
	len_ := -1

	// day
	if d, ok = parseInt(&s, &pos, 1, 2); !ok {
		return -1, -1, -1, -1
	}
	// suffix
	if len_ = scanSuffix(s, pos); len_ > 0 {
		pos += len_
	}
	_ = skipSpaces(&s, &pos)

	// month
	if m, len_ = scanMonth(s, pos); len_ <= 0 {
		return -1, -1, -1, -1
	}
	pos += len_

	_ = skipSpaces(&s, &pos)

	// year
	if y, ok = parseInt(&s, &pos, 4, 4); !ok {
		return -1, -1, -1, -1
	}
	return y, m, d, (pos - pos_s)
}

func scanYmd(s string, pos_s int) (y int, m int, d int, length int) {
	// (\d+)([/\-\.])(\d+)([/\-\.])(\d+)
	y = -1
	m = -1
	d = -1
	length = -1

	n1 := 0
	n2 := 0
	n3 := 0
	var sep byte

	s_len := len(s)
	pos := pos_s
	ok := false

	// n1
	if n1, ok = parseInt(&s, &pos, 1, 4); !ok {
		return -1, -1, -1, -1
	}
	// sep
	if pos >= s_len || (s[pos] != '-' && s[pos] != '/' && s[pos] != '.') {
		return -1, -1, -1, -1
	}
	sep = s[pos]
	pos++

	// n2
	if n2, ok = parseInt(&s, &pos, 1, 4); !ok {
		return -1, -1, -1, -1
	}

	// sep
	if pos >= s_len || s[pos] != sep {
		return -1, -1, -1, -1
	}
	pos++

	// n3
	if n3, ok = parseInt(&s, &pos, 1, 4); !ok {
		return -1, -1, -1, -1
	}
	// next character must not be a digit
	if pos < s_len && isNumeric(s[pos]) {
		return -1, -1, -1, -1
	}

	switch sep {
	case '.': //  d.m.y
		if !checkDate(n3, n2, n1) {
			return -1, -1, -1, -1
		}
		y, m, d = n3, n2, n1
	case '-': //  y-m-d d-m-y
		if checkDate(n3, n2, n1) {
			y, m, d = n3, n2, n1
		} else if checkDate(n1, n2, n3) {
			y, m, d = n1, n2, n3
		} else {
			return -1, -1, -1, -1
		}
	case '/': // y/m/d m/d/y
		if checkDate(n3, n1, n2) {
			y, m, d = n3, n1, n2
		} else if checkDate(n1, n2, n3) {
			y, m, d = n1, n2, n3
		} else {
			return -1, -1, -1, -1
		}
	}
	return y, m, d, (pos - pos_s)
}
func scanRelativePosition(s string, pos_s int) (n int, unit string, length int) {
	// ([\+\-]?)\s*(\d+|a)\s*(year|month|day|hour|minute|second|week|millisecond|microsecond|msec|ms|µsec|µs|usec|sec|min|forth?night)s?(\s+ago)?\b
	n = -1
	unit = ""
	length = -1

	s_len := len(s)
	pos := pos_s
	ok := false

	// sign
	sign_ := 1
	if s[pos] == '+' {
		pos++
	} else if s[pos] == '-' {
		sign_ = -1
		pos++
	}
	_ = skipSpaces(&s, &pos)

	// \d+ | a
	if s[pos] == 'a' && pos+1 < s_len && isSpace(s[pos+1]) {
		n = 1
		pos++
	} else if n, ok = parseInt(&s, &pos, 1, 29); !ok {
		return -1, "", -1
	}
	n = n * sign_

	_ = skipSpaces(&s, &pos)

	// unit
	units := []string{
		"year", "month", "day", "hour", "minute", "second",
		"week", "millisecond", "microsecond", "msec", "ms",
		"µsec", "µs", "usec", "nanosecond", "nsec", "ns",
		"sec", "min", "fortnight", "forthnight",
	}
	s_ := scanWords(s, pos, units, false)
	if s_ == nil {
		return -1, "", -1
	}
	unit = *s_
	pos += len(unit)

	if pos < s_len && (s[pos] == 's' || s[pos] == 'S') {
		pos++
	}

	if pos < s_len && isSpace(s[pos]) {
		_ = skipSpaces(&s, &pos)

		// ago
		if len_ := scanWord(s, pos, "ago", true); len_ > 0 {
			pos += len_
			n = n * -1
		}
	}

	return n, unit, (pos - pos_s)
}
func scanPosition(s string, pos_s int) (*timeAddition, int) {
	// `^((first|last) day of )?(next|last) ((year|month|day)|` + _months + `|` + _weeks + `)

	pos := pos_s
	len_ := 0

	var day_flg_ *string = nil
	var pos_flg_ *string = nil
	var word_ *string = nil
	var m_ int = -1
	var w_ int = -1

	// first | last day of
	day_flg_ = scanWords(s, pos, []string{"first", "last"}, true)
	if day_flg_ != nil {
		pos += len(*day_flg_)

		_ = skipSpaces(&s, &pos)

		// day of
		len_ = scanWord(s, pos, "day", true)
		if len_ < 0 {
			if *day_flg_ == "last" {
				// reset scanning "last day of ..." (it may continue to "last year|day|...")
				pos = pos_s
				day_flg_ = nil
				goto day_pos
			}
			return nil, -1
		}
		pos += len_
		_ = skipSpaces(&s, &pos)

		len_ = scanWord(s, pos, "of", true)
		if len_ < 0 {
			return nil, -1
		}
		pos += len_
		_ = skipSpaces(&s, &pos)
	}
day_pos:
	// next | last
	pos_flg_ = scanWords(s, pos, []string{"next", "last"}, true)
	if pos_flg_ == nil {
		return nil, -1
	}
	pos += len(*pos_flg_)
	_ = skipSpaces(&s, &pos)

	// year | month | day | $month_names | $weekday_names
	if word_ = scanWords(s, pos, []string{"year", "month", "day"}, true); word_ != nil {
		pos += len(*word_)
	} else if m_, len_ = scanMonth(s, pos); len_ >= 0 {
		pos += len_
	} else if w_, len_ = scanWeekday(s, pos); len_ >= 0 {
		pos += len_
	} else {
		return nil, -1
	}

	// new time addition
	a := newTimeAddition(0, "")
	if day_flg_ != nil {
		a.day_flg = *day_flg_
	}
	if pos_flg_ != nil {
		a.pos = *pos_flg_
	}

	a.month = m_
	a.weekday = w_
	if word_ != nil {
		a.word = *word_
	}
	return a, (pos - pos_s)
}
func scanISOInterval(s string, pos_s int) ([]*timeAddition, int) {
	// P1Y2M3DT4H5M6.7S
	a := make([]*timeAddition, 0, 7)

	s_len := len(s)
	pos := pos_s

	n := -1
	ok := false

	if pos >= s_len || (s[pos] != 'P' && s[pos] != 'p') {
		return nil, -1
	}
	pos++

	for pos < s_len && ('0' <= s[pos] && s[pos] <= '9') {
		pos_ := pos
		if n, ok = parseInt(&s, &pos, 1, 29); !ok {
			break
		}
		if pos >= s_len {
			pos = pos_
			break
		}
		c := s[pos]
		pos++

		if 'A' <= c && c <= 'Z' {
			c = c - 'A' + 'a'
		}

		switch {
		case c == 'y' || c == 'Y':
			a = append(a, newTimeAddition(n, "year"))
		case c == 'm' || c == 'M':
			a = append(a, newTimeAddition(n, "month"))
		case c == 'd' || c == 'D':
			a = append(a, newTimeAddition(n, "day"))
		default:
			break
		}
	}
	// return if not T
	if pos >= s_len || (s[pos] != 'T' && s[pos] != 't') {
		return a, (pos - pos_s)
	}
	pos++

	for pos < s_len && ('0' <= s[pos] && s[pos] <= '9') {
		pos_ := pos
		if n, ok = parseInt(&s, &pos, 1, 29); !ok {
			break
		}
		if pos >= s_len {
			pos = pos_
			break
		}
		c := s[pos]
		pos++

		// decimal
		var d float64 = 0.0
		if c == '.' {
			if d, ok = parseDecimal(&s, &pos, 1, 9); !ok {
				return nil, -1
			}
			c = s[pos]
			pos++
		}

		switch {
		case c == 'h' || c == 'H':
			a = append(a, newTimeAddition(n, "hour"))
		case c == 'm' || c == 'M':
			a = append(a, newTimeAddition(n, "minute"))
		case c == 's' || c == 'S':
			a = append(a, newTimeAddition(n, "second"))
			if d > 0 {
				a = append(a, newTimeAddition(int(d*1e6), "microsecond"))
			}
		default:
			break
		}
	}

	return a, (pos - pos_s)
}

// scan format chunk and add pos
func scanFormat(data *TimeData, s string, pos int) int {
	s_len := len(s)
	if pos >= s_len {
		return -1
	}
	// \s\n\r\t
	for s_len > pos && (s[pos] == 0x20 || s[pos] == '\n' || s[pos] == '\r' || s[pos] == '\t') {
		pos++
	}
	if pos >= s_len {
		return pos
	}

	_s := s[pos:]
	_s_len := len(_s)
	_ = _s_len

	// keywords
	if len_ := scanWord(s, pos, "now", true); len_ > 0 {
		data.setNow()
		pos += len_
	} else if len_ := scanWord(s, pos, "yesterday", true); len_ > 0 {
		data.appendAddition(newTimeAdditionWithTime(-1, "day", 0, 0, 0, 0))
		pos += len_
	} else if len_ := scanWord(s, pos, "midnight", true); len_ > 0 {
		data.appendAddition(newTimeAdditionWithTime(0, "day", 0, 0, 0, 0))
		pos += len_
	} else if len_ := scanWord(s, pos, "today", true); len_ > 0 {
		data.appendAddition(newTimeAdditionWithTime(0, "day", 0, 0, 0, 0))
		pos += len_
	} else if len_ := scanWord(s, pos, "noon", true); len_ > 0 {
		data.appendAddition(newTimeAdditionWithTime(0, "day", 0, 0, 0, 0))
		pos += len_
	} else if m_, len_ := scanMonth(s, pos); len_ >= 0 {
		// month name
		data.setMonth(m_)
		pos += len_
	} else if _, len_ := scanWeekday(s, pos); len_ >= 0 {
		// weekday name
		//data.setWeekday(num_)
		pos += len_
	} else if n_, unit_, len_ := scanRelativePosition(s, pos); len_ >= 0 {
		// relative format (1 year .. etc)
		data.appendAddition(newTimeAddition(n_, unit_))
		pos += len_
	} else if a_, len_ := scanISOInterval(s, pos); len_ >= 0 {
		// iso8601 interval format (P1Y2M3DT1H2M3S)
		for _, _a := range a_ {
			data.appendAddition(_a)
		}
		pos += len_
	} else if _a, len_ := scanPosition(s, pos); len_ >= 0 {
		// last|next day of the last month
		data.appendAddition(_a)
		pos += len_
	} else if y_, m_, d_, len_ := scanDmy(s, pos); len_ > 0 {
		// 10th January 2021
		data.setYear(y_)
		data.setMonth(m_)
		data.setDay(d_)
		pos += len_
	} else if d_, len_ := scanDayWithSuffix(s, pos); len_ > 0 {
		// 10th
		data.setDay(d_)
		pos += len_
	} else if y_, m_, d_, len_ := scanYmd(s, pos); len_ > 0 {
		// Y-m-d d.m.Y etc
		data.setYear(y_)
		data.setMonth(m_)
		data.setDay(d_)
		pos += len_
	} else if h_, m_, s_, ns_, len_ := scanTime(s, pos); len_ > 0 {
		// 00:00(:00)? (am|pm)?
		data.setHour(h_)
		data.setMinute(m_)
		data.setSecond(s_)
		data.setNanosecond(ns_)
		pos += len_
	} else if s_, len_ := scanTimezoneOffset(s, pos); len_ > 0 {
		// +00:00
		// Z00:00
		loc, err_ := time.LoadLocation("UTC")
		if err_ != nil {
			return -1
		}
		data.setLocation(loc)
		data.setTimezoneOffset(s_)
		pos += len_
	} else if loc_, len_ := scanLocation(s, pos); len_ > 0 {
		data.setLocation(loc_)
		data.setTimezoneOffset(0)
		pos += len_
	} else if y_, ok := parseInt(&s, &pos, 4, 4); ok { // pos has been already increased
		// Year
		data.setYear(y_)
	} else {
		return -1
	}

	return pos
}

// Convert string to a time.Time variable
func parseTimeStr(format string, base *time.Time) (*TimeData, error) {
	//s := strings.TrimSpace(strings.ToLower(format))
	s := strings.TrimSpace(format)
	if s == "" {
		return nil, errors.New("Failed to parse Time")
	}

	// data
	data := newTimeData()

	if base == nil {
		t_ := time.Now()
		base = &t_
	}
	data.setFromTime(base)

	// convert datetime
	s = preprocessScannedStr(s)

	// parse string
	s_len := len(s)
	pos := 0
	cnts := 0
	for cnts < 15 && pos < s_len {
		//pos_s := pos
		pos = scanFormat(data, s, pos)
		if pos < 0 {
			//fmt.Println(pos_s, s)
			return nil, errors.New("Failed to parse Time")
		}
		cnts++
	}

	// additions
	data.processAdditions()

	return data, nil
}

// Convert string to a time.Time variable
func ParseTimeStr(format string, base *time.Time) (*time.Time, error) {
	data, err := parseTimeStr(format, base)
	if err != nil {
		return nil, err
	}

	// result
	res := data.Time()

	// return &res, &data
	return res, nil
}
