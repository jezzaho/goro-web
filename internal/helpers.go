package internal

import (
	"bytes"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func FlattenJSON(data []byte) []byte {
	return bytes.Replace(data, []byte("]["), []byte(","), -1)
}
func SSIMtoDate(s string) string {
	var monthMap = make(map[string]string)
	monthMap["JAN"] = "01"
	monthMap["FEB"] = "02"
	monthMap["MAR"] = "03"
	monthMap["APR"] = "04"
	monthMap["MAY"] = "05"
	monthMap["JUN"] = "06"
	monthMap["JUL"] = "07"
	monthMap["AUG"] = "08"
	monthMap["SEP"] = "09"
	monthMap["OCT"] = "10"
	monthMap["NOV"] = "11"
	monthMap["DEC"] = "12"

	var dateString string

	length := len(s)
	switch length {
	case 6:
		// 4JUL24 012345
		dateString += "20" + s[4:] + "-" + monthMap[s[1:4]] + "-0" + s[:1]
	case 7:
		// 19JUL24 0123456
		dateString += "20" + s[5:] + "-" + monthMap[s[2:5]] + "-" + s[:2]
	}
	return dateString
}

// 845 - 840 (14) = 5
func NumberToTime(n int64) string {
	hours := n / 60
	minutes := n - hours*60
	hoursStr := strconv.FormatInt(hours, 10)
	minutesStr := strconv.FormatInt(minutes, 10)

	if hours < 10 {
		hoursStr = "0" + hoursStr
	}
	if minutes < 10 {
		minutesStr = "0" + minutesStr
	}

	return hoursStr + ":" + minutesStr
}
func DaysOfOperation(s string) string {
	return strings.ReplaceAll(s, " ", ".")
}

func moreThanOneNumberReg(s string) bool {
	digitCount := 0
	for _, char := range s {
		if unicode.IsDigit(char) {
			digitCount++
			if digitCount > 1 {
				return true
			}
		}
	}
	return false
}

func SeparateDays(r []string) [][]string {
	newLines := [][]string{}
	check := string(r[8])
	if moreThanOneNumberReg(check) {
		var days []int
		// Check if contains
		if strings.Contains(check, "1") {
			days = append(days, 1)
		}
		if strings.Contains(check, "2") {
			days = append(days, 2)
		}
		if strings.Contains(check, "3") {
			days = append(days, 3)
		}
		if strings.Contains(check, "4") {
			days = append(days, 4)
		}
		if strings.Contains(check, "5") {
			days = append(days, 5)
		}
		if strings.Contains(check, "6") {
			days = append(days, 6)
		}
		if strings.Contains(check, "7") {
			days = append(days, 7)
		}

		newLines = performSeparation(r, days)
		return newLines
	} else {
		newLines = append(newLines, r)
		return newLines
	}

}

func performSeparation(row []string, d []int) [][]string {
	var v int
	var newRows [][]string
	for _, v = range d {
		from := row[6]
		to := row[7]
		cpy := make([]string, len(row))
		copy(cpy, row)

		from_day, _ := time.Parse("2006-01-02", from)
		f_day_n := int(from_day.Weekday())

		to_day, _ := time.Parse("2006-01-02", to)
		t_day_n := int(to_day.Weekday())
		var l int

		// OD
		if v < f_day_n {
			l = 7 - (f_day_n - v)
		} else {
			l = v - f_day_n
		}
		// DO
		var m int
		if v > t_day_n {
			m = (v - t_day_n) - 7
		} else {
			m = t_day_n - v
		}

		cpy[8] = strings.Repeat(".", v-1) + strconv.Itoa(v) + strings.Repeat(".", 7-v)
		cpy[6] = string(from_day.AddDate(0, 0, l).Format("2006-01-02"))
		cpy[7] = string(to_day.AddDate(0, 0, m).Format("2006-01-02"))

		newRows = append(newRows, cpy)

		// OD -> do przodu, DO do tylu, sprawdzenie dat
	}
	return newRows
}

func operatorToICAO(operator string) string {
	operatorMap := map[string]string{
		"2L": "OAW",
		"BT": "BTI",
		"LX": "SWR",
		"CL": "CLH",
		"EN": "DLA",
		"LH": "DLH",
		"OS": "AUA",
		"SN": "BEL",
	}

	_, exists := operatorMap[operator]

	if !exists {
		return operator
	} else {
		return operatorMap[operator]
	}
}
