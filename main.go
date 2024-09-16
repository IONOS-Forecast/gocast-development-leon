package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
)

var weatherAPIURL string
var geoAPIURL string
var geoAPIKey string
var FDB_USER string
var FDB_PASS string
var FDB_DB string
var FDB_ADDRESS string
var latitude = 52.5170365
var longitude = 13.3888599
var date string
var cityName string
var today WeatherRecord
var notToday WeatherRecord
var opts options
var pgdb *pg.DB

type options struct {
	WeatherAPIURL  string `v:"wapiurl" long:"weather-api-url" env:"WEATHER_API_URL" description:"URL to interact with Weather provider"`
	GeoAPIURL      string `v:"gapiurl" long:"geo-api-url" env:"GEO_API_URL" description:"URL to interact with GEO provider"`
	GeoAPIKEY      string `v:"gapikey" long:"geo-api-key" env:"GEO_API_KEY" description:"KEY to interact with GEO provider"`
	MinutesRequest string `v:"reqmin" long:"req-aft-min" env:"REQ_AFT_MIN" description:"Minutes until the next request to the Weather provicer is made"`
	FDB_USER       string `v:"fdb-u" long:"fdb-user" env:"FDB_USER" description:"The user that connects to the forecast-database"`
	FDB_PASSWORD   string `v:"fdb-p" long:"fdb_password" env:"FDB_PASSWORD" description:"The password to the user for connecting to the forecast-database"`
	FDB_DATABASE   string `v:"fdb-db" long:"fdb_database" env:"FDB_DATABASE" description:"The database that the user connects to"`
	FDB_ADDRESS    string `v:"fdb-addr" long:"fdb_address" env:"FDB_ADDRESS" description:"The address of the database"`
}

type WeatherRecord struct {
	Hours []HourWeatherRecord `json:"weather"`
}

type HourWeatherRecord struct {
	TimeStamp                  string  `json:"timestamp" pg:"timestamp"`
	SourceID                   int     `json:"source_id" pg:"source_id"`
	Precipitation              float64 `json:"precipitation" pg:"precipitation"`
	PressureMSL                float64 `json:"pressure_msl" pg:"pressure_msl"`
	Sunshine                   float64 `json:"sunshine" pg:"sunshine"`
	Temperature                float64 `json:"temperature" pg:"temperature"`
	WindDirection              int     `json:"wind_direction" pg:"wind_direction"`
	WindSpeed                  float64 `json:"wind_speed" pg:"wind_speed"`
	CloudCover                 float64 `json:"cloud_cover" pg:"cloud_cover"`
	DewPoint                   float64 `json:"dew_point" pg:"dew_point"`
	RelativeHumidity           int     `json:"relative_humidity" pg:"relative_humidity"`
	Visibility                 int     `json:"visibility" pg:"visibility"`
	WindGustDirection          int     `json:"wind_gust_direction" pg:"wind_gust_direction"`
	WindGustSpeed              float64 `json:"wind_gust_speed" pg:"wind_gust_speed"`
	Condition                  string  `json:"condition" pg:"condition"`
	PrecipitationProbability   float64 `json:"precipitation_probability" pg:"precipitation_probability"`
	PrecipitationProbability6h float64 `json:"precipitation_probability_6h" pg:"precipitation_probability_6h"`
	Solar                      float64 `json:"solar" pg:"solar"`
	Icon                       string  `json:"icon" pg:"icon"`
}

type OWCity struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Country   string  `json:"country"`
}

type City struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func setDateAndLocationByCityName(year, month, day int, cityName string, cities map[string]City) {
	setDate(year, month, day)
	setLocationByCityName(cityName, cities)
}

func setDateAndLocation(year, month, day int, lat, lon float64) {
	setDate(year, month, day)
	setLocation(lat, lon)
}

func setDate(year, month, day int) {
	now := time.Now()
	_date := fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	_, err := checkDate(_date)
	if err != nil { // What happens when date is invalid
		date = now.Format("2006-01-02")
		fmt.Println("Error: The date is invalid!")
	}
	if year >= 2010 && year <= now.Year() && month >= 1 && month <= 12 && day >= 1 && day <= 31 {
		date = _date
	} else { // What happens when the date is out of range
		date = now.Format("2006-01-02")
		fmt.Println("Error: The given date is out of range.")
	}
	reloadWeatherURL()
}

func setLocation(lat, lon float64) {
	if lat <= 54 && lat >= 48 && lon <= 14 && lon >= 6 {
		latitude = lat
		longitude = lon
		reloadWeatherURL()
	} else { // When location is not in range
		fmt.Println("Error: The location is not in germany or not in range!")
	}
}

func setLocationByCityName(name string, cities map[string]City) {
	cities = readCities(name, cities)
	if city, exists := cities[strings.ToLower(name)]; exists { // When the city exists
		cityName = name
		setLocation(city.Lat, city.Lon)
	} else { // When the city doesn't exist
		saveCityByName(name, cities)
		setLocationByCityName(name, cities)
	}
}

func readCities(name string, cities map[string]City) map[string]City {
	file, err := os.Open("resources/data/cities.json")
	if err != nil {
		saveCityByName(name, cities)
		return readCities(name, cities)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cities)
	if err != nil {
		panic(err)
	}
	return cities
}

func saveCityByName(name string, cities map[string]City) string {
	cityName = name
	reloadGeoURL()
	var owcities []OWCity
	resp, err := http.Get(geoAPIURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &owcities)
	if err != nil {
		log.Fatal(err)
	}
	var foundcity OWCity
	for _, owcity := range owcities {
		if owcity.Country == "DE" {
			foundcity = owcity
		}
	}
	cities[strings.ToLower(name)] = City{Lat: foundcity.Latitude, Lon: foundcity.Longitude}
	data, err := json.MarshalIndent(cities, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	resourcesPath := "resources/data"
	saveFile(resourcesPath, "cities.json", data)
	citiesPath := "resources/data/cities.txt"
	if !pathExists(citiesPath) {
		var citiesData []byte
		for _, v := range []byte(strings.ToLower(name)) {
			citiesData = append(citiesData, v)
		}
		saveFile(resourcesPath, "cities.txt", citiesData)
	} else {
		file, err := os.Open(citiesPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		content, err := os.ReadFile(citiesPath)
		if err != nil {
			panic(err)
		}
		var citiesData []byte
		if !strings.Contains(string(content), strings.ToLower(name)) {
			citiesString := string(content) + "\n" + strings.ToLower(name)
			for _, v := range []byte(citiesString) {
				citiesData = append(citiesData, v)
			}
			saveFile("resources/data", "cities.txt", citiesData)
		}
	}
	return foundcity.Name
}

func reloadWeatherURL() {
	u, err := url.Parse(weatherAPIURL)
	if err != nil {
		log.Fatal(err)
	}

	v := u.Query()
	v.Set("lat", strconv.FormatFloat(latitude, 'f', 2, 64))
	v.Set("lon", strconv.FormatFloat(longitude, 'f', 2, 64))
	v.Set("date", date)
	u.RawQuery = v.Encode()
	weatherAPIURL = u.String()
}

func reloadGeoURL() {
	u, err := url.Parse(geoAPIURL)
	if err != nil {
		log.Fatal(err)
	}

	v := u.Query()
	v.Set("q", cityName)
	v.Set("appid", geoAPIKey)
	u.RawQuery = v.Encode()
	geoAPIURL = u.Redacted()
}

func main() {
	cities := make(map[string]City)
	godotenv.Load()
	godotenv.Load(".env.dev")
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	weatherAPIURL = opts.WeatherAPIURL
	geoAPIURL = opts.GeoAPIURL
	geoAPIKey = opts.GeoAPIKEY
	FDB_USER = opts.FDB_USER
	FDB_PASS = opts.FDB_PASSWORD
	FDB_DB = opts.FDB_DATABASE
	FDB_ADDRESS = opts.FDB_ADDRESS
	now := time.Now()
	setDate(now.Year(), int(now.Month()), now.Day())
	if cityName == "" {
		setLocationByCityName("Berlin", cities)
	}
	database := connectToDatabase()
	defer database.Close()
	getWeatherRecord(cityName, database)
	/*minutesRequest, err := strconv.Atoi(opts.MinutesRequest)
	if err != nil {
		log.Fatal(err)
	}
	requestWeatherEvery(time.Duration(minutesRequest*int(time.Minute)), showWeather)*/

}

func getWeatherRecord(city string, db pg.DB) {
	city = strings.ToLower(city)
	if !weatherDataExists(city, db) {
		fmt.Println("INFO: Weather records don't exist! Getting new weather records from API Server.")
		requestWeather()
		saveFutureWeatherInFile(cityName, date)
		insertCityWeatherRecordsToTable(city, db)
	}
	if !pathExists("resources/weather_records/berlin_0-orig.json") && weatherDataExists(city, db) {
		fmt.Println("INFO: Weather records don't exist! Getting weather records from Database.")
		getHourWeatherRecord(city, db)
	}
	getHourWeatherRecord(city, db)
}

func requestWeatherEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		requestWeather()
		f(x)
	}
}

func checkDate(s string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", s)
	if err != nil {
		return date, err
	}
	return date, nil
}

func saveFile(directory, filename string, data []byte) error {
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	filePath := fmt.Sprintf("%s/%s", directory, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func requestWeather() {
	resp, err := http.Get(weatherAPIURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &today)
	if err != nil {
		log.Fatal(err)
	}
}

func requestFutureWeather() {
	resp, err := http.Get(weatherAPIURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &notToday)
	if err != nil {
		log.Fatal(err)
	}
}

func saveFutureWeather(city string, count string) {
	data, err := json.MarshalIndent(notToday, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = saveFile("resources/weather_records", strings.ToLower(city)+"_"+count+"-orig.json", data)
	if err != nil {
		log.Fatal(err)
	}
}

func saveTodaysWeather(city string, count string) {
	data, err := json.MarshalIndent(today, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = saveFile("resources/weather_records", strings.ToLower(city)+"_"+count+"-orig.json", data)
	if err != nil {
		log.Fatal(err)
	}
}

func saveFutureWeatherInFile(city string, date string) {
	count := "0"
	requestWeather()
	saveTodaysWeather(city, count)
	newDate := date
	for i := 1; i <= 7; i++ { // Create for the next 6 days after the first
		newDate, count = setFutureDay(newDate, count)
		requestFutureWeather()
		saveFutureWeather(city, count)
	}
	year, month, day := splitDate(date)
	day += 1
	setDate(year, month, day)
}

func setFutureDay(date string, count string) (string, string) {
	splitDate := strings.Split(date, "-")
	day, err := strconv.Atoi(splitDate[2])
	if err != nil {
		log.Fatalf("Couldn't convert Day: %v", err)
	}
	day += 1
	month, err := strconv.Atoi(splitDate[1])
	if err != nil {
		log.Fatalf("Couldn't convert Month: %v", err)
	}
	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		log.Fatalf("Couldn't convert Year: %v", err)
	}
	_count, err := strconv.Atoi(count)
	if err != nil {
		log.Fatalf("Couldn't convert Counter: %v", err)
	}
	_count += 1
	count = strconv.Itoa(_count)
	setDate(year, month, day)
	return fmt.Sprintf("%v-%.2v-%.2v", year, month, day), count
}

func splitDate(date string) (year, month, day int) {
	splitDate := strings.Split(date, "-")
	day, err := strconv.Atoi(splitDate[2])
	if err != nil {
		log.Fatalf("Couldn't convert Day: %v", err)
	}
	month, err = strconv.Atoi(splitDate[1])
	if err != nil {
		log.Fatalf("Couldn't convert Month: %v", err)
	}
	year, err = strconv.Atoi(splitDate[0])
	if err != nil {
		log.Fatalf("Couldn't convert Year: %v", err)
	}
	return year, month, day
}

func connectToDatabase() pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     FDB_ADDRESS,
		User:     FDB_USER,
		Password: FDB_PASS,
		Database: FDB_DB,
	})

	return *db
}

func getHourWeatherRecord(city string, db pg.DB) {
	city = strings.ToLower(city)
	records, err := queryDayDatabase(city, db)
	if err != nil {
		panic(err)
	}
	today.Hours = records
}

func queryDayDatabase(city string, db pg.DB) ([]HourWeatherRecord, error) {
	var res []HourWeatherRecord
	year, month, day := splitDate(date)
	query := fmt.Sprintf("timestamp::date='%v-%.2v-%.2v 00:00:00+00' AND city='%v'", year, month, day, city)
	err := db.Model().Table("weather_records").
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

func queryDatabase(t any, value string, hour int, city string, db pg.DB) error {
	year, month, day := splitDate(date)
	query := fmt.Sprintf("SELECT %v FROM weather_records WHERE timestamp='%v-%.2v-%.2v %.2v:00:00+00' AND city='%v'", value, year, month, day, hour, city)
	_, err := db.Query(pg.Scan(&t), query)
	if err != nil {
		return err
	}
	return nil
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func insertCityWeatherRecordsToTable(city string, db pg.DB) {
	path := fmt.Sprintf("resources/pg/data/%v.csv", city)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TEMPORARY TABLE temp_weather_records (id int NOT NULL, timestamp TIMESTAMP, source_id INT, precipitation FLOAT, pressure_msl FLOAT, sunshine FLOAT, temperature FLOAT, wind_direction INT, wind_speed FLOAT, cloud_cover FLOAT, dew_point FLOAT, relative_humidity FLOAT, visibility FLOAT, wind_gust_direction INT, wind_gust_speed FLOAT, condition VARCHAR(100), precipitation_probability FLOAT, precipitation_probability_6h FLOAT, solar FLOAT, icon VARCHAR(100), city VARCHAR(100), PRIMARY KEY(ID));`)
	if err != nil {
		log.Fatal("Exec:", err)
	}
	var csvString []string
	var count int
	_, err = db.Query(pg.Scan(&count), "SELECT COUNT(*) FROM weather_records")
	if err != nil {
		log.Fatal("Count FAILED!:", err)
	}

	for _, inner := range records {
		inner = append([]string{strconv.Itoa(count + 1)}, inner...)
		csvString = append(csvString, strings.Join(inner, ","))
		count++
	}
	csvData := strings.Join(csvString, "\n")
	reader := strings.NewReader(csvData)
	_, err = db.CopyFrom(reader, `COPY temp_weather_records FROM STDIN WITH CSV`)
	if err != nil {
		log.Fatal("CopyFrom:", err)
	}

	_, err = db.Exec("INSERT INTO weather_records\nSELECT * FROM temp_weather_records WHERE timestamp NOT IN (SELECT timestamp FROM weather_records)")
	if err != nil {
		log.Fatal("Exec-Overwrite:", err)
	}
	fmt.Println("Weather Data inserted into Table!")
}

func weatherDataExists(city string, db pg.DB) bool {
	now := time.Now()
	var timestamp string
	queryDatabase(&timestamp, "timestamp", now.Hour(), strings.ToLower(city), db)
	timeString := date + fmt.Sprintf(" %.2v:00:00+00", now.Hour())
	if strings.Contains(timestamp, timeString) {
		return true
	}
	return false
}

/*
TODO:
-	Change log.Fatalf
	and give back error
*/
