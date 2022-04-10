package timeparser

//import (
//	"strings"
//)

// ==============================================================
// get*Table
// ==============================================================

// get month table
func getMonthTable() map[string]int {
	return map[string]int{
		"january": 1,
		//"jan":       1,
		"february": 2,
		//"feb":       2,
		"march": 3,
		//"mar":       3,
		"april": 4,
		//"apr":       4,
		"may":  5,
		"june": 6,
		//"jun":       6,
		"july": 7,
		//"jul":       7,
		"august": 8,
		//"aug":       8,
		"september": 9,
		//"sep":       9,
		"octobar": 10,
		//"oct":       10,
		"november": 11,
		//"nov":       11,
		"december": 12,
		//"dec":       12,
	}
}

// get weekday table
func getWeekdayTable() map[string]int {
	return map[string]int{
		"sunday":    0,
		"monday":    1,
		"tuesday":   2,
		"wednesday": 3,
		"thursday":  4,
		"friday":    5,
		"saturday":  6,
	}
}

// ==============================================================
// is*
// ==============================================================

func isSpace(c byte) bool {
	return c == 0 || c == 0x20 || c == '\n' || c == '\r' || c == '\t' || c == '\v'
}
func isSeparator(c byte) bool {
	return c == ';' || c == ':' || c == '/' || c == '.' || c == ',' || c == '-' || c == '(' || c == ')'
}
func isNumeric(c byte) bool {
	return '0' <= c && c <= '9'
}
func isAlphanumeric(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func cmpichr(a byte, b byte) bool {
	if 'A' <= a && a <= 'Z' {
		a = a - 'A' + 'a'
	}
	if 'A' <= b && b <= 'Z' {
		b = b - 'A' + 'a'
	}
	return a == b
}
func cmpiStartWith(a string, b string) bool {
	a_len := len(a)
	b_len := len(b)
	if a_len < b_len {
		return false
	}
	for i := 0; i < b_len; i++ {
		if !cmpichr(a[i], b[i]) {
			return false
		}
	}
	return true
}

// ==============================================================
// skip*
// ==============================================================

func skipSpaces(s *string, pos_s *int) int {
	return skipChars(s, pos_s, isSpace)
}
func skipSeparators(s *string, pos_s *int) int {
	return skipChars(s, pos_s, isSeparator)
}

func skipChars(s *string, pos_s *int, callback func(c byte) bool) int {
	start_pos := *pos_s
	s_len := len(*s)
	for *pos_s < s_len {
		if callback((*s)[*pos_s]) {
			(*pos_s)++
		} else {
			break
		}
	}
	return (*pos_s) - start_pos
}

// ==============================================================
// others
// ==============================================================

// check if ymd is correct
func checkDate(y int, m int, d int) bool {
	if y < 0 || (m < 1 || 12 < m) || (d < 1 || 31 < d) {
		return false
	}
	return d <= getLastDay(y, m)
}

// get the last day of a month
func getLastDay(y int, m int) int {
	if m == 4 || m == 6 || m == 9 || m == 11 {
		return 30
	}
	if m == 2 {
		if y%4 == 0 && (y%100 != 0 || y%400 == 0) {
			return 29
		} else {
			return 28
		}
	}
	return 31
}

// month name string to number
func getMonthNum(s string) int {
	// month_name should be the lower cases of a correct month name
	for month_name, month_num := range getMonthTable() {
		if s == month_name {
			return month_num
		}
	}
	return -1
}

// starts with month name
func startsWithMonthName(s string) (int, string) {
	if len(s) < 3 {
		return -1, ""
	}
	// month_name should be the lower cases of a correct month name
	for month_name, month_num := range getMonthTable() {
		if cmpiStartWith(s, month_name) {
			return month_num, month_name
		}
		if s[0:3] == month_name[0:3] {
			return month_num, month_name[0:3]
		}
	}
	return -1, ""
}

// weekday name string to number
func getWeekdayNum(s string) int {
	// day_name should be the lower cases of a correct day name
	for weekday_name, weekday_num := range getWeekdayTable() {
		if s == weekday_name {
			return weekday_num
		}
	}
	return -1
}

// starts with month name
func startsWithWeekdayName(s string) (int, string) {
	// month_name should be the lower cases of a correct month name
	if len(s) < 3 {
		return -1, ""
	}
	for weekday_name, weekday_num := range getWeekdayTable() {
		if cmpiStartWith(s, weekday_name) {
			return weekday_num, weekday_name
		}
		if s[0:3] == weekday_name[0:3] {
			return weekday_num, weekday_name[0:3]
		}
	}
	return -1, ""
}
