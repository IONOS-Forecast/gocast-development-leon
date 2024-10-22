package utils

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestSetDate(t *testing.T) {
	// Expecting no errors
	year := 2010
	month := 12
	day := 5
	expectedDate := fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	date, err := SetDate(year, month, day)
	if err != nil {
		t.Error("SetDate() returned an error: ", err)
	}
	if expectedDate != date {
		t.Errorf("SetDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	now := time.Now()
	year = now.Year()
	month = int(now.Month())
	day = now.Day()
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	date, err = SetDate(year, month, day)
	if err != nil {
		t.Error("SetDate() returned an error: ", err)
	}
	if expectedDate != date {
		t.Errorf("SetDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	// Expecting errors
	log.Print("INFO: 2 expected warnings following!")
	year = now.Year()
	month = int(now.Month())
	day = now.Day()
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	year = 2008
	month = 8
	day = 25
	expectedErrorDate := fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	expectedError := fmt.Errorf("date should not be before december 2010").Error()
	date, err = SetDate(year, month, day)
	if err == nil || err.Error() != expectedError {
		t.Errorf("SetDate() returned wrong error: got \"%v\" want \"%v\"", err, expectedError)
	}
	if expectedDate != date {
		t.Errorf("SetDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	year = 2010
	month = 2
	day = 29
	expectedErrorDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	expectedError = fmt.Errorf("parsing time \"%v\": day out of range", expectedErrorDate).Error()
	date, err = SetDate(year, month, day)
	if err == nil || err.Error() != expectedError {
		t.Errorf("SetDate() returned wrong error: got \"%v\" want \"%v\"", err, expectedError)
	}
	if expectedDate != date {
		t.Errorf("SetDate() returned wrong date: got %v want %v", date, expectedDate)
	}
}

func TestGetDate(t *testing.T) {
	year := 2010
	month := 12
	day := 5
	expectedDate, err := SetDate(year, month, day)
	if err != nil {
		t.Error("SetDate() in test getDateByString() returned an error: ", err)
	}
	if expectedDate != GetDate() {
		t.Errorf("getDateByString() returned wrong date: got %v want %v", date, expectedDate)
	}
}

func getFutureDate(year, month, day int) (string, string) {
	return fmt.Sprintf("%v-%.2v-%.2v", year, month, day), fmt.Sprintf("%v-%.2v-%.2v", year, month, day+1)
}

func TestSetFutureDate(t *testing.T) {
	testDate, expectedDate := getFutureDate(2011, 2, 27)
	date, _, err := SetFutureDay(expectedDate, testDate, "0")
	if err != nil {
		t.Error("SetFutureDate() returned an error: ", err)
	}
	if expectedDate != date {
		t.Errorf("SetFutureDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	testDate, expectedDate = getFutureDate(2015, 8, 1)
	date, _, err = SetFutureDay(expectedDate, testDate, "0")
	if err != nil {
		t.Error("SetFutureDate() returned an error: ", err)
	}
	if expectedDate != date {
		t.Errorf("SetFutureDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	t.Run("checking for unexpected errors", func(t *testing.T) {
		log.Print("INFO: 4 expected warnings following!")
		testDate, expectedDate = getFutureDate(2011, 2, 28)
		expectedError := fmt.Errorf("parsing time \"%v\": day out of range", expectedDate).Error()
		date, _, err = SetFutureDay(expectedDate, testDate, "0")
		if err == nil || err.Error() != expectedError {
			t.Errorf("SetFutureDate() returned wrong error: got \"%v\" want \"%v\"", err, expectedError)
		}
		if expectedDate == date {
			t.Errorf("SetFutureDate() returned wrong date: got %v want %v", date, expectedDate)
		}
		testDate, expectedDate = getFutureDate(2010, 12, 31)
		expectedError = fmt.Errorf("parsing time \"%v\": day out of range", expectedDate).Error()
		date, _, err = SetFutureDay(expectedDate, testDate, "0")
		if err == nil || err.Error() != expectedError {
			t.Errorf("SetFutureDate() returned wrong error: got \"%v\" want \"%v\"", err, expectedError)
		}
		if expectedDate == date {
			t.Errorf("SetFutureDate() returned wrong date: got %v want %v", date, expectedDate)
		}
		testDate, expectedDate = getFutureDate(2008, 8, 25)
		expectedError = fmt.Errorf("date should not be before december 2010").Error()
		date, _, err = SetFutureDay(expectedDate, testDate, "0")
		if err == nil || err.Error() != expectedError {
			t.Errorf("SetFutureDate() returned wrong error: got \"%v\" want \"%v\"", err, expectedError)
		}
		if expectedDate == date {
			t.Errorf("SetFutureDate() returned wrong date: got %v want %v", date, expectedDate)
		}
		now := time.Now()
		year := now.Year()
		month := int(now.Month())
		day := now.Day()
		testDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
		expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day+8)
		expectedError = fmt.Errorf("date should not be after next 7 days").Error()
		date, _, err = SetFutureDay(expectedDate, testDate, "7")
		if err == nil || err.Error() != expectedError {
			t.Errorf("SetFutureDate() returned wrong error: got \"%v\" want \"%v\"", err, expectedError)
		}
		if expectedDate == date {
			t.Errorf("SetFutureDate() returned wrong date: got %v want %v", date, testDate)
		}
	})
}

func TestSplitDate(t *testing.T) {
	year := 2010
	month := 12
	day := 5
	expectedDate := fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	year, month, day, err := SplitDate(expectedDate)
	if err != nil {
		t.Error("SplitDate() returned an error: ", err)
	}
	date := fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	if expectedDate != date {
		t.Errorf("SplitDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	year = 2020
	month = 5
	day = 12
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	year, month, day, err = SplitDate(expectedDate)
	if err != nil {
		t.Error("SplitDate() returned an error: ", err)
	}
	date = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	if expectedDate != date {
		t.Errorf("SplitDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	now := time.Now()
	year = now.Year()
	month = int(now.Month())
	day = now.Day()
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	year, month, day, err = SplitDate(expectedDate)
	if err != nil {
		t.Error("SplitDate() returned an error: ", err)
	}
	date = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	if expectedDate != date {
		t.Errorf("SplitDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	// Expecting errors
	t.Run("checking for unexpected errors", func(t *testing.T) {
		expectedDate = fmt.Sprintf("%v-%.2v-%.2v", 0, 0, 0)
		expectedError := fmt.Errorf("failed to convert day: strconv.Atoi: parsing \"%.2v\": invalid syntax", "five").Error()
		date = getDateByString("2018", "5", "five")
		year, month, day, err = SplitDate(date)
		if err == nil || err.Error() != expectedError {
			t.Errorf("SplitDate() returned unexpected error: got \"%v\" want \"%v\"", err, expectedError)
		}
		date = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
		if date != expectedDate {
			t.Errorf("SplitDate() returned unexpected date: got %v want %v", date, expectedDate)
		}
		date = getDateByString("2018", "five", "2")
		expectedError = fmt.Errorf("failed to convert month: strconv.Atoi: parsing \"%.2v\": invalid syntax", "five").Error()
		year, month, day, err = SplitDate(date)
		if err == nil || err.Error() != expectedError {
			t.Errorf("SplitDate() returned unexpected error: got \"%v\" want \"%v\"", err, expectedError)
		}
		date = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
		if date != expectedDate {
			t.Errorf("SplitDate() returned unexpected date: got %v want %v", date, expectedDate)
		}
		date = getDateByString("five", "5", "2")
		expectedError = fmt.Errorf("failed to convert year: strconv.Atoi: parsing \"%v\": invalid syntax", "five").Error()
		year, month, day, err = SplitDate(date)
		if err == nil || err.Error() != expectedError {
			t.Errorf("SplitDate() returned unexpected error: got \"%v\" want \"%v\"", err, expectedError)
		}
		date = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
		if date != expectedDate {
			t.Errorf("SplitDate() returned unexpected date: got %v want %v", date, expectedDate)
		}
	})
}

func getDateByString(year, month, day string) string {
	return fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
}

func TestCheckDate(t *testing.T) {
	tests := []struct {
		expectedDate  string
		wantErr       bool
		expectedError string
	}{
		{"2010-12-05", false, ""},
		{"2016-03-06", false, ""},
		{"2011-01-15", false, ""},
		{"2014-05-22", false, ""},
		{"2016-09-10", false, ""},
		{"2018-03-05", false, ""},
		{"2020-07-20", false, ""},
		{"2021-11-30", false, ""},
		{"2022-06-18", false, ""},
		{"2023-01-25", false, ""},
		{"2023-08-14", false, ""},
		{"2024-02-29", false, ""},
		{fmt.Sprintf("%v-%.2v-%.2v", time.Now().Year(), int(time.Now().Month()), time.Now().Day()), false, ""},
		{fmt.Sprintf("%v-%.2v-%.2v", time.Now().Year(), int(time.Now().Month()), time.Now().Day()+8), true, "date should not be after next 7 days"},
		{"2008-08-25", true, "date should not be before december 2010"},
		{"2010-11-15", true, "date should not be before december 2010"},
		{"2010-08-22", true, "date should not be before december 2010"},
		{"2009-05-10", true, "date should not be before december 2010"},
		{"2008-12-01", true, "date should not be before december 2010"},
		{"2007-07-19", true, "date should not be before december 2010"},
		{"2006-03-28", true, "date should not be before december 2010"},
		{"2005-09-15", true, "date should not be before december 2010"},
		{"2004-01-12", true, "date should not be before december 2010"},
		{"2003-06-30", true, "date should not be before december 2010"},
		{"2002-10-25", true, "date should not be before december 2010"},
		{"2012-07-32", true, "parsing time \"2012-07-32\": day out of range"},
		{"2023-06-31", true, "parsing time \"2023-06-31\": day out of range"},
		{"2024-02-30", true, "parsing time \"2024-02-30\": day out of range"},
		{"2022-11-31", true, "parsing time \"2022-11-31\": day out of range"},
		{"2023-09-31", true, "parsing time \"2023-09-31\": day out of range"},
		{"2024-01-32", true, "parsing time \"2024-01-32\": day out of range"},
		{"2022-05-32", true, "parsing time \"2022-05-32\": day out of range"},
		{"2023-03-32", true, "parsing time \"2023-03-32\": day out of range"},
		{"2024-08-32", true, "parsing time \"2024-08-32\": day out of range"},
		{"2022-12-32", true, "parsing time \"2022-12-32\": day out of range"},
		{"2023-07-32", true, "parsing time \"2023-07-32\": day out of range"},
		{"20241-15-82", true, "parsing time \"20241-15-82\" as \"2006-01-02\": cannot parse \"1-15-82\" as \"-\""},
		{"20239-04-67", true, "parsing time \"20239-04-67\" as \"2006-01-02\": cannot parse \"9-04-67\" as \"-\""},
		{"20242-29-53", true, "parsing time \"20242-29-53\" as \"2006-01-02\": cannot parse \"2-29-53\" as \"-\""},
		{"20240-11-18", true, "parsing time \"20240-11-18\" as \"2006-01-02\": cannot parse \"0-11-18\" as \"-\""},
		{"20243-23-74", true, "parsing time \"20243-23-74\" as \"2006-01-02\": cannot parse \"3-23-74\" as \"-\""},
		{"20238-06-29", true, "parsing time \"20238-06-29\" as \"2006-01-02\": cannot parse \"8-06-29\" as \"-\""},
		{"20244-19-85", true, "parsing time \"20244-19-85\" as \"2006-01-02\": cannot parse \"4-19-85\" as \"-\""},
		{"20240-12-41", true, "parsing time \"20240-12-41\" as \"2006-01-02\": cannot parse \"0-12-41\" as \"-\""},
		{"20245-05-90", true, "parsing time \"20245-05-90\" as \"2006-01-02\": cannot parse \"5-05-90\" as \"-\""},
		{"20237-08-34", true, "parsing time \"20237-08-34\" as \"2006-01-02\": cannot parse \"7-08-34\" as \"-\""},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test-%v", i), func(t *testing.T) {
			if tt.wantErr {
				log.Print("Checking for unexpected errors...")
				_, err := CheckDate(tt.expectedDate)
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("CheckDate() returned unexpected error: got \"%v\" want \"%v\"", err.Error(), tt.expectedError)
				}
				log.Print("No unexpected errors found")
			} else {
				date, err := CheckDate(tt.expectedDate)
				if err != nil {
					t.Error("CheckDate() returned an error: ", err)
				}
				_date := date.Format("2006-01-02")
				if tt.expectedDate != _date {
					t.Errorf("CheckDate() returned wrong date: got %v want %v", _date, tt.expectedDate)
				}
			}
		})
	}
}

func getDate(year, month, day int) string {
	return fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
}
