package metric_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/db"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/metric"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
)

/*
func TestGet(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/direct?date=2024-10-08&city=Berlin", nil)
	rr := httptest.NewRecorder()
	var h metric.Handler
	h.Get(rr, req)
	//handler := http.HandlerFunc(h.Get)
	//handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}*/

type dbMock struct {
	weatherRecord model.WeatherRecord
	err           error
	exist         bool
	city          string
}

func (m *dbMock) setResult(weatherRecord model.WeatherRecord, err error, exist bool, city string) {
	m.weatherRecord = weatherRecord
	m.err = err
	m.exist = exist
	m.city = city
}

func (m dbMock) QueryDayDatabase(city string, date string) ([]model.HourWeatherRecord, error) {
	return m.weatherRecord.Hours, m.err
}
func (m dbMock) QueryDatabase(t any, value string, date string, hour int, city string) error {
	return m.err
}
func (m dbMock) WeatherDataExists(city string, date string) (bool, error) {
	return m.exist, m.err
}
func (m dbMock) GetWeatherRecord(city string, date string) (model.WeatherRecord, error) {
	return m.weatherRecord, m.err
}
func (m dbMock) InsertCityIntoDatabase(name string) error {
	return m.err
}
func (m dbMock) InsertCityWeatherRecordsToTable(record model.WeatherRecord) error {
	return m.err
}
func (m dbMock) QueryCitiesDatabase(t any, value string, name string) error {
	return m.err
}
func (m dbMock) SetLocationByCityName(city string) (string, error) {
	return m.city, m.err
}

func TestGetError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/direct?date=2024-10-20&city=Berlin", nil)
	rr := httptest.NewRecorder()
	dbMock := dbMock{}
	dbMock.setResult(model.WeatherRecord{}, nil, true, "Berlin")
	h := metric.NewHandler(dbMock, db.NewCityDB(""))
	h.Get(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetOKOnEmptyResult(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/direct?date=2024-09-20&city=Berlin", nil)
	rr := httptest.NewRecorder()
	dbMock := dbMock{}
	dbMock.setResult(model.WeatherRecord{}, nil, true, "Berlin")
	h := metric.NewHandler(dbMock, db.NewCityDB(""))
	h.Get(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

//func TestError(t *testing.T) {
//	req := httptest.NewRequest("GET", "http://localhost:8080/error", nil)
//	rr := httptest.NewRecorder()
//	var h metric.Handler
//	h.Get(rr, req)
//	if status := rr.Code; status != http.StatusBadRequest {
//		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
//	}
//}
