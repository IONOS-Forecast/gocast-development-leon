package db

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/utils"
	"github.com/go-pg/pg/v10"
)

type DBI interface {
	QueryDayDatabase(city, date string) ([]model.HourWeatherRecord, error)
	QueryDatabase(t any, value string, date string, hour int, city string) error
	WeatherDataExists(city, date string) (bool, error)
	GetWeatherRecord(city, date string) (model.WeatherRecord, error)
	GetWeatherRecords(city, date string) (model.WeatherRecord, error)
	InsertCityIntoDatabase(name string) error
	InsertCityWeatherRecordsToTable(record model.WeatherRecord) error
	QueryCitiesDatabase(t any, value, name string) error
	GetDatabase() pg.DB
	SetLocationByCityName(city string) (string, error)
}

type postgresDB struct {
	pgDB pg.DB
}

func NewPG(user, password, name, address string) (DBI, error) {
	pgDB := connectToDatabase(user, password, name, address)
	return postgresDB{pgDB: pgDB}, nil
}

func connectToDatabase(user, password, database, address string) pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     address,
		User:     user,
		Password: password,
		Database: database,
	})

	return *db
}

func (p postgresDB) QueryDayDatabase(city string, date string) ([]model.HourWeatherRecord, error) {
	var res []model.HourWeatherRecord
	year, month, day, err := utils.SplitDate(date)
	if err != nil {
		return []model.HourWeatherRecord{}, fmt.Errorf("failed to query Database (date): %v", err)
	}
	if city == "" {
		return []model.HourWeatherRecord{}, fmt.Errorf("failed to query Database because city isn't set!")
	}
	query := fmt.Sprintf("timestamp::date='%v-%.2v-%.2v 00:00:00+00' AND city='%v'", year, month, day, city)
	err = p.pgDB.Model().Table("weather_records").
		Column("timestamp", "source_id", "precipitation", "pressure_msl", "sunshine", "temperature",
			"wind_direction", "wind_speed", "cloud_cover", "dew_point", "relative_humidity", "visibility",
			"wind_gust_direction", "wind_gust_speed", "condition", "precipitation_probability",
			"precipitation_probability_6h", "solar", "icon").
		Where(query).
		Select(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p postgresDB) QueryDatabase(t any, value string, date string, hour int, city string) error {
	year, month, day, err := utils.SplitDate(date)
	if err != nil {
		return fmt.Errorf("failed to query Database (date failure): %v", city)
	}
	query := fmt.Sprintf("SELECT %v FROM weather_records WHERE timestamp='%v-%.2v-%.2v %.2v:00:00+00' AND city='%v'", value, year, month, day, hour, city)
	_, err = p.pgDB.Query(pg.Scan(&t), query)
	if err != nil {
		return err
	}
	return nil
}

func (p postgresDB) GetDatabase() pg.DB {
	return p.pgDB
}

func (p postgresDB) WeatherDataExists(city, date string) (bool, error) {
	now := time.Now()
	var timestamp string
	err := p.QueryDatabase(&timestamp, "timestamp", date, now.Hour(), strings.ToLower(city))
	if err != nil {
		return false, err
	}
	timeString := date + fmt.Sprintf(" %.2v:00:00+00", now.Hour())
	if strings.Contains(timestamp, timeString) {
		return true, nil
	}
	return false, nil
}

func (p postgresDB) GetWeatherRecords(city, date string) (model.WeatherRecord, error) {
	var today model.WeatherRecord
	city = strings.ToLower(city)
	year, month, day, err := utils.SplitDate(date)
	if err != nil {
		return model.WeatherRecord{}, err
	}
	utils.SetDate(year, month, day)
	dataExists, err := p.WeatherDataExists(city, date)
	if err != nil {
		return model.WeatherRecord{}, err
	}
	if !dataExists {
		records, err := utils.GetWeatherRecordsFromFiles(city)
		if err != nil {
			return model.WeatherRecord{}, err
		}
		if len(records) != 0 {
			if len(records[0].Hours) != 0 {
				timestamp := records[0].Hours[0].TimeStamp[:10]
				if timestamp != date {
					log.Print("INFO: Weather records don't exist! Getting new weather records from API Server.")
					today, err = utils.RequestWeather()
					if err != nil {
						return model.WeatherRecord{}, fmt.Errorf("failed to request weather: %v", err)
					}
					_, err = utils.SaveFutureWeatherInFile(city, date)
					if err != nil {
						return model.WeatherRecord{}, fmt.Errorf("failed to save future weather: %v", err)
					}
					utils.SetDate(year, month, day)
				}
			} else {
				log.Print("INFO: Weather records don't exist! Getting new weather records from API Server.")
				today, err = utils.RequestWeather()
				if err != nil {
					return model.WeatherRecord{}, fmt.Errorf("failed to request weather: %v", err)
				}
				_, err = utils.SaveFutureWeatherInFile(city, date)
				if err != nil {
					return model.WeatherRecord{}, fmt.Errorf("failed to save future weather: %v", err)
				}
				utils.SetDate(year, month, day)
			}
		}
		records, err = utils.GetWeatherRecordsFromFiles(city)
		if err != nil {
			return model.WeatherRecord{}, err
		}
		if len(records) == 0 {
			return model.WeatherRecord{}, fmt.Errorf("failed to get weather records! (internal error)")
		}
		for i := 0; i < len(records); i++ {
			err := p.InsertCityWeatherRecordsToTable(records[i])
			if err != nil {
				return model.WeatherRecord{}, err
			}
		}
	}
	if !utils.PathExists(fmt.Sprintf("resources/weather_records/%v_0-orig.json", city)) && dataExists {
		log.Print("INFO: Weather records don't exist! Getting weather records from Database.")
		_, err := p.getHourWeatherRecord(city, date)
		if err != nil {
			return model.WeatherRecord{}, err
		}
		_, err = utils.SaveFutureWeatherInFile(city, date)
		if err != nil {
			return model.WeatherRecord{}, fmt.Errorf("failed to save future weather: %v", err)
		}
	}
	today, err = p.getHourWeatherRecord(city, date)
	if err != nil {
		return model.WeatherRecord{}, err
	}
	return today, nil
}

func (p postgresDB) GetWeatherRecord(city, date string) (model.WeatherRecord, error) {
	var record model.WeatherRecord
	city = strings.ToLower(city)
	year, month, day, err := utils.SplitDate(date)
	if err != nil {
		return model.WeatherRecord{}, err
	}
	utils.SetDate(year, month, day)
	date = utils.GetDate()
	_, err = utils.CheckDate(date)
	if err != nil {
		return model.WeatherRecord{}, nil
	}
	dataExists, err := p.WeatherDataExists(city, date)
	if err != nil {
		return model.WeatherRecord{}, err
	}
	if !dataExists {
		record, err = utils.GetWeatherRecordFromFile(city)
		if err != nil {
			return model.WeatherRecord{}, err
		}
		if len(record.Hours) != 0 {
			timestamp := record.Hours[0].TimeStamp[:10]
			if timestamp != date {
				log.Print("INFO: Weather records don't exist! Getting new weather records from API Server.")
				record, err = utils.RequestWeather()
				if err != nil {
					return model.WeatherRecord{}, fmt.Errorf("failed to request weather: %v", err)
				}
				err = utils.SaveWeather(city, "0", record)
				if err != nil {
					return model.WeatherRecord{}, fmt.Errorf("failed to save future weather: %v", err)
				}
				utils.SetDate(year, month, day)
			}
		} else {
			log.Print("INFO: Weather records don't exist! Getting new weather records from API Server.")
			record, err = utils.RequestWeather()
			if err != nil {
				return model.WeatherRecord{}, fmt.Errorf("failed to request weather: %v", err)
			}
			err = utils.SaveWeather(city, "0", record)
			if err != nil {
				return model.WeatherRecord{}, fmt.Errorf("failed to save future weather: %v", err)
			}
			utils.SetDate(year, month, day)
		}
		if len(record.Hours) == 0 {
			return model.WeatherRecord{}, fmt.Errorf("failed to get weather records! (internal error)")
		}
		err = p.InsertCityWeatherRecordsToTable(record)
		if err != nil {
			return model.WeatherRecord{}, err
		}
	}
	if !utils.PathExists(fmt.Sprintf("resources/weather_records/%v_0-orig.json", city)) && dataExists {
		log.Print("INFO: Weather records don't exist! Getting weather records from Database.")
		record, err := p.getHourWeatherRecord(city, date)
		if err != nil {
			return model.WeatherRecord{}, err
		}
		err = utils.SaveWeather(city, "0", record)
		if err != nil {
			return model.WeatherRecord{}, fmt.Errorf("failed to save future weather: %v", err)
		}
	}
	if len(record.Hours) == 0 {
		record, err = p.getHourWeatherRecord(city, date)
		if err != nil {
			return model.WeatherRecord{}, err
		}
	}
	return record, nil
}

func (p postgresDB) getHourWeatherRecord(city, date string) (model.WeatherRecord, error) {
	city = strings.ToLower(city)
	records, err := p.QueryDayDatabase(city, date)
	if err != nil {
		return model.WeatherRecord{}, err
	}
	for i := range records {
		records[i].City = city
	}
	var today model.WeatherRecord
	today.Hours = records
	return today, nil
}

func (p postgresDB) InsertCityWeatherRecordsToTable(record model.WeatherRecord) error {
	city := utils.GetCityName()
	for i := 0; i < len(record.Hours); i++ {
		record.Hours[i].City = strings.ToLower(city)
	}
	res, err := p.pgDB.Model(&record.Hours).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return fmt.Errorf("failed to insert: %v", err)
	}
	log.Printf("INFO: Weather-Data inserted into DB: %d", res.RowsAffected())
	return nil
}

func (p postgresDB) Close() error {
	err := p.pgDB.Close()
	if err != nil {
		return err
	}
	log.Print("INFO: Connection closed")
	return nil
}

func (p postgresDB) SetLocationByCityName(city string) (string, error) {
	cityName, err := utils.SetLocationByCityName(city)
	if err != nil {
		return "", err
	}
	lcity := strings.ToLower(cityName)
	lat, lon, err := utils.GetLocation()
	if err != nil {
		return "", err
	}
	log.Printf("INFO: Location set to \"%v\" with Lat: %v, Lon: %v", lcity, lat, lon)
	err = p.InsertCityIntoDatabase(cityName)
	if err != nil {
		return "", err
	}
	return cityName, err
}

func (p postgresDB) InsertCityIntoDatabase(name string) error {
	cities := utils.GetCities()
	lat, lon, err := utils.GetLocation()
	cityName := strings.ToLower(name)
	city := model.City{
		Name: cityName,
		Lat:  lat,
		Lon:  lon,
	}
	_, err = p.pgDB.Model(&city).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return fmt.Errorf("failed insert: %v", err)
	}
	utils.SetCities(cities)
	log.Printf("INFO: City \"%v\" inserted into database", cityName)
	return nil
}

func (p postgresDB) QueryCitiesDatabase(t any, value, name string) error {
	query := fmt.Sprintf("SELECT %v FROM cities WHERE name='%v'", value, name)
	_, err := p.pgDB.Query(pg.Scan(&t), query)
	if err != nil {
		return err
	}
	return nil
}
