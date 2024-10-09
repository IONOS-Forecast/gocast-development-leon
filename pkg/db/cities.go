package db

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/utils"
)

type CityDBI interface {
	GetCities() (map[string]model.City, error)
	CityExists(name string) bool
	SetLocationByCityName(name string) (string, error)
}

type cities struct {
	file string
}

func NewCityDB(file string) CityDBI {
	if file == "" {
		file = "resources/data/cities.json"
	}
	return cities{
		file: file,
	}
}

func (c cities) GetCities() (map[string]model.City, error) {
	file, err := os.Open(c.file)
	if err != nil {
		return map[string]model.City{}, err
	}
	defer file.Close()

	var cities map[string]model.City
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cities)
	if err != nil {
		return map[string]model.City{}, err
	}
	return cities, nil
}

func (c cities) CityExists(name string) bool {
	cities, err := c.GetCities()
	if err != nil {
		return false
	}
	cities, err = utils.ReadCities(name, cities)
	if err != nil {
		return false
	}
	if city, exists := cities[strings.ToLower(name)]; exists { // When the city exists
		name = city.Name
		err := utils.SetLocation(city.Lat, city.Lon)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func (c cities) SetLocationByCityName(name string) (string, error) {
	cities, err := c.GetCities()
	if err != nil {
		return "", err
	}
	cities, err = utils.ReadCities(name, cities)
	if err != nil {
		return "", err
	}
	if city, exists := cities[strings.ToLower(name)]; exists { // When the city exists
		name = city.Name
		err := utils.SetLocation(city.Lat, city.Lon)
		if err != nil {
			return "", err
		}
		return name, nil
	}
	return "", fmt.Errorf("couldn't set locaton because city doesn't exist!")
}
