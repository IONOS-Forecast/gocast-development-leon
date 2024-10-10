package metric_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/db"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/metric"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/utils"
)

type dbMock struct {
	weatherRecord model.WeatherRecord
	err           error
	exists        bool
	city          string
	cities        map[string]model.City
}

func (m *dbMock) setResult(weatherRecord model.WeatherRecord, err error, exists bool, city string, cities map[string]model.City) {
	m.weatherRecord = weatherRecord
	m.err = err
	m.exists = exists
	m.city = city
	m.cities = cities
}

func (m dbMock) QueryDayDatabase(city, date string) ([]model.HourWeatherRecord, error) {
	return m.weatherRecord.Hours, m.err
}
func (m dbMock) QueryDatabase(t any, value string, date string, hour int, city string) error {
	return m.err
}
func (m dbMock) WeatherDataExists(city, date string) (bool, error) {
	return m.exists, m.err
}
func (m dbMock) GetWeatherRecord(city, date string) (model.WeatherRecord, error) {
	return m.weatherRecord, m.err
}
func (m dbMock) GetWeatherRecords(city, date string) (model.WeatherRecord, error) {
	return m.weatherRecord, m.err
}
func (m dbMock) InsertCityIntoDatabase(name string) error {
	return m.err
}
func (m dbMock) InsertCityWeatherRecordsToTable(record model.WeatherRecord) error {
	return m.err
}
func (m dbMock) QueryCitiesDatabase(t any, value, name string) error {
	return m.err
}
func (m dbMock) GetCities() ([]model.City, error) {
	var cities []model.City
	for _, value := range m.cities {
		cities = append(cities, value)
	}
	return cities, nil
}
func (m dbMock) SetLocationByCityName(city string) (string, error) {
	return m.city, m.err
}

func TestGetOK(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/direct?date=2024-10-09&city=Berlin", nil)
	rr := httptest.NewRecorder()
	dbMock := dbMock{}
	cities := map[string]model.City{"Berlin": {Name: "berlin", Lat: 52.5170365, Lon: 13.3888599}}
	dbMock.setResult(model.WeatherRecord{}, nil, true, "Berlin", cities)
	cityDB := db.NewCityDB(dbMock, "")
	h := metric.NewHandler(dbMock, cityDB, db.NewWeatherMapDB(cities, utils.Options.GeoAPIURL, utils.Options.GeoAPIKey))
	h.Get(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetOKNoData(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/direct?date=2024-10-21&city=Berlin", nil)
	rr := httptest.NewRecorder()
	dbMock := dbMock{}
	cities := map[string]model.City{"Berlin": {Name: "berlin", Lat: 52.5170365, Lon: 13.3888599}}
	dbMock.setResult(model.WeatherRecord{}, nil, true, "Berlin", cities)
	cityDB := db.NewCityDB(dbMock, "")
	h := metric.NewHandler(dbMock, cityDB, db.NewWeatherMapDB(cities, utils.Options.GeoAPIURL, utils.Options.GeoAPIKey))
	h.Get(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetErrorNoCity(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/direct?date=2024-10-21&city=Berlin", nil)
	rr := httptest.NewRecorder()
	dbMock := dbMock{}
	cities := map[string]model.City{}
	dbMock.setResult(model.WeatherRecord{}, nil, true, "Berlin", cities)
	cityDB := db.NewCityDB(dbMock, "")
	h := metric.NewHandler(dbMock, cityDB, db.NewWeatherMapDB(cities, utils.Options.GeoAPIURL, utils.Options.GeoAPIKey))
	h.Get(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetErrorIncorrectDate(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/direct?date=2008-08-25&city=Berlin", nil)
	rr := httptest.NewRecorder()
	dbMock := dbMock{}
	cities := map[string]model.City{"Berlin": {Name: "berlin", Lat: 52.5170365, Lon: 13.3888599}}
	dbMock.setResult(model.WeatherRecord{}, nil, true, "Berlin", cities)
	cityDB := db.NewCityDB(dbMock, "")
	h := metric.NewHandler(dbMock, cityDB, db.NewWeatherMapDB(cities, utils.Options.GeoAPIURL, utils.Options.GeoAPIKey))
	h.Get(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
