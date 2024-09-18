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

var pgDB pg.DB

type DBI interface {
	Query(city, date string)
}

func (f postgresDB) query(t any, value, city, date string, hour int) error {
	err := f.QueryDatabase(&t, value, date, hour, city)
	return err
}

func (f postgresDB) getDay(city, date string) (model.WeatherRecord, error) {
	var today model.WeatherRecord
	records, err := f.QueryDayDatabase(city, date)
	today.Hours = records
	return today, err
}

type postgresDB struct {
	pgDB pg.DB
}

func NewPG(user, password, name, address string) (postgresDB, error) {
	pgDB = connectToDatabase(user, password, name, address)
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

func (f postgresDB) QueryDayDatabase(city string, date string) ([]model.Weather_record, error) {
	var res []model.Weather_record
	year, month, day, err := utils.SplitDate(date)
	if err != nil {
		return []model.Weather_record{}, fmt.Errorf("failed to query Database (date): %v", err)
	}
	if city == "" {
		return []model.Weather_record{}, fmt.Errorf("failed to query Database because city isn't set!")
	}
	query := fmt.Sprintf("timestamp::date='%v-%.2v-%.2v 00:00:00+00' AND city='%v'", year, month, day, city)
	err = pgDB.Model().Table("weather_records").
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

func (f postgresDB) QueryDatabase(t any, value string, date string, hour int, city string) error {
	year, month, day, err := utils.SplitDate(date)
	if err != nil {
		return fmt.Errorf("failed to query Database (date failure): %v", city)
	}
	query := fmt.Sprintf("SELECT %v FROM weather_records WHERE timestamp='%v-%.2v-%.2v %.2v:00:00+00' AND city='%v'", value, year, month, day, hour, city)
	_, err = pgDB.Query(pg.Scan(&t), query)
	if err != nil {
		return err
	}
	return nil
}

func (f postgresDB) WeatherDataExists(city, date string) (bool, error) {
	now := time.Now()
	var timestamp string
	err := f.QueryDatabase(&timestamp, "timestamp", date, now.Hour(), strings.ToLower(city))
	if err != nil {
		return false, err
	}
	timeString := date + fmt.Sprintf(" %.2v:00:00+00", now.Hour())
	if strings.Contains(timestamp, timeString) {
		return true, nil
	}
	return false, nil
}

func (f postgresDB) GetWeatherRecord(city, date string) (model.WeatherRecord, error) {
	var today model.WeatherRecord
	city = strings.ToLower(city)
	dataExists, err := f.WeatherDataExists(city, date)
	if err != nil {
		return model.WeatherRecord{}, err
	}
	//if !dataExists {
	if !utils.PathExists(fmt.Sprintf("resources/weather_records/%v_0-orig.json", strings.ToLower(city))) {
		log.Print("INFO: Weather records don't exist! Getting new weather records from API Server.")
		today, err = utils.RequestWeather()
		if err != nil {
			return model.WeatherRecord{}, fmt.Errorf("failed to request weather: %v", err)
		}
		_, err := utils.SaveFutureWeatherInFile(city, date)
		if err != nil {
			return model.WeatherRecord{}, fmt.Errorf("failed to save future weather: %v", err)
		}
	}
	records, err := utils.GetWeatherRecordsFromFiles(city)
	if err != nil {
		return model.WeatherRecord{}, err
	}
	for i := 0; i < len(records); i++ {
		err := f.InsertCityWeatherRecordsToTable(records[i])
		if err != nil {
			return model.WeatherRecord{}, err
		}
	}
	year, month, day, err := utils.SplitDate(date)
	if err != nil {
		return model.WeatherRecord{}, err
	}
	utils.SetDate(year, month, day)
	//}
	if !utils.PathExists(fmt.Sprintf("resources/weather_records/%v_0-orig.json", city)) && dataExists {
		log.Print("INFO: Weather records don't exist! Getting weather records from Database.")
		_, err := f.getHourWeatherRecord(city, date)
		if err != nil {
			return model.WeatherRecord{}, err
		}
	}
	today, err = f.getHourWeatherRecord(city, date)
	if err != nil {
		return model.WeatherRecord{}, err
	}
	return today, nil
}

func (f postgresDB) getHourWeatherRecord(city, date string) (model.WeatherRecord, error) {
	city = strings.ToLower(city)
	_, err := f.QueryDayDatabase(city, date)
	if err != nil {
		return model.WeatherRecord{}, err
	}
	var today model.WeatherRecord
	//today.Hours = records
	return today, nil
}

func (f postgresDB) InsertCityWeatherRecordsToTable(record model.WeatherRecord) error {
	city := utils.GetCityName()
	for i := 0; i < len(record.Hours); i++ {
		record.Hours[i].City = strings.ToLower(city)
	}
	res, err := pgDB.Model(&record.Hours).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return fmt.Errorf("failed to insert: %v", err)
	}
	log.Printf("INFO: Weather-Data inserted into DB %d", res.RowsAffected())
	return nil
}

func (f postgresDB) Close() {
	pgDB.Close()
}

func (f postgresDB) SetLocationByCityName(city string) (string, error) {
	cityName, err := utils.SetLocationByCityName(city)
	if err != nil {
		return "", err
	}
	err = f.InsertIntoCitiesDatabase(city, city)
	if err != nil {
		return "", err
	}
	return cityName, err
}

func (f postgresDB) InsertIntoCitiesDatabase(value, name string) error {
	cities := utils.GetCities()
	lat, lon, err := utils.GetLocation()
	cityName := strings.ToLower(name)
	city := model.Cities{
		Name: cityName,
		Lat:  lat,
		Lon:  lon,
	}
	_, err = pgDB.Model(&city).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return fmt.Errorf("failed insert: %v", err)
	}
	utils.SetCities(cities)
	return nil
}

func (f postgresDB) QueryCitiesDatabase(t any, value, name string) error {
	query := fmt.Sprintf("SELECT %v FROM cities WHERE name='%v'", value, name)
	_, err := pgDB.Query(pg.Scan(&t), query)
	if err != nil {
		return err
	}
	return nil
}

/* OLD WAY didnt quite work
path := fmt.Sprintf("resources/pg/data/%v.csv", city)
if !utils.PathExists(path) {
	return fmt.Errorf("ERROR: Path \"%v\" ", path)
}
utils.ConvertWeatherRecords()
file, err := os.Open(path)
if err != nil {
	return fmt.Errorf("ERROR: Couldn't open file (%v)", path)
}
defer file.Close()
csvReader := csv.NewReader(file)
records, err := csvReader.ReadAll()
if err != nil {
	return fmt.Errorf("ERROR: Couldn't read file %v.csv\nERROR: %v", city, err)
}
_, err = pgDB.Exec(`CREATE TEMPORARY TABLE IF NOT EXISTS temp_weather_records (id INT NOT NULL, timestamp TIMESTAMP, source_id INT, precipitation FLOAT, pressure_msl FLOAT, sunshine FLOAT, temperature FLOAT, wind_direction INT, wind_speed FLOAT, cloud_cover FLOAT, dew_point FLOAT, relative_humidity FLOAT, visibility FLOAT, wind_gust_direction INT, wind_gust_speed FLOAT, condition VARCHAR(100), precipitation_probability FLOAT, precipitation_probability_6h FLOAT, solar FLOAT, icon VARCHAR(100), city VARCHAR(100), PRIMARY KEY(ID));`)
if err != nil {
	return fmt.Errorf("ERROR: Couldn't create temporary table temp_weather_records\nERROR: %v", err)
}
var csvString []string
var count int
_, err = pgDB.Query(pg.Scan(&count), "SELECT COUNT(*) FROM weather_records")
if err != nil {
	return fmt.Errorf("ERROR: Count failed!\nERROR: %v", err)
}
for _, inner := range records {
	inner = append([]string{strconv.Itoa(count + 1)}, inner...)
	csvString = append(csvString, strings.Join(inner, ","))
	count++
}
csvData := strings.Join(csvString, "\n")
reader := strings.NewReader(csvData)
_, err = pgDB.CopyFrom(reader, `COPY temp_weather_records FROM STDIN WITH CSV`)
if err != nil {
	return fmt.Errorf("ERROR: Couldn't copy temp_weather_records\nERROR: %v", err)
}
_, err = pgDB.Exec("INSERT INTO weather_records SELECT * FROM temp_weather_records ON CONFLICT (id, timestamp, city) DO NOTHING")
if err != nil {
	return fmt.Errorf("failed to insert temp_weather_records into weather_records: %v", err)
}
_, err = pgDB.Exec("DROP TABLE temp_weather_records")
if err != nil {
	return fmt.Errorf("failed to drop temp_weather_records: %v", err)
}*/
