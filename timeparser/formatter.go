package timeparser

import (
	"strconv"
	"time"
)

// used in date
func paddingZero(n string, length int) string {
	if len(n) > length {
		return n
	}
	dst := make([]byte, length)
	n_len := len(n)

	for i := 0; i < length; i++ {
		if length-n_len <= i {
			dst[i] = n[i-(length-n_len)]
		} else {
			dst[i] = '0'
		}
	}

	return string(dst)
}

func tzFormat(d *time.Time, o_flg bool) string {
	dst := ""

	_, offset_ := d.Zone()
	if offset_ >= 0 {
		dst += "+"
	} else {
		dst += "-"
		offset_ *= -1
	}

	hour_ := (offset_ / 3600) % 60
	min_ := (offset_ / 60) % 60

	dst += paddingZero(strconv.Itoa(hour_), 2)
	if o_flg == true {
		dst += ":"
	}
	dst += paddingZero(strconv.Itoa(min_), 2)

	return dst
}
func iso8601FirstDate(d *time.Time) *time.Time {
	_d := time.Date(d.Year(), 1, 1, 0, 0, 0, 0, d.Location())
	_t := _d.Unix()
	_w := int64(_d.Weekday())
	if _w > 3 {
		_t += (7 - _w) * 86400
	} else {
		_t -= (7 - _w) * 86400
	}
	t := time.Unix(_t, 0)
	return &t
}

func timeFormatChr(f byte, d *time.Time) (string, bool) {
	switch f {
	// Date
	case 'd':
		return strconv.Itoa(d.Day()), true
	case 'j':
		return paddingZero(strconv.Itoa(d.Day()), 2), true
	case 'S':
		d_ := d.Day()
		switch {
		case d_%10 == 1:
			return "st", true
		case d_%10 == 2:
			return "nd", true
		case d_%10 == 3:
			return "rd", true
		default:
			return "th", true
		}
	// Day
	case 'D':
		return d.Weekday().String()[0:3], true
	case 'l':
		return d.Weekday().String(), true
	case 'w':
		return strconv.Itoa(int(d.Weekday())), true
	case 'W':
		_d := iso8601FirstDate(d)
		return strconv.Itoa(int((d.Unix()-_d.Unix())/86400/7) + 1), true
	case 'N':
		return strconv.Itoa(int((d.Weekday()+6)%7 + 1)), true

	// Month
	case 'F':
		return d.Month().String(), true
	case 'm':
		return paddingZero(strconv.Itoa(int(d.Month())), 2), true
	case 'M':
		return d.Month().String()[0:3], true
	case 'n':
		return strconv.Itoa(int(d.Month())), true
	case 't':
		y_ := d.Year()
		m_ := d.Month()
		if m_ > 12 {
			m_ = 1
			y_++
		} else {
			m_++
		}
		d_ := time.Date(y_, m_, 1, 0, 0, 0, 0, d.Location()).AddDate(0, 0, -1)
		return strconv.Itoa(d_.Day()), true

	// Year
	case 'Y':
		return strconv.Itoa(int(d.Year())), true
	case 'y':
		return strconv.Itoa(int(d.Year()))[2:4], true
	case 'L':
		y_ := int(d.Year())
		if y_%400 == 0 || (y_%4 == 0 && y_%100 != 0) {
			return "1", true
		}
		return "0", true
	case 'o':
		_d := iso8601FirstDate(d)
		_tmp := _d.AddDate(0, 0, 365+7)
		_next_d := iso8601FirstDate(&_tmp)
		_y := d.Year()
		if d.Unix() < _d.Unix() {
			_y--
		} else if d.Unix() > _next_d.Unix() {
			_y++
		}
		return strconv.Itoa(int(_y)), true

	// Time
	case 'a':
		if d.Hour() < 12 {
			return "am", true
		} else {
			return "pm", true
		}
	case 'A':
		if d.Hour() < 12 {
			return "AM", true
		} else {
			return "PM", true
		}
	case 'B':
		//Swatch Internet Time
		h_ := d.Hour()
		m_ := d.Minute()
		s_ := d.Second()
		b_ := strconv.Itoa((int(1000*(3600*h_+60*m_+s_-28800)/86400) + 1000) % 1000)

		return b_, true

	// hours
	case 'g':
		return strconv.Itoa(int(d.Hour()) % 12), true
	case 'G':
		return strconv.Itoa(int(d.Hour())), true
	case 'h':
		return paddingZero(strconv.Itoa(int(d.Hour())%12), 2), true
	case 'H':
		return paddingZero(strconv.Itoa(int(d.Hour())), 2), true

	// minute / seconds
	case 'i':
		return paddingZero(strconv.Itoa(int(d.Minute())), 2), true
	case 's':
		return paddingZero(strconv.Itoa(int(d.Second())), 2), true
	case 'v':
		return paddingZero(strconv.Itoa(int(d.Nanosecond())), 9)[0:3], true
	case 'u':
		return paddingZero(strconv.Itoa(int(d.Nanosecond())), 9)[0:6], true

	// Full Date/Time
	case 'c':
		return FormatTime("Y-m-d\\TH:i:sP", d), true
	case 'r':
		return FormatTime("D, d M Y H:i:s O", d), true
	case 'U':
		return strconv.FormatInt(d.Unix(), 10), true

	// timezone
	case 'e':
		name_, _ := d.Zone()
		return name_, true
	case 'I':
		return "", true // (Not supported) @TODO
	case 'Z':
		_, offset_ := d.Zone()
		return strconv.Itoa(offset_), true
	case 'T':
		name_, _ := d.Zone()
		return name_, true
	case 'P':
		return tzFormat(d, true), true
	case 'O':
		return tzFormat(d, false), true
	}
	return "", false
}

// Format a time.Time variable to a string
func FormatTime(s string, dt *time.Time) string {
	if dt == nil {
		d := time.Now()
		dt = &d
	}

	var dst []byte
	s_len := len(s)
	pos := 0
	for i := 0; i < s_len; i++ {
		// backslash
		if s[i] == '\\' {
			dst = append(dst, s[pos:i]...)
			if i+1 < s_len {
				dst = append(dst, s[i+1])
			}
			pos = i + 2
			i++ // skip
			continue
		}
		if b, ok := timeFormatChr(s[i], dt); ok == true {
			if pos < i {
				dst = append(dst, s[pos:i]...)
			}
			dst = append(dst, b...)
			pos = i + 1
		}
	}
	if pos < s_len {
		dst = append(dst, s[pos:]...)
	}
	return string(dst)
}
