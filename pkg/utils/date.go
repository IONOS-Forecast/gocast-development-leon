package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
)

var date string

func SplitDate(date string) (year, month, day int, err error) {
	splitDate := strings.Split(date, "-")
	day, err = strconv.Atoi(splitDate[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("ERROR: Couldn't convert day!\nERROR: %v", err)
	}
	month, err = strconv.Atoi(splitDate[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("ERROR: Couldn't convert month!\nERROR: %v", err)
	}
	year, err = strconv.Atoi(splitDate[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("ERROR: Couldn't convert year!\nERROR: %v", err)
	}
	return year, month, day, nil
}

func SetDate(year, month, day int) (string, error) {
	now := time.Now()
	_date := fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	_, err := CheckDate(_date)
	if err != nil { // What happens when date is invalid
		date = now.Format("2006-01-02")
		return date, fmt.Errorf("date (\"%v\") invalid! WARNING: Date set to today (\"%v\")", _date, date)
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
		return newDate, count, fmt.Errorf("converting counter failed: %v", err)
	}
	year, month, day, err := SplitDate(oldDate)
	if err != nil {
		return newDate, count, err
	}
	_count += 1
	day = day + _count
	count = strconv.Itoa(_count)
	_, err = SetDate(year, month, day)
	if err != nil {
		return "", "", err
	}
	return fmt.Sprintf("%v-%.2v-%.2v", year, month, day), count, nil
}

func CheckDate(s string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", s)
	if err != nil {
		return date, err
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
