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
	ContainsCity(city string) (bool, error)
	GetCity(cityName string) (model.City, error)
}

type cities struct {
	db   DBI
	file string
}

func NewCityDB(database DBI, file string) CityDBI {
	if file == "" {
		file = "resources/data/cities.json"
	}
	return cities{
		db:   database,
		file: file,
	}
}

func (c cities) ContainsCity(cityName string) (bool, error) {
	cityName = strings.ToLower(cityName)
	var cities []model.City
	cities, err := c.db.GetCities()
	if err != nil {
		return false, err
	}
	for _, v := range cities {
		if v.Name == cityName {
			return true, nil
		}
	}
	return false, nil
}

func (c cities) GetCity(cityName string) (model.City, error) {
	cityName = strings.ToLower(cityName)
	var cities []model.City
	//err := c.db.GetDatabase().Model(&cities).Table("cities").Where("cities.name = ?", cityName).Select()
	cities, err := c.db.GetCities()
	if err != nil {
		return model.City{}, err
	}
	var city model.City
	for _, v := range cities {
		if v.Name == cityName {
			city = v
			continue
		}
	}
	return city, nil
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
