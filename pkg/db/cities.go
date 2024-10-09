package db

import (
	"encoding/json"
	"os"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
)

type CityDB interface {
	GetCities() map[string]model.City
}

type cities struct {
	file string
}

func NewCityDB(file string) CityDB {
	if file == "" {
		file = "resources/data/cities.json"
	}
	return cities{
		file: file,
	}
}

func (c cities) GetCities() map[string]model.City {
	file, err := os.Open(c.file)
	if err != nil {
		return map[string]model.City{}
	}
	defer file.Close()

	var cities map[string]model.City
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cities)
	if err != nil {
		return map[string]model.City{}
	}
	return cities
}
