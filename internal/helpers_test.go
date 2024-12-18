package internal

import (
	"bytes"
	"testing"
)

func TestFlattenJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "Valid JSON with one pair of ']['",
			input:    []byte(`[1][2][3]`),
			expected: []byte(`[1,2,3]`),
		},
		{
			name:     "Valid JSON with multiple pairs of ']['",
			input:    []byte(`[1][2][3][4][5]`),
			expected: []byte(`[1,2,3,4,5]`),
		},
		{
			name:     "No '][' in JSON",
			input:    []byte(`[1,2,3]`),
			expected: []byte(`[1,2,3]`), // No modification should occur
		},
		{
			name:     "Empty input",
			input:    []byte(``),
			expected: []byte(``), // Should return empty byte slice
		},
		{
			name:     "Input has empty object '[]'",
			input:    []byte(`[1][][2]`),
			expected: []byte(`[1],[2]`),
		},
		{
			name:     "Input ends with ']]'",
			input:    []byte(`[1][2]]`),
			expected: []byte(`[1,2]]`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlattenJSON(tt.input)

			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Test %s failed: expected %s, got %s", tt.name, string(tt.expected), string(result))
			}
		})
	}
}

func TestgetMonthMap(t *testing.T) {
	tests := []struct {
		name     string
		expected map[string]string
	}{
		{
			name: "Check if month map is correctly populated",
			expected: map[string]string{
				"JAN": "01",
				"FEB": "02",
				"MAR": "03",
				"APR": "04",
				"MAY": "05",
				"JUN": "06",
				"JUL": "07",
				"AUG": "08",
				"SEP": "09",
				"OCT": "10",
				"NOV": "11",
				"DEC": "12",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMonthMap()

			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("test %s failed: expected %s for key %s, got %s", tt.name, expectedValue, key, result[key])
				}
			}
		})
	}
}

func TestSSIMToDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid 6-character SSIM (single-digit day)",
			input:    "4JUL24",
			expected: "2024-07-04",
		},
		{
			name:     "Valid 7-character SSIM (two-digit day)",
			input:    "19JUL24",
			expected: "2024-07-19",
		},
		{
			name:     "Valid 6-character SSIM (single-digit day, different month)",
			input:    "1DEC24",
			expected: "2024-12-01",
		},
		{
			name:     "Valid 7-character SSIM (two-digit day, different month)",
			input:    "25DEC24",
			expected: "2024-12-25",
		},
		{
			name:     "Valid 7-character SSIM (February)",
			input:    "28FEB24",
			expected: "2024-02-28",
		},
		{
			name:     "Invalid input (incorrect month abbreviation)",
			input:    "19XYZ24",
			expected: "", // Expect empty string or error handling
		},
		{
			name:     "Invalid input (incorrect length)",
			input:    "JUL24",
			expected: "", // Expect empty string or error handling
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SSIMtoDate(tt.input)
			if result != tt.expected {
				t.Errorf("Test %s failed: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}

func TestDateToSSIM(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid Date (single-digit day)",
			input:    "2024-07-04",
			expected: "04JUL24",
		},
		{
			name:     "Valid Date (two-digit day)",
			input:    "2024-07-19",
			expected: "19JUL24",
		},
		{
			name:     "Valid Date (December)",
			input:    "2024-12-01",
			expected: "01DEC24",
		},
		{
			name:     "Valid Date (February)",
			input:    "2024-02-28",
			expected: "28FEB24",
		},
		{
			name:     "Valid Date (September)",
			input:    "2024-09-15",
			expected: "15SEP24",
		},
		{
			name:     "Invalid Date (incorrect month format)",
			input:    "2024-13-15",
			expected: "", // Expect empty string or panic based on error handling
		},
		{
			name:     "Invalid Date (incorrect length)",
			input:    "2024-7-4",
			expected: "", // Expect empty string or panic based on error handling
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && r != tt.expected {
					t.Errorf("Test %s failed: expected %v, got panic %v", tt.name, tt.expected, r)
				}
			}()

			result := DateToSSIM(tt.input)
			if result != tt.expected {
				t.Errorf("Test %s failed: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}

func TestNumberToTime(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "Exact hour",
			input:    60,
			expected: "01:00",
		},
		{
			name:     "Exact hour and minute",
			input:    90,
			expected: "01:30",
		},
		{
			name:     "Multiple hours",
			input:    150,
			expected: "02:30",
		},
		{
			name:     "No hours, only minutes",
			input:    45,
			expected: "00:45",
		},
		{
			name:     "Single digit hours and minutes",
			input:    9,
			expected: "00:09",
		},
		{
			name:     "Midnight",
			input:    0,
			expected: "00:00",
		},
		{
			name:     "Large value",
			input:    12345,
			expected: "205:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NumberToTime(tt.input)
			if result != tt.expected {
				t.Errorf("Test %s failed: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}

func TestDaysOfOperation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Days with spaces (Mon = 1, Tue = 2, etc.)",
			input:    "1  4567",
			expected: "1..4567", // Multiple spaces should be replaced with periods
		},
		{
			name:     "Days with single space",
			input:    "23 7",
			expected: ".23...7", // A space between 23 and 7 should become a single period
		},
		{
			name:     "Only single day code",
			input:    "1......",
			expected: "1......", // No space to replace, should remain the same
		},
		{
			name:     "Days with multiple consecutive spaces",
			input:    "1  2   3",
			expected: "1..2...3", // Multiple spaces between days replaced by periods
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "", // No change expected
		},
		{
			name:     "String with no spaces, all days",
			input:    "1234567",
			expected: "1234567", // No spaces, no change expected
		},
		{
			name:     "String with leading and trailing spaces",
			input:    " 23  7 ",
			expected: ".23...7", // Leading and trailing spaces replaced by periods
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DaysOfOperation(tt.input)
			if result != tt.expected {
				t.Errorf("Test %s failed: expected %s, got %s", tt.name, tt.expected, result)
			}
		})
	}
}

func TestmoreThanOneDigit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Single digit",
			input:    "1......",
			expected: false,
		},
		{
			name:     "Two digits",
			input:    "12.....",
			expected: true,
		},
		{
			name:     "Two digits on edges",
			input:    "1.....7",
			expected: true,
		},
		{
			name:     "All digits",
			input:    "1234567",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := moreThanOneNumberReg(tt.input)

			if result != tt.expected {
				t.Errorf("Test %s failed: expected: %v, got: %v", tt.name, tt.expected, result)
			}
		})
	}
}
func TestSeparateDays(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected [][]string
	}{
		{
			name:     "Valid days with digits",
			input:    []string{"Test", "data", "1234567"},
			expected: [][]string{{"Test", "data", "1"}, {"Test", "data", "2"}, {"Test", "data", "3"}, {"Test", "data", "4"}, {"Test", "data", "5"}, {"Test", "data", "6"}, {"Test", "data", "7"}},
		},
		{
			name:     "Single day, one day in the middle",
			input:    []string{"Test", "data", "12"},
			expected: [][]string{{"Test", "data", "1"}, {"Test", "data", "2"}},
		},
		{
			name:     "No days to separate",
			input:    []string{"Test", "data", "0"},
			expected: [][]string{{"Test", "data", "0"}},
		},
		{
			name:     "Multiple days, scattered",
			input:    []string{"Test", "data", "1357"},
			expected: [][]string{{"Test", "data", "1"}, {"Test", "data", "3"}, {"Test", "data", "5"}, {"Test", "data", "7"}},
		},
		{
			name:     "Invalid format, non-digit",
			input:    []string{"Test", "data", "abc"},
			expected: [][]string{{"Test", "data", "abc"}}, // Shouldn't be separated, invalid days string
		},
		{
			name:     "Empty input",
			input:    []string{"", "", ""},
			expected: [][]string{{"", "", ""}}, // No days to separate, just return as is
		},
		{
			name:     "Edge case with mixed valid and invalid days",
			input:    []string{"Test", "data", "1a3"},
			expected: [][]string{{"Test", "data", "1a3"}}, // Should not separate due to invalid input
		},
		{
			name:     "Single day with multiple occurrences",
			input:    []string{"Test", "data", "222"},
			expected: [][]string{{"Test", "data", "2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SeparateDays(tt.input)
			if !equal(result, tt.expected) {
				t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expected, result)
			}
		})
	}
}

// Helper function to check equality of 2D slices
func equal(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func TestPerformSeparation(t *testing.T) {
	tests := []struct {
		name     string
		row      []string
		days     []int
		expected [][]string
	}{
		{
			name: "Valid case with normal dates",
			row:  []string{"data1", "data2", "data3", "data4", "data5", "data6", "2023-05-01", "2023-05-07", "......."},
			days: []int{1, 3, 5},
			expected: [][]string{
				{"data1", "data2", "data3", "data4", "data5", "data6", "2023-05-01", "2023-05-01", "..1......"},
				{"data1", "data2", "data3", "data4", "data5", "data6", "2023-05-03", "2023-05-03", "...3....."},
				{"data1", "data2", "data3", "data4", "data5", "data6", "2023-05-05", "2023-05-05", "....5...."},
			},
		},
		{
			name: "Edge case with same start and end day",
			row:  []string{"data1", "data2", "data3", "data4", "data5", "data6", "2023-05-01", "2023-05-01", "......."},
			days: []int{1, 7},
			expected: [][]string{
				{"data1", "data2", "data3", "data4", "data5", "data6", "2023-05-01", "2023-05-01", "..1......"},
				{"data1", "data2", "data3", "data4", "data5", "data6", "2023-05-01", "2023-05-01", "......7.."},
			},
		},
		{
			name: "Case with backward day calculation",
			row:  []string{"data1", "data2", "data3", "data4", "data5", "data6", "2023-05-07", "2023-05-01", "......."},
			days: []int{6},
			expected: [][]string{
				{"data1", "data2", "data3", "data4", "data5", "data6", "2023-05-06", "2023-05-06", ".....6..."},
			},
		},
		{
			name: "Invalid date format in row",
			row:  []string{"data1", "data2", "data3", "data4", "data5", "data6", "invalid-date", "2023-05-07", "......."},
			days: []int{1, 3, 5},
			expected: [][]string{
				{"data1", "data2", "data3", "data4", "data5", "data6", "invalid-date", "2023-05-07", "......."},
			}, // Should return the original row with invalid date without processing
		},
		{
			name: "Case with empty row",
			row:  []string{"", "", "", "", "", "", "", "", ""},
			days: []int{1},
			expected: [][]string{
				{"", "", "", "", "", "", "", "", ""},
			}, // Empty row should be returned as is
		},
		{
			name: "Edge case with no days",
			row:  []string{"data1", "data2", "data3", "data4", "data5", "data6", "2023-05-01", "2023-05-07", "......."},
			days: []int{},
			expected: [][]string{
				{"data1", "data2", "data3", "data4", "data5", "data6", "2023-05-01", "2023-05-07", "......."},
			}, // No days, should return the row as is
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := performSeparation(tt.row, tt.days)
			if !equal2d(result, tt.expected) {
				t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expected, result)
			}
		})
	}
}

// Helper function to check equality of 2D slices
func equal2d(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
func TestOperatorToICAO(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		expected string
	}{
		{
			name:     "Valid operator 2L",
			operator: "2L",
			expected: "OAW",
		},
		{
			name:     "Valid operator BT",
			operator: "BT",
			expected: "BTI",
		},
		{
			name:     "Valid operator LX",
			operator: "LX",
			expected: "SWR",
		},
		{
			name:     "Valid operator CL",
			operator: "CL",
			expected: "CLH",
		},
		{
			name:     "Invalid operator",
			operator: "XX",
			expected: "XX", // Unmapped operator should return as is
		},
		{
			name:     "Another invalid operator",
			operator: "ZZ",
			expected: "ZZ", // Unmapped operator should return as is
		},
		{
			name:     "Valid operator LH",
			operator: "LH",
			expected: "DLH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := operatorToICAO(tt.operator)
			if result != tt.expected {
				t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expected, result)
			}
		})
	}
}
func TestGetQueryListForAirline(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		beg      string
		end      string
		expected []ApiQuery
	}{
		{
			name: "Test code 0",
			code: 0,
			beg:  "2024-01-01",
			end:  "2024-01-31",
			expected: []ApiQuery{
				{
					Airline:         "LH",
					StartDate:       "2024-01-01",
					EndDate:         "2024-01-31",
					DaysOfOperation: "1234567",
					TimeMode:        "LT",
					Origin:          "KRK",
					Destination:     "FRA",
				},
				{
					Airline:         "LH",
					StartDate:       "2024-01-01",
					EndDate:         "2024-01-31",
					DaysOfOperation: "1234567",
					TimeMode:        "LT",
					Origin:          "KRK",
					Destination:     "MUC",
				},
			},
		},
		{
			name: "Test code 1",
			code: 1,
			beg:  "2024-01-01",
			end:  "2024-01-31",
			expected: []ApiQuery{
				{
					Airline:         "OS",
					StartDate:       "2024-01-01",
					EndDate:         "2024-01-31",
					DaysOfOperation: "1234567",
					TimeMode:        "LT",
					Origin:          "KRK",
					Destination:     "VIE",
				},
			},
		},
		{
			name: "Test code 2",
			code: 2,
			beg:  "2024-01-01",
			end:  "2024-01-31",
			expected: []ApiQuery{
				{
					Airline:         "LX",
					StartDate:       "2024-01-01",
					EndDate:         "2024-01-31",
					DaysOfOperation: "1234567",
					TimeMode:        "LT",
					Origin:          "KRK",
					Destination:     "ZRH",
				},
			},
		},
		{
			name: "Test code 3",
			code: 3,
			beg:  "2024-01-01",
			end:  "2024-01-31",
			expected: []ApiQuery{
				{
					Airline:         "SN",
					StartDate:       "2024-01-01",
					EndDate:         "2024-01-31",
					DaysOfOperation: "1234567",
					TimeMode:        "LT",
					Origin:          "KRK",
					Destination:     "BRU",
				},
			},
		},
		{
			name: "Test code 4",
			code: 4,
			beg:  "2024-01-01",
			end:  "2024-01-31",
			expected: []ApiQuery{
				{
					Airline:         "EN",
					StartDate:       "2024-01-01",
					EndDate:         "2024-01-31",
					DaysOfOperation: "1234567",
					TimeMode:        "LT",
					Origin:          "KRK",
					Destination:     "MUC",
				},
			},
		},
		{
			name:     "Test default case",
			code:     99, // Invalid code
			beg:      "2024-01-01",
			end:      "2024-01-31",
			expected: []ApiQuery{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetQueryListForAirline(tt.code, tt.beg, tt.end)
			if len(result) != len(tt.expected) {
				t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expected, result)
				return
			}

			// Compare each field of the ApiQuery struct
			for i, query := range result {
				if query != tt.expected[i] {
					t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expected[i], query)
				}
			}
		})
	}
}
func TestAreValidForMerge(t *testing.T) {
	tests := []struct {
		name      string
		record1   []string
		record2   []string
		expected  bool
		expectErr bool
	}{
		{
			name: "Valid merge",
			record1: []string{
				"1", "ABC", "XYZ", "123", "456", "789", "2024-01-01", "2024-01-08", "Mon", "Tue", "Wed", "Thu",
			},
			record2: []string{
				"1", "ABC", "XYZ", "123", "456", "789", "2024-01-08", "2024-01-15", "Mon", "Tue", "Wed", "Thu",
			},
			expected:  true,
			expectErr: false,
		},
		{
			name: "Invalid merge - different columns",
			record1: []string{
				"1", "ABC", "XYZ", "123", "456", "789", "2024-01-01", "2024-01-08", "Mon", "Tue", "Wed", "Thu",
			},
			record2: []string{
				"2", "DEF", "XYZ", "123", "456", "789", "2024-01-08", "2024-01-15", "Mon", "Tue", "Wed", "Thu",
			},
			expected:  false,
			expectErr: false,
		},
		{
			name: "Invalid merge - different dates",
			record1: []string{
				"1", "ABC", "XYZ", "123", "456", "789", "2024-01-01", "2024-01-08", "Mon", "Tue", "Wed", "Thu",
			},
			record2: []string{
				"1", "ABC", "XYZ", "123", "456", "789", "2024-01-10", "2024-01-15", "Mon", "Tue", "Wed", "Thu",
			},
			expected:  false,
			expectErr: false,
		},
		{
			name: "Invalid date format",
			record1: []string{
				"1", "ABC", "XYZ", "123", "456", "789", "2024-01-01", "2024-01-08", "Mon", "Tue", "Wed", "Thu",
			},
			record2: []string{
				"1", "ABC", "XYZ", "123", "456", "789", "invalid-date", "2024-01-15", "Mon", "Tue", "Wed", "Thu",
			},
			expected:  false,
			expectErr: true,
		},
		{
			name: "Same data, same dates, valid merge",
			record1: []string{
				"1", "ABC", "XYZ", "123", "456", "789", "2024-01-01", "2024-01-08", "Mon", "Tue", "Wed", "Thu",
			},
			record2: []string{
				"1", "ABC", "XYZ", "123", "456", "789", "2024-01-08", "2024-01-15", "Mon", "Tue", "Wed", "Thu",
			},
			expected:  true,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AreValidForMerge(tt.record1, tt.record2)

			// Check if error state matches expectations
			if (err != nil) != tt.expectErr {
				t.Errorf("Test %s failed: expected error: %v, got error: %v", tt.name, tt.expectErr, err)
			}

			// Compare the result to the expected value
			if result != tt.expected {
				t.Errorf("Test %s failed: expected: %v, got: %v", tt.name, tt.expected, result)
			}
		})
	}
}
func TestPerformMerge(t *testing.T) {
	tests := []struct {
		name     string
		record1  []string
		record2  []string
		expected []string
	}{
		{
			name: "Merge records with matching columns and updated date",
			record1: []string{
				"1", "ABC", "XYZ", "123", "456", "789", "2024-01-01", "2024-01-08", "Mon", "Tue", "Wed", "Thu",
			},
			record2: []string{
				"1", "ABC", "XYZ", "123", "456", "789", "2024-01-08", "2024-01-15", "Mon", "Tue", "Wed", "Thu",
			},
			expected: []string{
				"1", "ABC", "XYZ", "123", "456", "789", "2024-01-01", "2024-01-15", "Mon", "Tue", "Wed", "Thu",
			},
		},
		{
			name: "Merge records with different data but same structure",
			record1: []string{
				"2", "DEF", "ABC", "111", "222", "333", "2024-01-01", "2024-01-08", "Mon", "Tue", "Wed", "Thu",
			},
			record2: []string{
				"2", "DEF", "ABC", "111", "222", "333", "2024-01-08", "2024-01-15", "Mon", "Tue", "Wed", "Thu",
			},
			expected: []string{
				"2", "DEF", "ABC", "111", "222", "333", "2024-01-01", "2024-01-15", "Mon", "Tue", "Wed", "Thu",
			},
		},
		{
			name: "Merge records with identical dates",
			record1: []string{
				"3", "XYZ", "DEF", "456", "789", "101", "2024-02-01", "2024-02-08", "Mon", "Tue", "Wed", "Thu",
			},
			record2: []string{
				"3", "XYZ", "DEF", "456", "789", "101", "2024-02-08", "2024-02-15", "Mon", "Tue", "Wed", "Thu",
			},
			expected: []string{
				"3", "XYZ", "DEF", "456", "789", "101", "2024-02-01", "2024-02-15", "Mon", "Tue", "Wed", "Thu",
			},
		},
		{
			name: "Merge with unchanged date but different data",
			record1: []string{
				"4", "GHI", "JKL", "555", "666", "777", "2024-03-01", "2024-03-08", "Mon", "Tue", "Wed", "Thu",
			},
			record2: []string{
				"4", "GHI", "JKL", "555", "666", "777", "2024-03-08", "2024-03-15", "Mon", "Tue", "Wed", "Thu",
			},
			expected: []string{
				"4", "GHI", "JKL", "555", "666", "777", "2024-03-01", "2024-03-15", "Mon", "Tue", "Wed", "Thu",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PerformMerge(tt.record1, tt.record2)

			// Check if result matches expected
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("Test %s failed: expected: %v, got: %v", tt.name, tt.expected, result)
					break
				}
			}
		})
	}
}
func TestConvertFlightResponseToCSVRows(t *testing.T) {
	tests := []struct {
		name     string
		input    FlightResponse
		expected [][]string
	}{
		{
			name: "Basic Flight Response",
			input: FlightResponse{
				Airline:      "LH",
				FlightNumber: 123,
				Legs: []Leg{
					{
						Origin:                  "KRK",
						Destination:             "FRA",
						AircraftOwner:           "OAW",
						AircraftType:            "A320",
						ServiceType:             "Regular",
						AircraftDepartureTimeLT: 1590000000,
						AircraftArrivalTimeLT:   1590015000,
					},
				},
				PeriodOfOperationLT: PeriodOfOperation{
					StartDate:       "1JAN24",
					EndDate:         "31JAN24",
					DaysOfOperation: "1234567",
				},
			},
			expected: [][]string{
				{
					"KRK", "FRA", "LH", "123", "02:00", "02:30", "2024-01-01", "2024-01-31", "1234567", "A320", "OAW", "Regular",
				},
			},
		},
		{
			name: "Flight with different airline and time",
			input: FlightResponse{
				Airline:      "LX",
				FlightNumber: 456,
				Legs: []Leg{
					{
						Origin:                  "MUC",
						Destination:             "ZRH",
						AircraftOwner:           "SWR",
						AircraftType:            "A321",
						ServiceType:             "VIP",
						AircraftDepartureTimeLT: 1600000000,
						AircraftArrivalTimeLT:   1600015000,
					},
				},
				PeriodOfOperationLT: PeriodOfOperation{
					StartDate:       "5FEB24",
					EndDate:         "10FEB24",
					DaysOfOperation: "1357",
				},
			},
			expected: [][]string{
				{
					"MUC", "ZRH", "LX", "456", "03:00", "03:30", "2024-02-05", "2024-02-10", "1357", "A321", "SWR", "VIP",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertFlightResponseToCSVRows(tt.input)

			// Compare result with expected
			if len(result) != len(tt.expected) {
				t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expected, result)
				return
			}

			for i, row := range result {
				for j, val := range row {
					if val != tt.expected[i][j] {
						t.Errorf("Test %s failed: expected %v at row %d, column %d, got %v", tt.name, tt.expected[i][j], i, j, val)
					}
				}
			}
		})
	}
}
func TestMergeRecords(t *testing.T) {
	tests := []struct {
		name      string
		input     [][]string
		expected  [][]string
		expectErr bool
	}{
		{
			name: "Valid merge of two records",
			input: [][]string{
				{"KRK", "FRA", "LH", "123", "02:00", "02:30", "2024-01-01", "2024-01-31", "1234567", "A320", "OAW", "Regular"},
				{"KRK", "FRA", "LH", "123", "02:00", "02:30", "2024-01-01", "2024-01-31", "1234567", "A320", "OAW", "Regular"},
			},
			expected: [][]string{
				{"KRK", "FRA", "LH", "123", "02:00", "02:30", "2024-01-01", "2024-01-31", "1234567", "A320", "OAW", "Regular"},
			},
			expectErr: false,
		},
		{
			name: "Non-mergeable records",
			input: [][]string{
				{"KRK", "FRA", "LH", "123", "02:00", "02:30", "2024-01-01", "2024-01-31", "1234567", "A320", "OAW", "Regular"},
				{"MUC", "ZRH", "LX", "456", "03:00", "03:30", "2024-02-01", "2024-02-28", "1357", "A321", "SWR", "VIP"},
			},
			expected: [][]string{
				{"KRK", "FRA", "LH", "123", "02:00", "02:30", "2024-01-01", "2024-01-31", "1234567", "A320", "OAW", "Regular"},
				{"MUC", "ZRH", "LX", "456", "03:00", "03:30", "2024-02-01", "2024-02-28", "1357", "A321", "SWR", "VIP"},
			},
			expectErr: false,
		},
		{
			name: "Merge records with date overlap",
			input: [][]string{
				{"KRK", "FRA", "LH", "123", "02:00", "02:30", "2024-01-01", "2024-01-31", "1234567", "A320", "OAW", "Regular"},
				{"KRK", "FRA", "LH", "123", "02:00", "02:30", "2024-01-15", "2024-01-31", "1234567", "A320", "OAW", "Regular"},
			},
			expected: [][]string{
				{"KRK", "FRA", "LH", "123", "02:00", "02:30", "2024-01-01", "2024-01-31", "1234567", "A320", "OAW", "Regular"},
			},
			expectErr: false,
		},
		{
			name: "Merge with invalid records",
			input: [][]string{
				{"KRK", "FRA", "LH", "123", "02:00", "02:30", "2024-01-01", "2024-01-31", "1234567", "A320", "OAW", "Regular"},
				{"KRK", "FRA", "LH", "123", "02:00", "02:30", "INVALID_DATE", "2024-01-31", "1234567", "A320", "OAW", "Regular"},
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MergeRecords(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
				return
			}

			for i, row := range result {
				for j, val := range row {
					if val != tt.expected[i][j] {
						t.Errorf("at row %d, column %d: expected %v, got %v", i, j, tt.expected[i][j], val)
					}
				}
			}
		})
	}
}
func TestSortRecordsByDateCol(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]string
		expected [][]string
	}{
		{
			name: "Valid dates in ascending order",
			input: [][]string{
				{"KRK", "FRA", "2024-01-03"},
				{"KRK", "FRA", "2024-01-01"},
				{"KRK", "FRA", "2024-01-02"},
			},
			expected: [][]string{
				{"KRK", "FRA", "2024-01-01"},
				{"KRK", "FRA", "2024-01-02"},
				{"KRK", "FRA", "2024-01-03"},
			},
		},
		{
			name: "Valid and invalid dates",
			input: [][]string{
				{"KRK", "FRA", "2024-01-03"},
				{"KRK", "FRA", "INVALID_DATE"},
				{"KRK", "FRA", "2024-01-01"},
			},
			expected: [][]string{
				{"KRK", "FRA", "2024-01-01"},
				{"KRK", "FRA", "2024-01-03"},
				{"KRK", "FRA", "INVALID_DATE"},
			},
		},
		{
			name: "Already sorted records",
			input: [][]string{
				{"KRK", "FRA", "2024-01-01"},
				{"KRK", "FRA", "2024-01-02"},
				{"KRK", "FRA", "2024-01-03"},
			},
			expected: [][]string{
				{"KRK", "FRA", "2024-01-01"},
				{"KRK", "FRA", "2024-01-02"},
				{"KRK", "FRA", "2024-01-03"},
			},
		},
		{
			name:     "Empty list",
			input:    [][]string{},
			expected: [][]string{},
		},
		{
			name: "Single record",
			input: [][]string{
				{"KRK", "FRA", "2024-01-01"},
			},
			expected: [][]string{
				{"KRK", "FRA", "2024-01-01"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SortRecordsByDateCol(tt.input, 2)

			// Check if the result matches the expected output
			for i := range tt.input {
				for j := range tt.input[i] {
					if tt.input[i][j] != tt.expected[i][j] {
						t.Errorf("Test %s failed: expected %v, got %v", tt.name, tt.expected, tt.input)
					}
				}
			}
		})
	}
}
