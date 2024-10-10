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
		t.Error("SetDate() in test GetDate() returned an error: ", err)
	}
	if expectedDate != GetDate() {
		t.Errorf("GetDate() returned wrong date: got %v want %v", date, expectedDate)
	}
}

func TestSetFutureDate(t *testing.T) {
	// Expecting results
	year := 2011
	month := 2
	day := 27
	testDate := fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	expectedDate := fmt.Sprintf("%v-%.2v-%.2v", year, month, day+1)
	date, _, err := SetFutureDay(expectedDate, testDate, "0")
	if err != nil {
		t.Error("SetFutureDate() returned an error: ", err)
	}
	if expectedDate != date {
		t.Errorf("SetFutureDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	year = 2015
	month = 8
	day = 1
	testDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day+1)
	date, _, err = SetFutureDay(expectedDate, testDate, "0")
	if err != nil {
		t.Error("SetFutureDate() returned an error: ", err)
	}
	if expectedDate != date {
		t.Errorf("SetFutureDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	// Expecting errors
	log.Print("INFO: 4 expected warnings following!")
	year = 2011
	month = 2
	day = 28
	testDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day+1)
	expectedError := fmt.Errorf("parsing time \"%v\": day out of range", expectedDate).Error()
	date, _, err = SetFutureDay(expectedDate, testDate, "0")
	if err == nil || err.Error() != expectedError {
		t.Errorf("SetFutureDate() returned wrong error: got \"%v\" want \"%v\"", err, expectedError)
	}
	if expectedDate == date {
		t.Errorf("SetFutureDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	year = 2010
	month = 12
	day = 31
	testDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day+1)
	expectedError = fmt.Errorf("parsing time \"%v\": day out of range", expectedDate).Error()
	date, _, err = SetFutureDay(expectedDate, testDate, "0")
	if err == nil || err.Error() != expectedError {
		t.Errorf("SetFutureDate() returned wrong error: got \"%v\" want \"%v\"", err, expectedError)
	}
	if expectedDate == date {
		t.Errorf("SetFutureDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	year = 2008
	month = 8
	day = 25
	testDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day+1)
	expectedError = fmt.Errorf("date should not be before december 2010").Error()
	date, _, err = SetFutureDay(expectedDate, testDate, "0")
	if err == nil || err.Error() != expectedError {
		t.Errorf("SetFutureDate() returned wrong error: got \"%v\" want \"%v\"", err, expectedError)
	}
	if expectedDate == date {
		t.Errorf("SetFutureDate() returned wrong date: got %v want %v", date, expectedDate)
	}
	now := time.Now()
	year = now.Year()
	month = int(now.Month())
	day = now.Day()
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
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", 0, 0, 0)

	yearString := "2018"
	monthString := "5"
	dayString := "five"
	date = fmt.Sprintf("%v-%.2v-%.2v", yearString, monthString, dayString)
	expectedError := fmt.Errorf("failed to convert day: strconv.Atoi: parsing \"%.2v\": invalid syntax", dayString).Error()
	year, month, day, err = SplitDate(date)
	if err == nil || err.Error() != expectedError {
		t.Errorf("SplitDate() returned unexpected error: got \"%v\" want \"%v\"", err, expectedError)
	}
	date = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	if date != expectedDate {
		t.Errorf("SplitDate() returned unexpected date: got %v want %v", date, expectedDate)
	}
	yearString = "2018"
	monthString = "five"
	dayString = "2"
	date = fmt.Sprintf("%v-%.2v-%.2v", yearString, monthString, dayString)
	expectedError = fmt.Errorf("failed to convert month: strconv.Atoi: parsing \"%.2v\": invalid syntax", monthString).Error()
	year, month, day, err = SplitDate(date)
	if err == nil || err.Error() != expectedError {
		t.Errorf("SplitDate() returned unexpected error: got \"%v\" want \"%v\"", err, expectedError)
	}
	date = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	if date != expectedDate {
		t.Errorf("SplitDate() returned unexpected date: got %v want %v", date, expectedDate)
	}
	yearString = "five"
	monthString = "5"
	dayString = "2"
	date = fmt.Sprintf("%v-%.2v-%.2v", yearString, monthString, dayString)
	expectedError = fmt.Errorf("failed to convert year: strconv.Atoi: parsing \"%v\": invalid syntax", yearString).Error()
	year, month, day, err = SplitDate(date)
	if err == nil || err.Error() != expectedError {
		t.Errorf("SplitDate() returned unexpected error: got \"%v\" want \"%v\"", err, expectedError)
	}
	date = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	if date != expectedDate {
		t.Errorf("SplitDate() returned unexpected date: got %v want %v", date, expectedDate)
	}
}

func TestCheckDate(t *testing.T) {
	year := 2010
	month := 12
	day := 5
	expectedDate := fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	date, err := CheckDate(expectedDate)
	if err != nil {
		t.Error("CheckDate() returned an error: ", err)
	}
	_date := date.Format("2006-01-02")
	if expectedDate != _date {
		t.Errorf("CheckDate() returned wrong date: got %v want %v", _date, expectedDate)
	}
	year = 2016
	month = 3
	day = 6
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	date, err = CheckDate(expectedDate)
	if err != nil {
		t.Error("CheckDate() returned an error: ", err)
	}
	_date = date.Format("2006-01-02")
	if expectedDate != _date {
		t.Errorf("CheckDate() returned wrong date: got %v want %v", _date, expectedDate)
	}
	now := time.Now()
	year = now.Year()
	month = int(now.Month())
	day = now.Day()
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	date, err = CheckDate(expectedDate)
	if err != nil {
		t.Error("CheckDate() returned an error: ", err)
	}
	_date = date.Format("2006-01-02")
	if expectedDate != _date {
		t.Errorf("CheckDate() returned wrong date: got %v want %v", _date, expectedDate)
	}
	// Expecting errors
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day+8)
	expectedError := fmt.Errorf("date should not be after next 7 days").Error()
	_, err = CheckDate(expectedDate)
	if err == nil || err.Error() != expectedError {
		t.Errorf("CheckDate() returned unexpected error: got \"%v\" want \"%v\"", err.Error(), expectedError)
	}
	year = 2008
	month = 8
	day = 25
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	expectedError = fmt.Errorf("date should not be before december 2010").Error()
	_, err = CheckDate(expectedDate)
	if err == nil || err.Error() != expectedError {
		t.Errorf("CheckDate() returned unexpected error: got \"%v\" want \"%v\"", err.Error(), expectedError)
	}
	year = 2012
	month = 7
	day = 32
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	expectedError = fmt.Errorf("parsing time \"%v\": day out of range", expectedDate).Error()
	_, err = CheckDate(expectedDate)
	if err == nil || err.Error() != expectedError {
		t.Errorf("CheckDate() returned unexpected error: got \"%v\" want \"%v\"", err.Error(), expectedError)
	}
	year = 20240
	month = 7
	day = 32
	expectedDate = fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	expectedError = fmt.Errorf("parsing time \"%v\" as \"2006-01-02\": cannot parse \"%v\" as \"-\"", expectedDate, expectedDate[4:]).Error()
	_, err = CheckDate(expectedDate)
	if err == nil || err.Error() != expectedError {
		t.Errorf("CheckDate() returned unexpected error: got \"%v\" want \"%v\"", err.Error(), expectedError)
	}
}
