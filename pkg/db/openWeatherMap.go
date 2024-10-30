package db

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/utils"
)

type WeatherMapDB interface {
	SaveCityByName(name string) (string, error)
	GetRefactoredGeoURL(cityName string) (string, error)
}

type weatherMap struct {
	cities map[string]model.City
	url    string
	key    string
}

func NewWeatherMapDB(cities map[string]model.City, url string, key string) WeatherMapDB {
	return weatherMap{cities: cities, key: key}
}

func (wm weatherMap) SaveCityByName(name string) (string, error) {
	cities := wm.cities
	url, err := wm.GetRefactoredGeoURL(name)
	if err != nil {
		return "", err
	}
	var owcities []model.OWCity
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &owcities)
	if err != nil {
		return "", err
	}
	var foundcity model.OWCity
	for _, owcity := range owcities {
		if owcity.Country == "DE" {
			foundcity = owcity
		}
	}
	err = utils.SetLocation(foundcity.Latitude, foundcity.Longitude) // maybe change later?
	if err != nil {
		return "", err
	}
	cities[strings.ToLower(name)] = model.City{Name: foundcity.Name, Lat: foundcity.Latitude, Lon: foundcity.Longitude}
	data, err := json.MarshalIndent(cities, "", "  ")
	if err != nil {
		return "", err
	}
	resourcesPath := "resources/data"
	utils.SaveFile(resourcesPath, "cities.json", data)
	citiesPath := "resources/data/cities.txt"
	if !utils.PathExists(citiesPath) {
		var citiesData []byte
		for _, v := range []byte(strings.ToLower(name)) {
			citiesData = append(citiesData, v)
		}
		utils.SaveFile(resourcesPath, "cities.txt", citiesData)
	} else {
		file, err := os.Open(citiesPath)
		if err != nil {
			return "", err
		}
		defer file.Close()
		content, err := os.ReadFile(citiesPath)
		if err != nil {
			return "", err
		}
		var citiesData []byte
		if !strings.Contains(string(content), strings.ToLower(name)) {
			citiesString := string(content) + "\n" + strings.ToLower(name)
			for _, v := range []byte(citiesString) {
				citiesData = append(citiesData, v)
			}
			utils.SaveFile("resources/data", "cities.txt", citiesData)
		}
	}
	return name, nil
}

func (wm weatherMap) GetRefactoredGeoURL(cityName string) (string, error) {
	u, err := url.Parse(wm.url)
	if err != nil {
		return "", fmt.Errorf("incorrect values in geocodingAPI URL (%v): %v", wm.url, err)
	}

	v := u.Query()
	v.Set("q", cityName)
	v.Set("appid", wm.key)
	u.RawQuery = v.Encode()
	return u.Redacted(), nil
}
