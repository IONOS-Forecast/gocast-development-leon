package db

import (
	"encoding/json"
	"os"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
)

type CityDBI interface {
}

type cities struct {
	file string
}

func NewCityDB(file string) CityDBI {
	if file == "" {
		file = ""
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
