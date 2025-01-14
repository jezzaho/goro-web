package internal

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func FlattenJSON(data []byte) []byte {
	// Convert input to string for easier manipulation
	strData := string(data)

	// Remove all empty arrays `[]` and replace them with a comma
	strData = strings.ReplaceAll(strData, "[]", ",")

	// Replace all occurrences of `][` with a comma
	strData = strings.ReplaceAll(strData, "][", ",")

	// Clean up any leftover double commas `,,` caused by replacements
	strData = strings.ReplaceAll(strData, ",,", ",")

	// Trim leading and trailing commas if they exist
	strData = strings.Trim(strData, ",")

	return []byte(strData)
}

func getMonthMap() map[string]string {
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

	return monthMap
}

func SSIMtoDate(s string) string {
	monthMap := getMonthMap()

	//var dateString string
	var year string
	var month string
	var day string
	length := len(s)
	switch length {
	case 6:
		// 4JUL24 012345
		//dateString += "20" + s[4:] + "-" + monthMap[s[1:4]] + "-0" + s[:1]
		year = "20" + s[4:]
		month = monthMap[s[1:4]]
		day = "0" + s[:1]
	case 7:
		// 19JUL24 0123456
		//dateString += "20" + s[5:] + "-" + monthMap[s[2:5]] + "-" + s[:2]
		year = "20" + s[5:]
		month = monthMap[s[2:5]]
		day = s[:2]
	}
	if month == "" {
		return ""
	}
	return year + "-" + month + "-" + day
}

func DateToSSIM(s string) string {
	// FORMAT
	// YYYY-MM-DD
	year := s[2:4]
	month := s[5:7]
	day := s[8:]

	var monthStr string
	foundFlag := false
	monthMap := getMonthMap()
	// Iterate over the map to find the key for the value
	for key, value := range monthMap {
		if value == month {
			monthStr = key
			foundFlag = true
			break
		}
	}
	if !foundFlag {
		return ""
	}
	return day + monthStr + year

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

func performSeparation(record []string, weekdays []int) [][]string {
	var separatedRecords [][]string

	startDate := record[6]
	endDate := record[7]

	startDateTime, _ := time.Parse("2006-01-02", startDate)
	endDateTime, _ := time.Parse("2006-01-02", endDate)

	startWeekday := int(startDateTime.Weekday())
	if startWeekday == 0 {
		startWeekday = 7
	}
	endWeekday := int(endDateTime.Weekday())
	if endWeekday == 0 {
		endWeekday = 7
	}
	for _, targetWeekday := range weekdays {
		// Create a copy of the original record
		newRecord := make([]string, len(record))
		copy(newRecord, record)

		// Calculate days to adjust for start date
		daysToAdjustStart := 0
		if targetWeekday < startWeekday {
			daysToAdjustStart = 7 - (startWeekday - targetWeekday)
		} else if targetWeekday > startWeekday {
			daysToAdjustStart = targetWeekday - startWeekday
		}

		// Calculate days to adjust for end date
		daysToAdjustEnd := 0
		if targetWeekday > endWeekday {
			daysToAdjustEnd = (targetWeekday - endWeekday) - 7
		} else if targetWeekday < endWeekday {
			daysToAdjustEnd = targetWeekday - endWeekday
		}

		// Create weekday marker (e.g., "...3....")
		weekdayMarker := strings.Repeat(".", targetWeekday-1) +
			strconv.Itoa(targetWeekday) +
			strings.Repeat(".", 7-targetWeekday)

		// Update the record with new values
		newRecord[8] = weekdayMarker
		if daysToAdjustStart != 0 {
			newRecord[6] = startDateTime.AddDate(0, 0, daysToAdjustStart).Format("2006-01-02")
		}
		if daysToAdjustEnd != 0 {
			newRecord[7] = endDateTime.AddDate(0, 0, daysToAdjustEnd).Format("2006-01-02")
		}

		separatedRecords = append(separatedRecords, newRecord)
	}

	return separatedRecords
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

// Querying for specific Airline should output specyfic Querylist
// Code: 0 - LH || 1 - OS || 2 - LX || 3 - SN || 4 - EN
// beg && end format in SSIM date format DDMMMYY eg. 15MAR25
func GetQueryListForAirline(code int, beg, end string) (QueryList []ApiQuery) {
	switch code {
	case 0:
		return []ApiQuery{
			{
				Airline:         "LH",
				StartDate:       beg,
				EndDate:         end,
				DaysOfOperation: "1234567",
				TimeMode:        "LT",
				Origin:          "KRK",
				Destination:     "FRA",
			},
			{
				Airline:         "LH",
				StartDate:       beg,
				EndDate:         end,
				DaysOfOperation: "1234567",
				TimeMode:        "LT",
				Origin:          "KRK",
				Destination:     "MUC",
			},
		}
	case 1:
		return []ApiQuery{
			{
				Airline:         "OS",
				StartDate:       beg,
				EndDate:         end,
				DaysOfOperation: "1234567",
				TimeMode:        "LT",
				Origin:          "KRK",
				Destination:     "VIE",
			},
		}
	case 2:
		return []ApiQuery{
			{
				Airline:         "LX",
				StartDate:       beg,
				EndDate:         end,
				DaysOfOperation: "1234567",
				TimeMode:        "LT",
				Origin:          "KRK",
				Destination:     "ZRH",
			},
		}
	case 3:
		return []ApiQuery{
			{
				Airline:         "SN",
				StartDate:       beg,
				EndDate:         end,
				DaysOfOperation: "1234567",
				TimeMode:        "LT",
				Origin:          "KRK",
				Destination:     "BRU",
			},
		}
	case 4:
		return []ApiQuery{
			{
				Airline:         "EN",
				StartDate:       beg,
				EndDate:         end,
				DaysOfOperation: "1234567",
				TimeMode:        "LT",
				Origin:          "KRK",
				Destination:     "MUC",
			},
		}
		// Change default?
	default:
		return []ApiQuery{}
	}
}

func AreValidForMerge(record1, record2 []string) (bool, error) {
	columnsToCompare := []int{0, 1, 2, 3, 4, 5, 8, 9, 10, 11}

	for _, col := range columnsToCompare {
		if record1[col] != record2[col] {

			return false, nil
		}
	}
	// Compare dates
	record1To := record1[7]
	record2From := record2[6]

	dateOne, err := time.Parse("2006-01-02", record1To)
	if err != nil {
		return false, err
	}
	dateOne = dateOne.AddDate(0, 0, 7)

	dateTwo, err := time.Parse("2006-01-02", record2From)
	if err != nil {
		return false, err
	}

	if dateOne.Compare(dateTwo) == 0 {
		return true, nil
	}
	return false, nil

}

func PerformMerge(record1, record2 []string) []string {
	temp := make([]string, len(record1))
	copy(temp, record1)

	temp[7] = record2[7]

	return temp
}

// Helper function to convert FlightResponse to CSV rows
func convertFlightResponseToCSVRows(d FlightResponse) [][]string {
	var csvRows [][]string

	flightNumberWrite := strconv.Itoa(d.FlightNumber)
	startTimeWrite := NumberToTime(d.Legs[0].AircraftDepartureTimeLT)
	endTimeWrite := NumberToTime(d.Legs[0].AircraftArrivalTimeLT)
	startDateWrite := SSIMtoDate(d.PeriodOfOperationLT.StartDate)
	endDateWrite := SSIMtoDate(d.PeriodOfOperationLT.EndDate)
	daysOfOperationWrite := DaysOfOperation(d.PeriodOfOperationLT.DaysOfOperation)
	operator := operatorToICAO(d.Legs[0].AircraftOwner)

	row := []string{
		d.Legs[0].Origin,
		d.Legs[0].Destination,
		d.Airline,
		flightNumberWrite,
		startTimeWrite,
		endTimeWrite,
		startDateWrite,
		endDateWrite,
		daysOfOperationWrite,
		d.Legs[0].AircraftType,
		operator,
		d.Legs[0].ServiceType,
	}

	csvRows = append(csvRows, row)
	return csvRows
}

func MergeRecords(records [][]string) ([][]string, error) {
	if len(records) <= 1 {
		return records, nil
	}

	mergedRecords := [][]string{} // Resulting list of merged records
	merged := make(map[int]bool)  // Tracks indices of merged records

	for i := 0; i < len(records); i++ {
		if merged[i] { // Skip already merged records
			continue
		}

		currentRecord := records[i]
		for j := 0; j < len(records); j++ {
			if i == j || merged[j] { // Skip self and already merged records
				continue
			}

			// Check if the records can be merged
			validFlag, err := AreValidForMerge(currentRecord, records[j])
			if err != nil {
				return nil, err
			}

			if validFlag {
				currentRecord = PerformMerge(currentRecord, records[j]) // Merge the records
				merged[j] = true                                        // Mark the merged record
			}
		}

		// Add the merged record to the result
		mergedRecords = append(mergedRecords, currentRecord)
	}

	return mergedRecords, nil
}

func SortRecordsByDateCol(data [][]string, columnIndex int) {
	sort.Slice(data, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", data[i][columnIndex])
		dateJ, errJ := time.Parse("2006-01-02", data[j][columnIndex])

		// If parsing fails, maintain original order
		if errI != nil || errJ != nil {
			log.Printf("Parsing for Sorting Failed because: %v and %v", errI, errJ)
			return false
		}

		// Compare dates
		return dateI.Before(dateJ)
	})
}
