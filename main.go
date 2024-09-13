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
var db *pg.DB

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
	TimeStamp                  string  `json:"timestamp"`
	SourceID                   int     `json:"source_id"`
	Precipitation              float64 `json:"precipitation"`
	PressureMSL                float64 `json:"pressure_msl"`
	Sunshine                   float64 `json:"sunshine"`
	Temperature                float64 `json:"temperature"`
	WindDirection              int     `json:"wind_direction"`
	WindSpeed                  float64 `json:"wind_speed"`
	CloudCover                 float64 `json:"cloud_cover"`
	DewPoint                   float64 `json:"dew_point"`
	RelativeHumidity           int     `json:"relative_humidity"`
	Visibility                 int     `json:"visibility"`
	WindGustDirection          int     `json:"wind_gust_direction"`
	WindGustSpeed              float64 `json:"wind_gust_speed"`
	Condition                  string  `json:"condition"`
	PrecipitationProbability   float64 `json:"precipitation_probability"`
	PrecipitationProbability6h float64 `json:"precipitation_probability_6h"`
	Solar                      float64 `json:"solar"`
	Icon                       string  `json:"icon"`
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
	file, err := os.Open("resources/json-data/cities.json")
	if err != nil {
		saveCityByName(name, cities)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cities)
	if err != nil {
		fmt.Println(err)
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
	saveJSONFile("resources/json-data", "cities.json", data)
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
	if !(pathExists("resources/weather_records/")) {
		fmt.Println("INFO: Weather records don't exist! Getting new weather records from API Server.")
		requestWeather()
		saveFutureWeatherInFile(cityName, date)
	}
	connectToDatabase()
	/*minutesRequest, err := strconv.Atoi(opts.MinutesRequest)
	if err != nil {
		log.Fatal(err)
	}
	requestWeatherEvery(time.Duration(minutesRequest*int(time.Minute)), showWeather)*/
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

func saveJSONFile(directory, filename string, data []byte) error {
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

	err = saveJSONFile("resources/weather_records", strings.ToLower(city)+"_"+count+"-orig.json", data)
	if err != nil {
		log.Fatal(err)
	}
}

func saveTodaysWeather(city string, count string) {
	data, err := json.MarshalIndent(today, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = saveJSONFile("resources/weather_records", strings.ToLower(city)+"_"+count+"-orig.json", data)
	if err != nil {
		log.Fatal(err)
	}
}

func saveFutureWeatherInFile(city string, date string) {
	count := "0"
	requestWeather()
	saveTodaysWeather(city, count)
	newDate := date
	newDate, count = setFutureDay(newDate, count) // Create the first next day
	requestFutureWeather()
	saveFutureWeather(city, count)
	for i := 1; i <= 6; i++ { // Create for the next 6 days after the first
		newDate, count = setFutureDay(newDate, count)
		requestFutureWeather()
		saveFutureWeather(city, count)
	}
	year, month, day := splitDate(date)
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
	day += 1
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

func connectToDatabase() {
	db = pg.Connect(&pg.Options{
		Addr:     FDB_ADDRESS,
		User:     FDB_USER,
		Password: FDB_PASS,
		Database: FDB_DB,
	})
	defer db.Close()
	var day WeatherRecord
	var hours = [25]HourWeatherRecord{}
	day.Hours = hours[:]
	for i := 0; i <= 24; i++ {
		getHourWeatherRecord(day, i, db)
	}
	insertCityWeatherRecordsToTable(strings.ToLower(cityName), db)
}

func getHourWeatherRecord(day WeatherRecord, hour int, db *pg.DB) {
	var city, timestamp, condition, icon string
	var source_id, wind_direction, relative_humidity, visibility, wind_gust_direction int
	var precipitation, pressuemsl, sunshine, temperature, wind_speed, cloud_cover, dew_point,
		wind_gust_speed, precipitation_probability, precipitation_probability_6h, solar float64
	queryDatabase(&timestamp, "timestamp", hour, db)
	queryDatabase(&source_id, "source_id", hour, db)
	queryDatabase(&precipitation, "precipitation", hour, db)
	queryDatabase(&pressuemsl, "pressure_msl", hour, db)
	queryDatabase(&sunshine, "sunshine", hour, db)
	queryDatabase(&temperature, "temperature", hour, db)
	queryDatabase(&wind_direction, "wind_direction", hour, db)
	queryDatabase(&wind_speed, "wind_speed", hour, db)
	queryDatabase(&cloud_cover, "cloud_cover", hour, db)
	queryDatabase(&dew_point, "dew_point", hour, db)
	queryDatabase(&relative_humidity, "relative_humidity", hour, db)
	queryDatabase(&visibility, "visibility", hour, db)
	queryDatabase(&wind_gust_direction, "wind_gust_direction", hour, db)
	queryDatabase(&wind_gust_speed, "wind_gust_speed", hour, db)
	queryDatabase(&condition, "condition", hour, db)
	queryDatabase(&precipitation_probability, "precipitation_probability", hour, db)
	queryDatabase(&precipitation_probability_6h, "precipitation_probability_6h", hour, db)
	queryDatabase(&solar, "solar", hour, db)
	queryDatabase(&icon, "icon", hour, db)
	queryDatabase(&city, "city", hour, db)
	day.Hours[hour].TimeStamp = timestamp
	day.Hours[hour].SourceID = source_id
	day.Hours[hour].Precipitation = precipitation
	day.Hours[hour].PressureMSL = pressuemsl
	day.Hours[hour].Sunshine = sunshine
	day.Hours[hour].Temperature = temperature
	day.Hours[hour].WindDirection = wind_direction
	day.Hours[hour].WindSpeed = wind_speed
	day.Hours[hour].CloudCover = cloud_cover
	day.Hours[hour].DewPoint = dew_point
	day.Hours[hour].RelativeHumidity = relative_humidity
	day.Hours[hour].Visibility = visibility
	day.Hours[hour].WindGustDirection = wind_gust_direction
	day.Hours[hour].WindGustSpeed = wind_gust_speed
	day.Hours[hour].Condition = condition
	day.Hours[hour].PrecipitationProbability = precipitation_probability
	day.Hours[hour].PrecipitationProbability6h = precipitation_probability_6h
	day.Hours[hour].Solar = solar
	day.Hours[hour].Icon = icon
	today.Hours = day.Hours
	today = day
}

func queryDatabase(t interface{}, value string, hour int, db *pg.DB) (interface{}, error) {
	year, month, day := splitDate(date)
	query := fmt.Sprintf("SELECT %v FROM weather_records WHERE timestamp='%v-%.2v-%.2v %.2v:00:00+00'", value, year, month, day, hour)
	_, err := db.Query(pg.Scan(&t), query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return t, nil
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func insertCityWeatherRecordsToTable(city string, db *pg.DB) {
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

/*
TODO:
-	Change log.Fatalf
	and give back error
*/
