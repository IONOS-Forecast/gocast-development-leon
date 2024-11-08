package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
)

var date string
var earliestDate = time.Date(2010, time.December, 1, 0, 0, 0, 0, time.UTC)

func SplitDate(date string) (year, month, day int, err error) {
	splitDate := strings.Split(date, "-")
	day, err = strconv.Atoi(splitDate[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to convert day: %v", err)
	}
	month, err = strconv.Atoi(splitDate[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to convert month: %v", err)
	}
	year, err = strconv.Atoi(splitDate[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to convert year: %v", err)
	}
	return year, month, day, nil
}

func SetDate(year, month, day int) (string, error) {
	now := time.Now()
	_date := fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	_, err := CheckDate(_date)
	if err != nil {
		date = now.Format("2006-01-02")
		log.Printf("WARNING: Date set to today (\"%v\")", date)
		return date, fmt.Errorf("%v", err)
	}
	if year >= 2010 && year <= now.Year() && month >= 1 && month <= 12 && day >= 1 && day <= 31 {
		date = _date
	} else { // What happens when the date is out of range
		date = now.Format("2006-01-02")
		return date, fmt.Errorf("date out of range (\"%v\") invalid! WARNING: Date set to today (\"%v\")", _date, date)
	}
	lat, lon, err := GetLocation()
	if err != nil { // What happens when location is invalid
		date = now.Format("2006-01-02")
		_, err := SetLocationByCityName(GetCityName())
		newLat, newLon, err := GetLocation()
		if err != nil {
			return date, err
		}
		return date, fmt.Errorf("location (Lat: \"%v\" Lon: \"%v\") is invalid! WARNING: Location set to Berlin (Lat: \"%v\" Lon: \"%v\")", lat, lon, newLat, newLon)
	}
	err = ReloadWeatherURL(date, lat, lon)
	if err != nil {
		return "", err
	}
	return date, nil
}

func SetFutureDay(newDate, oldDate, count string) (string, string, error) {
	_count, err := strconv.Atoi(count)
	if err != nil {
		return oldDate, count, fmt.Errorf("converting counter failed: %v", err)
	}
	year, month, day, err := SplitDate(oldDate)
	if err != nil {
		return oldDate, count, err
	}
	_count += 1
	day = day + _count
	count = strconv.Itoa(_count)
	date, err := SetDate(year, month, day)
	if err != nil {
		return date, "", err
	}
	return date, count, nil
}

func CheckDate(s string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", s)
	if err != nil {
		return date, err
	}
	var lastDate = time.Now().Add(7 * 24 * time.Hour)
	if date.Before(time.Date(2010, time.December, 1, 0, 0, 0, 0, time.UTC)) {
		return date, fmt.Errorf("date should not be before december 2010")
	}
	if date.After(lastDate) {
		return date, fmt.Errorf("date should not be after next 7 days")
	}
	return date, nil
}

func GetDate() string {
	return date
}

func SetDateAndLocationByCityName(year, month, day int, cityName string, cities map[string]model.City) (string, error) {
	_, err := SetDate(year, month, day)
	if err != nil {
		return "", err
	}
	cityName, err = SetLocationByCityName(cityName)
	if err != nil {
		return "", err
	}
	return cityName, nil
}

func SetDateAndLocation(year, month, day int, lat, lon float64) error {
	_, err := SetDate(year, month, day)
	if err != nil {
		return err
	}
	err = SetLocation(lat, lon)
	if err != nil {
		return err
	}
	return nil
}
