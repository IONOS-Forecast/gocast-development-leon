package utils

import (
	"fmt"
	"testing"
)

func TestGetWeatherAPIKeyAPIURL(t *testing.T) {
	expectedURL := Options.WeatherAPIURL
	SetWeatherAPIURL(expectedURL)
	url := GetWeatherAPIURL()
	if url != expectedURL {
		t.Errorf("GetWeatherAPIURL() returned wrong url: got \"%v\" want \"%v\"", url, expectedURL)
	}
}

func TestGetGeocodingAPIURL(t *testing.T) {
	expectedURL := Options.GeoAPIURL
	SetWeatherAPIURL(expectedURL)
	url := GetGeocodingAPIURL()
	if url != expectedURL {
		t.Errorf("GetGeocodingAPIURL() returned wrong url: got \"%v\" want \"%v\"", url, expectedURL)
	}
}

func TestSetWeatherAPIKeyAPIURL(t *testing.T) {
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
	url := "https://api.com/weather?date=2010-12-31&lat=28&lon=5.8"
	SetWeatherAPIURL(url)
	expectedURL, err := reloadWeatherURL("2024-08-25", 52.52, 13.39)
	if err != nil {
		t.Error("ReloadWeatherURL() returned an error: ", err)
	}
	if GetWeatherAPIURL() != expectedURL {
		t.Errorf("ReloadWeatherURL() returned wrong url: got \"%v\" want \"%v\"", weatherAPIURL, expectedURL)
	}

	expectedURL, err = reloadWeatherURL("2018-05-01", 53.55, 10.00)
	if err != nil {
		t.Error("ReloadWeatherURL() returned an error: ", err)
	}
	if GetWeatherAPIURL() != expectedURL {
		t.Errorf("ReloadWeatherURL() returned wrong url: got \"%v\" want \"%v\"", weatherAPIURL, expectedURL)
	}

	expectedURL, err = reloadWeatherURL("2015-12-31", 48.14, 11.58)
	if err != nil {
		t.Error("ReloadWeatherURL() returned an error: ", err)
	}
	if GetWeatherAPIURL() != expectedURL {
		t.Errorf("ReloadWeatherURL() returned wrong url: got \"%v\" want \"%v\"", weatherAPIURL, expectedURL)
	}
}

func TestReloadGeocodingURL(t *testing.T) {
	url := "https://api.com/geo?appid=KEY&q=Brandenburg"
	SetGeocodingAPIURL(url)
	expectedURL, err := reloadGeoURL("Berlin")
	if err != nil {
		t.Error("ReloadWeatherURL() returned an error: ", err)
	}
	if GetGeocodingAPIURL() != expectedURL {
		t.Errorf("ReloadGeocodingURL() returned wrong url: got \"%v\" want \"%v\"", geocodingAPIURL, expectedURL)
	}
	expectedURL, err = reloadGeoURL("Hamburg")
	if err != nil {
		t.Error("ReloadWeatherURL() returned an error: ", err)
	}
	if GetGeocodingAPIURL() != expectedURL {
		t.Errorf("ReloadGeocodingURL() returned wrong url: got \"%v\" want \"%v\"", geocodingAPIURL, expectedURL)
	}
	// TODO: München doesnt work because of special character 'ü'
	expectedURL, err = reloadGeoURL("München")
	if err != nil {
		t.Error("ReloadWeatherURL() returned an error: ", err)
	}
	if GetGeocodingAPIURL() != expectedURL {
		t.Errorf("ReloadGeocodingURL() returned wrong url: got \"%v\" want \"%v\"", geocodingAPIURL, expectedURL)
	}
}

func reloadWeatherURL(date string, lat, lon float64) (string, error) {
	return fmt.Sprintf("https://api.com/weather?date=%v&lat=%.2f&lon=%.2f", date, lat, lon), ReloadWeatherURL(date, lat, lon)
}

func reloadGeoURL(city string) (string, error) {
	SetGeocodingAPIKey("API_KEY")
	return fmt.Sprintf("https://api.com/geo?appid=%v&q=%v", geocodingAPIKey, city), ReloadGeoURL(city)
}
