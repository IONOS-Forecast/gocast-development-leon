package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
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
		return fmt.Errorf("failed to make all directories (MkdirAll): %v", err)
	}
	filePath := fmt.Sprintf("%s/%s", directory, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create: %v", err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write: %v", err)
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
	log.Printf("INFO: Weather for %v saved", strings.ToLower(city))
	return records, nil
}

func SaveWeather(city string, count string, record model.WeatherRecord) error {
	for i := 0; i < len(record.Hours); i++ {
		record.Hours[i].City = city
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to MarshalIndent: %v", err)
	}

	err = SaveFile("resources/weather_records", strings.ToLower(city)+"_"+count+"-orig.json", data)
	if err != nil {
		return err
	}
	return nil
}

func ConvertWeatherRecordss() {
	path := "scripts/convert.sh"
	exec.Command("/bin/bash", path)
	log.Printf("INFO: Weather Records successfully converted")
}

func GetWeatherRecordsFromFiles(city string) ([]model.WeatherRecord, error) {
	var records []model.WeatherRecord
	count := 0
	for i := 0; i <= 7; i++ {
		var record model.WeatherRecord
		path := fmt.Sprintf("resources/weather_records/%v_%v-orig.json", strings.ToLower(city), strconv.Itoa(count))
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(content, &record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
		count++
	}
	log.Printf("INFO: Weather Records received from File")
	return records, nil
}
