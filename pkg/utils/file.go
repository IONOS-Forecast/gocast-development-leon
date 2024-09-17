package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
)

func PathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func SaveFile(directory, filename string, data []byte) error {
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	filePath := fmt.Sprintf("%s/%s", directory, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func SaveFutureWeatherInFile(city string, date string) ([]model.WeatherRecord, error) {
	count := "0"
	var records []model.WeatherRecord
	record, err := RequestWeather()
	if err != nil {
		return records, err
	}
	err = SaveWeather(city, count, record)
	if err != nil {
		return records, err
	}
	oldDate := date
	newDate := date
	for i := 1; i <= 7; i++ { // Create for the next 6 days after the first
		newDate, count, err = SetFutureDay(newDate, oldDate, count)
		if err != nil {
			return []model.WeatherRecord{}, err
		}
		record, err = RequestFutureWeather()
		if err != nil {
			return []model.WeatherRecord{}, err
		}
		records = append(records, record)
		err = SaveWeather(city, count, record)
		if err != nil {
			return []model.WeatherRecord{}, err
		}
	}
	year, month, day, err := SplitDate(date)
	if err != nil {
		return []model.WeatherRecord{}, err
	}
	day += 1
	_, err = SetDate(year, month, day)
	if err != nil {
		return []model.WeatherRecord{}, err
	}
	return records, nil
}

func SaveWeather(city string, count string, record model.WeatherRecord) error {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("ERROR: MarshalIndent threw an error!\nERROR: %v", err)
	}

	err = SaveFile("resources/weather_records", strings.ToLower(city)+"_"+count+"-orig.json", data)
	if err != nil {
		return err
	}
	return nil
}

func GetAllWeatherRecords(city string, date string) ([]model.WeatherRecord, error) {
	count := "0"
	var records []model.WeatherRecord
	record, err := RequestWeather()
	if err != nil {
		return []model.WeatherRecord{}, err
	}
	err = SaveWeather(city, count, record)
	if err != nil {
		return []model.WeatherRecord{}, err
	}
	records = append(records, record)
	oldDate := date
	newDate := date
	for i := 1; i <= 7; i++ { // Create for the next 6 days after the first
		newDate, count, err = SetFutureDay(newDate, oldDate, count)
		if err != nil {
			return []model.WeatherRecord{}, err
		}
		record, err = RequestFutureWeather()
		if err != nil {
			return []model.WeatherRecord{}, err
		}
		records = append(records, record)
	}
	year, month, day, err := SplitDate(date)
	if err != nil {
		return []model.WeatherRecord{}, err
	}
	day += 1
	_, err = SetDate(year, month, day)
	if err != nil {
		return []model.WeatherRecord{}, err
	}
	return records, nil
}

func ConvertWeatherRecordss() {
	path := "scripts/convert.sh"
	exec.Command("/bin/bash", path)
}
