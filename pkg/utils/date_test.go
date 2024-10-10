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
