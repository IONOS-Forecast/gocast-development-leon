package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
)

var weatherAPIURL string
var geoAPIURL string
var geoAPIKey string
var latitude = 52.5170365
var longitude = 13.3888599
var date string
var cityName string
var today WeatherRecord
var notToday WeatherRecord
var opts options

type options struct {
	WeatherAPIURL  string `v:"wapiurl" long:"weather-api-url" env:"WEATHER_API_URL" description:"URL to interact with Weather provider"`
	GeoAPIURL      string `v:"gapiurl" long:"geo-api-url" env:"GEO_API_URL" description:"URL to interact with GEO provider"`
	GeoAPIKEY      string `v:"gapikey" long:"geo-api-key" env:"GEO_API_KEY" description:"KEY to interact with GEO provider"`
	MinutesRequest string `v:"reqmin" long:"req-aft-min" env:"REQ_AFT_MIN" description:"Minutes until the next request to the Weather provicer is made"`
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
	CloudCover                 int     `json:"cloud_cover"`
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
	//minutesRequest, err := strconv.Atoi(opts.MinutesRequest)
	if err != nil {
		log.Fatal(err)
	}
	setDateAndLocationByCityName(2024, 9, 11, "Berlin", cities)
	requestWeather()
	fmt.Println(today.Hours[time.Now().Hour()])
	saveFutureWeatherInFile(cityName, date)
	//requestWeatherEvery(time.Duration(minutesRequest*int(time.Minute)), showWeather)
}

func showWeather(time.Time) {
	fmt.Println(today.Hours[time.Now().Hour()])
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
	newDate, count := setFutureDay(date, count) // Create the first next day
	requestFutureWeather()
	saveFutureWeather(city, count)
	for i := 1; i <= 6; i++ { // Create for the next 6 days after the first
		newDate, count = setFutureDay(newDate, count)
		requestFutureWeather()
		saveFutureWeather(city, count)
	}
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

/*
TODO:
-	Change log.Fatalf
	and give back error
-	Retrieve Data From DB
*/
