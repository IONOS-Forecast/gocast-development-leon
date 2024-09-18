package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
)

var latitude, longitude float64
var cityName string
var cities map[string]model.City

func MakeCities() {
	cities = make(map[string]model.City)
}

func GetCityName() string {
	return cityName
}

func GetCities() map[string]model.City {
	cities = ReadSavedCities(cities)
	return cities
}

func SetCities(map[string]model.City) map[string]model.City {
	return cities
}

func GetLocation() (lat, lon float64, err error) {
	return latitude, longitude, nil
}

func SetLocation(lat, lon float64) error {
	if lat <= 54 && lat >= 48 && lon <= 14 && lon >= 6 {
		latitude = lat
		longitude = lon
		err := ReloadWeatherURL(GetDate(), latitude, longitude)
		if err != nil {
			return err
		}
		log.Printf("INFO: Location set to (Lat: \"%v\"Lon: \"%v\")", lat, lon)
		return nil
	} else { // When location is not in range
		return fmt.Errorf("location (Lat: \"%v\" Lon: \"%v\") is not in range!", lat, lon)
	}
}

func SetLocationByCityName(name string) (string, error) {
	cities, err := ReadCities(name, GetCities())
	if err != nil {
		return "", err
	}
	if city, exists := cities[strings.ToLower(name)]; exists { // When the city exists
		cityName = name
		err := SetLocation(city.Lat, city.Lon)
		if err != nil {
			return "", err
		}
		return cityName, nil
	} else { // When the city doesn't exist
		log.Printf("INFO: City (\"%v\") doesn't exist!", strings.ToLower(name))
		log.Printf("INFO: Getting city (\"%v\") from API!", strings.ToLower(name))
		cityName, err := SaveCityByName(name)
		if err != nil {
			return "", err
		}
		cityName, err = SetLocationByCityName(cityName)
		if err != nil {
			return "", err
		}
		log.Printf("INFO: Location set to (Lat: \"%v\"Lon: \"%v\")", city.Lat, city.Lon)
		return cityName, err
	}
}

func ReadCities(name string, cities map[string]model.City) (map[string]model.City, error) {
	file, err := os.Open("resources/data/cities.json")
	if err != nil {
		_, err = SaveCityByName(name)
		if err != nil {
			return map[string]model.City{}, err
		}
		return ReadCities(name, cities)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cities)
	if err != nil {
		return map[string]model.City{}, err
	}
	return cities, nil
}

func ReadSavedCities(cities map[string]model.City) map[string]model.City {
	file, err := os.Open("resources/data/cities.json")
	if err != nil {
		return ReadSavedCities(cities)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cities)
	if err != nil {
		panic(err)
	}
	return cities
}

func SaveCityByName(name string) (string, error) {
	oldCityName := cityName
	cityName = name
	err := ReloadGeoURL(cityName)
	if err != nil {
		return "", err
	}
	var owcities []model.OWCity
	resp, err := http.Get(GetGeocodingAPIURL())
	if err != nil {
		cityName = oldCityName
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		cityName = oldCityName
		return "", err
	}
	err = json.Unmarshal(body, &owcities)
	if err != nil {
		cityName = oldCityName
		return "", err
	}
	var foundcity model.OWCity
	for _, owcity := range owcities {
		if owcity.Country == "DE" {
			foundcity = owcity
		}
	}
	err = SetLocation(foundcity.Latitude, foundcity.Longitude)
	if err != nil {
		cityName = oldCityName
		return "", err
	}
	cities[strings.ToLower(name)] = model.City{Lat: foundcity.Latitude, Lon: foundcity.Longitude}
	data, err := json.MarshalIndent(cities, "", "  ")
	if err != nil {
		cityName = oldCityName
		return "", err
	}
	resourcesPath := "resources/data"
	SaveFile(resourcesPath, "cities.json", data)
	citiesPath := "resources/data/cities.txt"
	if !PathExists(citiesPath) {
		var citiesData []byte
		for _, v := range []byte(strings.ToLower(name)) {
			citiesData = append(citiesData, v)
		}
		SaveFile(resourcesPath, "cities.txt", citiesData)
	} else {
		file, err := os.Open(citiesPath)
		if err != nil {
			cityName = oldCityName
			return "", err
		}
		defer file.Close()
		content, err := os.ReadFile(citiesPath)
		if err != nil {
			cityName = oldCityName
			return "", err
		}
		var citiesData []byte
		if !strings.Contains(string(content), strings.ToLower(name)) {
			citiesString := string(content) + "\n" + strings.ToLower(name)
			for _, v := range []byte(citiesString) {
				citiesData = append(citiesData, v)
			}
			SaveFile("resources/data", "cities.txt", citiesData)
		}
	}
	return cityName, nil
}
