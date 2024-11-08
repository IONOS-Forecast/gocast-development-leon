package utils

import (
	"fmt"
	"net/url"
	"testing"
)

func TestGetWeatherAPIURL(t *testing.T) {
	expectedURL := Options.WeatherAPIURL
	SetWeatherAPIURL(expectedURL)
	url := GetWeatherAPIURL()
	if url != expectedURL {
		t.Errorf("GetWeatherAPIURL() returned wrong url: got \"%v\" want \"%v\"", url, expectedURL)
	}
}

func TestGetGeocodingAPIURL(t *testing.T) {
	expectedURL := "TestURL"
	SetGeocodingAPIURL(expectedURL)
	url := GetGeocodingAPIURL()
	if url != expectedURL {
		t.Errorf("GetGeocodingAPIURL() returned wrong url: got \"%v\" want \"%v\"", url, expectedURL)
	}
}

func TestSetGeocodingAPIKey(t *testing.T) {
	expectedKey := "API_KEY"
	SetGeocodingAPIKey(expectedKey)
	if geocodingAPIKey != expectedKey {
		t.Errorf("ReloadWeatherURL() returned wrong key: got \"%v\" want \"%v\"", geocodingAPIKey, expectedKey)
	}
}

func TestSetWeatherAPIURL(t *testing.T) {
	expectedURL := "TestURL"
	SetWeatherAPIURL(expectedURL)
	if weatherAPIURL != expectedURL {
		t.Errorf("SetWeatherAPIURL() returned wrong url: got \"%v\" want \"%v\"", weatherAPIURL, expectedURL)
	}
}

func TestSetGeocodingAPIURL(t *testing.T) {
	expectedURL := "TestURL"
	SetGeocodingAPIURL(expectedURL)
	if geocodingAPIURL != expectedURL {
		t.Errorf("SetGeocodingAPIURL() returned wrong url: got \"%v\" want \"%v\"", geocodingAPIURL, expectedURL)
	}
}

func TestReloadWeatherURL(t *testing.T) {
	tests := []struct {
		date string
		lat  float64
		lon  float64
	}{
		{"2024-08-25", 52.52, 13.39},
		{"2018-05-01", 53.55, 10.00},
		{"2015-12-31", 48.14, 11.58},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test-%v", i), func(t *testing.T) {
			defaultURL := "https://api.com/weather?date=2010-12-31&lat=28&lon=5.8"
			SetWeatherAPIURL(defaultURL)
			expectedURL, err := reloadWeatherURL("2024-08-25", tt.lat, tt.lon)
			if err != nil {
				t.Error("ReloadWeatherURL() returned an error: ", err)
			}
			if GetWeatherAPIURL() != expectedURL {
				t.Errorf("ReloadWeatherURL() returned wrong url: got \"%v\" want \"%v\"", weatherAPIURL, expectedURL)
			}
		})
	}
}

func TestReloadGeocodingURL(t *testing.T) {
	tests := []struct {
		city string
	}{
		{"Berlin"},
		{"München"},
		{"Hamburg"},
		{"Köln"},
		{"Frankfurt am Main"},
		{"Stuttgart"},
		{"Düsseldorf"},
		{"Leipzig"},
		{"Dresden"},
		{"Nürnberg"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test-%v", i), func(t *testing.T) {
			defaultURL := "https://api.com/geo?appid=KEY&q=Brandenburg"
			SetGeocodingAPIURL(defaultURL)
			expectedURL, err := reloadGeoURL(tt.city)
			if err != nil {
				t.Error("ReloadWeatherURL() returned an error: ", err)
			}
			if GetGeocodingAPIURL() != expectedURL {
				t.Errorf("ReloadGeocodingURL() returned wrong url: got \"%v\" want \"%v\"", geocodingAPIURL, expectedURL)
			}
		})
	}
}

func reloadWeatherURL(date string, lat, lon float64) (string, error) {
	return fmt.Sprintf("https://api.com/weather?date=%v&lat=%.2f&lon=%.2f", date, lat, lon), ReloadWeatherURL(date, lat, lon)
}

func reloadGeoURL(city string) (string, error) {
	SetGeocodingAPIKey("API_KEY")
	return fmt.Sprintf("https://api.com/geo?appid=%v&q=%s", geocodingAPIKey, url.QueryEscape(city)), ReloadGeoURL(city)
}
