package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
)

var weatherAPIURL, geocodingAPIURL, geocodingAPIKey string

func GetWeatherAPIURL() string {
	return weatherAPIURL
}

func SetWeatherAPIURL(url string) {
	weatherAPIURL = url
}

func GetGeocodingAPIURL() string {
	return geocodingAPIURL
}

func SetGeocodingAPIURL(url string) {
	geocodingAPIURL = url
}

func SetGeocodingAPIKey(url string) {
	geocodingAPIKey = url
}

func RequestWeather() (model.WeatherRecord, error) {
	var today model.WeatherRecord
	resp, err := http.Get(weatherAPIURL)
	if err != nil {
		return model.WeatherRecord{}, fmt.Errorf("failed to get weatherAPI from URL (\"%v\"): %v", weatherAPIURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.WeatherRecord{}, fmt.Errorf("failed to read response from URL (\"%v\"): %v", weatherAPIURL, err)
	}
	err = json.Unmarshal(body, &today)
	if err != nil {
		return model.WeatherRecord{}, fmt.Errorf("failed to unmarshal response body from URL (\"%v\"): %v", weatherAPIURL, err)
	}
	hours := []model.HourWeatherRecord{}
	for k, v := range today.Hours {
		if k <= 23 {
			hours = append(hours, v)
		}
	}
	today.Hours = hours
	return today, nil
}

func ReloadWeatherURL(date string, latitude, longitude float64) error {
	u, err := url.Parse(weatherAPIURL)
	if err != nil {
		return fmt.Errorf("incorrect values in weatherAPI URL (%v): %v", weatherAPIURL, err)
	}

	v := u.Query()
	v.Set("lat", strconv.FormatFloat(latitude, 'f', 2, 64))
	v.Set("lon", strconv.FormatFloat(longitude, 'f', 2, 64))
	v.Set("date", date)
	u.RawQuery = v.Encode()
	weatherAPIURL = u.String()
	return nil
}

func ReloadGeoURL(cityName string) error {
	u, err := url.Parse(geocodingAPIURL)
	if err != nil {
		return fmt.Errorf("incorrect values in geocodingAPI URL (%v): %v", geocodingAPIURL, err)
	}

	v := u.Query()
	v.Set("q", cityName)
	v.Set("appid", geocodingAPIKey)
	u.RawQuery = v.Encode()
	geocodingAPIURL = u.Redacted()
	return nil
}

func RequestFutureWeather() (model.WeatherRecord, error) {
	var notToday model.WeatherRecord
	resp, err := http.Get(weatherAPIURL)
	if err != nil {
		return model.WeatherRecord{}, fmt.Errorf("failed to get weatherAPI from URL (\"%v\"): %v", weatherAPIURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.WeatherRecord{}, fmt.Errorf("failed to read response from URL (\"%v\"): %v", weatherAPIURL, err)
	}
	err = json.Unmarshal(body, &notToday)
	if err != nil {
		return model.WeatherRecord{}, fmt.Errorf("failed to unmarshal response body from URL (\"%v\"): %v", weatherAPIURL, err)
	}
	hours := []model.HourWeatherRecord{}
	for k, v := range notToday.Hours {
		if k <= 23 {
			hours = append(hours, v)
		}
	}
	notToday.Hours = hours
	return notToday, nil
}

func RequestWeatherEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		RequestWeather()
		f(x)
	}
}
