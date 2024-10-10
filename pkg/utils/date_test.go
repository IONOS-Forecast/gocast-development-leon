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
	expectedDate := getDate(2010, 12, 5)
	date, err := CheckDate(expectedDate)
	if err != nil {
		t.Error("CheckDate() returned an error: ", err)
	}
	_date := date.Format("2006-01-02")
	if expectedDate != _date {
		t.Errorf("CheckDate() returned wrong date: got %v want %v", _date, expectedDate)
	}
	expectedDate = getDate(2016, 3, 6)
	date, err = CheckDate(expectedDate)
	if err != nil {
		t.Error("CheckDate() returned an error: ", err)
	}
	_date = date.Format("2006-01-02")
	if expectedDate != _date {
		t.Errorf("CheckDate() returned wrong date: got %v want %v", _date, expectedDate)
	}
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	day := now.Day()
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	date, err = CheckDate(expectedDate)
	if err != nil {
		t.Error("CheckDate() returned an error: ", err)
	}
	_date = date.Format("2006-01-02")
	if expectedDate != _date {
		t.Errorf("CheckDate() returned wrong date: got %v want %v", _date, expectedDate)
	}
	t.Run("checking for unexpected errors", func(t *testing.T) {
		expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day+8)
		expectedError := fmt.Errorf("date should not be after next 7 days").Error()
		_, err = CheckDate(expectedDate)
		if err == nil || err.Error() != expectedError {
			t.Errorf("CheckDate() returned unexpected error: got \"%v\" want \"%v\"", err.Error(), expectedError)
		}
		expectedDate = getDate(2008, 8, 25)
		expectedError = fmt.Errorf("date should not be before december 2010").Error()
		_, err = CheckDate(expectedDate)
		if err == nil || err.Error() != expectedError {
			t.Errorf("CheckDate() returned unexpected error: got \"%v\" want \"%v\"", err.Error(), expectedError)
		}
		expectedDate = getDate(2012, 7, 32)
		expectedError = fmt.Errorf("parsing time \"%v\": day out of range", expectedDate).Error()
		_, err = CheckDate(expectedDate)
		if err == nil || err.Error() != expectedError {
			t.Errorf("CheckDate() returned unexpected error: got \"%v\" want \"%v\"", err.Error(), expectedError)
		}
		expectedDate = getDate(20240, 7, 32)
		expectedError = fmt.Errorf("parsing time \"%v\" as \"2006-01-02\": cannot parse \"%v\" as \"-\"", expectedDate, expectedDate[4:]).Error()
		_, err = CheckDate(expectedDate)
		if err == nil || err.Error() != expectedError {
			t.Errorf("CheckDate() returned unexpected error: got \"%v\" want \"%v\"", err.Error(), expectedError)
		}
	})
}

func getDate(year, month, day int) string {
	return fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
}
