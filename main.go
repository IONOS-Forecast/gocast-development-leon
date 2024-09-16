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

func setDateAndLocationByCityName(year, month, day int, cityName string, cities map[string]City) error {
	err := setDate(year, month, day)
	if err != nil {
		return err
	}
	err = setLocationByCityName(cityName, cities)
	if err != nil {
		return err
	}
	return nil
}

func setDateAndLocation(year, month, day int, lat, lon float64) error {
	err := setDate(year, month, day)
	if err != nil {
		return err
	}
	err = setLocation(lat, lon)
	if err != nil {
		return err
	}
	return nil
}

func setDate(year, month, day int) error {
	now := time.Now()
	_date := fmt.Sprintf("%v-%.2v-%.2v", year, month, day)
	_, err := checkDate(_date)
	if err != nil { // What happens when date is invalid
		date = now.Format("2006-01-02")
		return fmt.Errorf("ERROR: The date \"%v\" is invalid!\nWARNING: Date set to today (\"%v\")", _date, date)
	}
	if year >= 2010 && year <= now.Year() && month >= 1 && month <= 12 && day >= 1 && day <= 31 {
		date = _date
	} else { // What happens when the date is out of range
		date = now.Format("2006-01-02")
		return fmt.Errorf("ERROR: The given date \"%v\" is out of range.\nWARNING: Date set to today (\"%v\")\n", _date, date)
	}
	err = reloadWeatherURL()
	if err != nil {
		return err
	}
	return nil
}

func setLocation(lat, lon float64) error {
	if lat <= 54 && lat >= 48 && lon <= 14 && lon >= 6 {
		latitude = lat
		longitude = lon
		err := reloadWeatherURL()
		if err != nil {
			return err
		}
		return nil
	} else { // When location is not in range
		return fmt.Errorf("ERROR: The location \"Lat:%f; Lon:%f\" is not in germany or not in range!", lat, lon)
	}
}

func setLocationByCityName(name string, cities map[string]City) error {
	cities = readCities(name, cities)
	if city, exists := cities[strings.ToLower(name)]; exists { // When the city exists
		cityName = name
		err := setLocation(city.Lat, city.Lon)
		if err != nil {
			return err
		}
		return nil
	} else { // When the city doesn't exist
		_, err := saveCityByName(name, cities)
		if err != nil {
			return err
		}
		err = setLocationByCityName(name, cities)
		if err != nil {
			return err
		}
		return fmt.Errorf("INFO: City \"%v\" doesn't exist!\nINFO: Getting city \"%v\" from API!", strings.ToLower(name), strings.ToLower(name))
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

func saveCityByName(name string, cities map[string]City) (string, error) {
	oldCityName := cityName
	cityName = name
	err := reloadGeoURL()
	if err != nil {
		return "", err
	}
	var owcities []OWCity
	resp, err := http.Get(geoAPIURL)
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
	var foundcity OWCity
	for _, owcity := range owcities {
		if owcity.Country == "DE" {
			foundcity = owcity
		}
	}
	err = setLocation(foundcity.Latitude, foundcity.Longitude)
	if err != nil {
		cityName = oldCityName
		return "", fmt.Errorf("ERROR: The given location is not in germany!\n")
	}
	cities[strings.ToLower(name)] = City{Lat: foundcity.Latitude, Lon: foundcity.Longitude}
	data, err := json.MarshalIndent(cities, "", "  ")
	if err != nil {
		cityName = oldCityName
		return "", err
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
			saveFile("resources/data", "cities.txt", citiesData)
		}
	}
	return foundcity.Name, nil
}

func reloadWeatherURL() error {
	u, err := url.Parse(weatherAPIURL)
	if err != nil {
		return fmt.Errorf("ERROR: GeocodingAPIURL has incorrect values!\nERROR: %v", err)
	}

	v := u.Query()
	v.Set("lat", strconv.FormatFloat(latitude, 'f', 2, 64))
	v.Set("lon", strconv.FormatFloat(longitude, 'f', 2, 64))
	v.Set("date", date)
	u.RawQuery = v.Encode()
	weatherAPIURL = u.String()
	return nil
}

func reloadGeoURL() error {
	u, err := url.Parse(geoAPIURL)
	if err != nil {
		return fmt.Errorf("ERROR: GeocodingAPIURL has incorrect values!\nERROR: %v", err)
	}

	v := u.Query()
	v.Set("q", cityName)
	v.Set("appid", geoAPIKey)
	u.RawQuery = v.Encode()
	geoAPIURL = u.Redacted()
	return nil
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
	err = setDate(now.Year(), int(now.Month()), now.Day())
	if err != nil {
		log.Print(err)
	}
	if cityName == "" {
		err = setLocationByCityName("Berlin", cities)
		if err != nil {
			log.Print(err)
		}
	}
	database := connectToDatabase()
	defer database.Close()
	err = getWeatherRecord(cityName, database)
	if err != nil {
		log.Print(err)
	}
	// Test with Second city
	err = setLocationByCityName("München", cities)
	if err != nil {
		log.Print(err)
	}
	err = getWeatherRecord(cityName, database)
	if err != nil {
		log.Print(err)
	}
	// Error Examples
	fmt.Println(":----------------------------------------------------------------------------:")
	fmt.Println("\t\t\tERROR EXAMPLES BEGIN HERE")
	fmt.Println(":----------------------------------------------------------------------------:")
	err = setDate(2024, 2, 30)
	if err != nil {
		log.Print(err)
	}
	err = setLocationByCityName("Afrika", cities)
	if err != nil {
		log.Print(err)
	}
	fmt.Println(":----------------------------------------------------------------------------:")
	fmt.Println("\t\t\t ERROR EXAMPLES END HERE")
	fmt.Println(":----------------------------------------------------------------------------:")
	/*minutesRequest, err := strconv.Atoi(opts.MinutesRequest)
	if err != nil {
		log.Fatal(err)
	}
	requestWeatherEvery(time.Duration(minutesRequest*int(time.Minute)), showWeather)*/

}

func getWeatherRecord(city string, db pg.DB) error {
	city = strings.ToLower(city)
	dataExists, err := weatherDataExists(city, db)
	if err != nil {
		return err
	}
	if !dataExists {
		fmt.Println("INFO: Weather records don't exist! Getting new weather records from API Server.")
		err = requestWeather()
		if err != nil {
			return fmt.Errorf("ERROR: Requesting weather threw an error!\nERROR: %v", err)
		}
		err = saveFutureWeatherInFile(cityName, date)
		if err != nil {
			return fmt.Errorf("ERROR: Saving future weather threw an error!\nERROR: %v", err)
		}
		err := insertCityWeatherRecordsToTable(city, db)
		if err != nil {
			return err
		}
	}
	if !pathExists(fmt.Sprintf("resources/weather_records/%v_0-orig.json", city)) && dataExists {
		fmt.Println("INFO: Weather records don't exist! Getting weather records from Database.")
		err := getHourWeatherRecord(city, db)
		if err != nil {
			return err
		}
	}
	err = getHourWeatherRecord(city, db)
	if err != nil {
		return err
	}
	return nil
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

func requestWeather() error {
	resp, err := http.Get(weatherAPIURL)
	if err != nil {
		return fmt.Errorf("ERROR: Couldn't get weatherAPI from URL \"%v\"\nERROR: %v", weatherAPIURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ERROR: Couldn't read Response from URL \"%v\"\nERROR: %v", weatherAPIURL, err)
	}
	err = json.Unmarshal(body, &today)
	if err != nil {
		return fmt.Errorf("ERROR: Couldn't unmarshal Response-Body from URL \"%v\"\nERROR: %v", weatherAPIURL, err)
	}
	return nil
}

func requestFutureWeather() error {
	resp, err := http.Get(weatherAPIURL)
	if err != nil {
		return fmt.Errorf("ERROR: Couldn't get weatherAPI from URL \"%v\"\nERROR: %v", weatherAPIURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ERROR: Couldn't read Response from URL \"%v\"\nERROR: %v", weatherAPIURL, err)
	}
	err = json.Unmarshal(body, &notToday)
	if err != nil {
		return fmt.Errorf("ERROR: Couldn't unmarshal Response-Body from URL \"%v\"\nERROR: %v", weatherAPIURL, err)
	}
	return nil
}

func saveFutureWeather(city string, count string) error {
	data, err := json.MarshalIndent(notToday, "", "  ")
	if err != nil {
		return fmt.Errorf("ERROR: MarshalIndent threw an error!\nERROR: %v", err)
	}

	err = saveFile("resources/weather_records", strings.ToLower(city)+"_"+count+"-orig.json", data)
	if err != nil {
		return err
	}
	return nil
}

func saveTodaysWeather(city string, count string) error {
	data, err := json.MarshalIndent(today, "", "  ")
	if err != nil {
		return fmt.Errorf("ERROR: MarshalIndent threw an error!\nERROR: %v", err)
	}

	err = saveFile("resources/weather_records", strings.ToLower(city)+"_"+count+"-orig.json", data)
	if err != nil {
		return err
	}
	return nil
}

func saveFutureWeatherInFile(city string, date string) error {
	count := "0"
	err := requestWeather()
	if err != nil {
		return err
	}
	err = saveTodaysWeather(city, count)
	if err != nil {
		return err
	}
	newDate := date
	for i := 1; i <= 7; i++ { // Create for the next 6 days after the first
		newDate, count, err = setFutureDay(newDate, count)
		if err != nil {
			return err
		}
		err = requestFutureWeather()
		if err != nil {
			return err
		}
		err = saveFutureWeather(city, count)
		if err != nil {
			return err
		}
	}
	year, month, day, err := splitDate(date)
	if err != nil {
		return err
	}
	day += 1
	err = setDate(year, month, day)
	if err != nil {
		return err
	}
	return nil
}

func setFutureDay(date string, count string) (string, string, error) {
	_count, err := strconv.Atoi(count)
	if err != nil {
		return date, count, fmt.Errorf("ERROR: Couldn't convert Counter!\nERROR: %v", err)
	}
	year, month, day, err := splitDate(date)
	if err != nil {
		return date, count, err
	}
	_count += 1
	count = strconv.Itoa(_count)
	setDate(year, month, day)
	return fmt.Sprintf("%v-%.2v-%.2v", year, month, day), count, nil
}

func splitDate(date string) (year, month, day int, err error) {
	splitDate := strings.Split(date, "-")
	day, err = strconv.Atoi(splitDate[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("ERROR: Couldn't convert day!\nERROR: %v", err)
	}
	month, err = strconv.Atoi(splitDate[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("ERROR: Couldn't convert month!\nERROR: %v", err)
	}
	year, err = strconv.Atoi(splitDate[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("ERROR: Couldn't convert year!\nERROR: %v", err)
	}
	return year, month, day, nil
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

func getHourWeatherRecord(city string, db pg.DB) error {
	city = strings.ToLower(city)
	records, err := queryDayDatabase(city, db)
	if err != nil {
		return err
	}
	today.Hours = records
	return nil
}

func queryDayDatabase(city string, db pg.DB) ([]HourWeatherRecord, error) {
	var res []HourWeatherRecord
	year, month, day, err := splitDate(date)
	if err != nil {
		return today.Hours, fmt.Errorf("%v\nERROR: Can't query Database because of date failure!", err)
	}
	if city != "" {
		return today.Hours, fmt.Errorf("ERROR: Can't query Database because city isn't set!")
	}
	query := fmt.Sprintf("timestamp::date='%v-%.2v-%.2v' AND city='%v'", year, month, day, city)
	err = db.Model().Table("weather_records").
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
	year, month, day, err := splitDate(date)
	if err != nil {
		return fmt.Errorf("%v\nERROR: Can't query Database because of date failure!", err)
	}
	query := fmt.Sprintf("SELECT %v FROM weather_records WHERE timestamp='%v-%.2v-%.2v %.2v:00:00+00' AND city='%v'", value, year, month, day, hour, city)
	_, err = db.Query(pg.Scan(&t), query)
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

func insertCityWeatherRecordsToTable(city string, db pg.DB) error {
	path := fmt.Sprintf("resources/pg/data/%v.csv", city)
	if !pathExists(path) {
		return fmt.Errorf("ERROR: Path \"%v\" ", path)
	}
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
	_, err = db.Exec(`CREATE TEMPORARY TABLE temp_weather_records (id int NOT NULL, timestamp TIMESTAMP, source_id INT, precipitation FLOAT, pressure_msl FLOAT, sunshine FLOAT, temperature FLOAT, wind_direction INT, wind_speed FLOAT, cloud_cover FLOAT, dew_point FLOAT, relative_humidity FLOAT, visibility FLOAT, wind_gust_direction INT, wind_gust_speed FLOAT, condition VARCHAR(100), precipitation_probability FLOAT, precipitation_probability_6h FLOAT, solar FLOAT, icon VARCHAR(100), city VARCHAR(100), PRIMARY KEY(ID));`)
	if err != nil {
		return fmt.Errorf("ERROR: Couldn't create temporary table temp_weather_records\nERROR: %v", err)
	}
	var csvString []string
	var count int
	_, err = db.Query(pg.Scan(&count), "SELECT COUNT(*) FROM weather_records")
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
	_, err = db.CopyFrom(reader, `COPY temp_weather_records FROM STDIN WITH CSV`)
	if err != nil {
		return fmt.Errorf("ERROR: Couldn't copy temp_weather_records\nERROR: %v", err)
	}

	_, err = db.Exec("INSERT INTO weather_records\nSELECT * FROM temp_weather_records WHERE timestamp NOT IN (SELECT timestamp FROM weather_records)")
	if err != nil {
		return fmt.Errorf("ERROR: Couldn't insert temp_weather_records into weather_records\nERROR: %v", err)
	}
	fmt.Println("Weather Data inserted into Table!")
	return nil
}

func weatherDataExists(city string, db pg.DB) (bool, error) {
	now := time.Now()
	var timestamp string
	err := queryDatabase(&timestamp, "timestamp", now.Hour(), strings.ToLower(city), db)
	if err != nil {
		return false, err
	}
	timeString := date + fmt.Sprintf(" %.2v:00:00+00", now.Hour())
	if strings.Contains(timestamp, timeString) {
		return true, nil
	}
	return false, nil
}

/*
TODO:
-	Change log.Fatalf
	and give back error
*/
