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

	"github.com/joho/godotenv"
)

type DayWeather struct {
	Weather []HourWeather `json:"weather"`
	Source  []Source      `json:"sources"`
}

type HourWeather struct {
	Timestamp                string  `json:"timestamp"`
	SourceId                 int     `json:"source_id"`
	Precipitation            float64 `json:"precipitation"`
	PressureMsl              float64 `json:"pressure_msl"`
	Sunshine                 float64 `json:"sunshine"`
	Temperature              float64 `json:"temperature"`
	WindDirection            int     `json:"wind_direction"`
	WindSpeed                float64 `json:"wind_speed"`
	CloudCover               int     `json:"cloud_cover"`
	DewPoint                 float64 `json:"dew_point"`
	RelativeHumidity         int     `json:"relative_humidity"`
	Visibility               int     `json:"visibility"`
	WindGustDirection        int     `json:"wind_gust_direction"`
	WindGustSpeed            float64 `json:"wind_gust_speed"`
	Condition                string  `json:"condition"`
	PrecipationProbability   float64 `json:"precipation_probability"`
	PrecipationProbability6h float64 `json:"precipation_probability_6h"`
	Solar                    float64 `json:"solar"`
	Icon                     string  `json:"icon"`
}

type Source struct {
	ID              int     `json:"id"`
	DwdStationId    string  `json:"dwd_station_id"`
	ObservationType string  `json:"observation_type"`
	Latitude        float64 `json:"lat"`
	Longitude       float64 `json:"lon"`
	Height          float64 `json:"height"`
	StationName     string  `json:"station_name"`
	WmoStationId    string  `json:"wmo_station_id"`
	FirstRecord     string  `json:"first_record"`
	LastRecord      string  `json:"last_record"`
	Distance        float64 `json:"distance"`
}

type OWCity struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}

type City struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

var weatherAPIURL string
var geoAPIURL string
var geoAPIKey string
var date = "2020-04-21"
var latitude float64 = 52.5
var longitude float64 = 13.4
var selectedCity string

// Date Format year-month-day
func SetDate(year, month, day int) {
	if month <= 12 && month >= 1 && year >= 2010 && year <= time.Now().Year() && day >= 0 && day <= 31 {
		date = fmt.Sprintf("%d-02%d-02%d", year, month, day)
		reloadWEATHERAPIURL()
	} else {
		layout := "2006-01-02"
		date = time.Now().Format(layout)
		// Make a message in app for user: "Incorrect Date" (something like that)
	}
}

func SetLocationByCityName(name string, cities map[string]City) {
	cities = ReadCities(name, cities)
	if city, exists := cities[strings.ToLower(name)]; exists { // When the city exists
		fmt.Println("City:", city)
	} else { // When the city doesn't exist
		SaveCityByName(name, cities)
		SetLocationByCityName(name, cities)
	}
}

func ReadCities(name string, cities map[string]City) map[string]City {
	// Implement a check if the file below exists!
	file, err := os.Open("resources/cities.json")
	if err != nil {
		SaveCityByName(name, cities)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cities)
	return cities
}

func SaveCityByName(name string, cities map[string]City) {
	selectedCity = name
	reloadGEOAPIURL()
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
	cities[strings.ToLower(foundcity.Name)] = City{Lat: foundcity.Lat, Lon: foundcity.Lon}
	data, err := json.MarshalIndent(cities, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = saveJSONFile("resources", "cities.json", data)
	if err != nil {
		log.Fatal(err)
	}
}

func SetLocation(Latitude float64, Longitude float64) {
	if Latitude <= 54 && Latitude >= 48 && Longitude <= 14 && Longitude >= 6 {
		latitude = Latitude
		longitude = Longitude
		reloadWEATHERAPIURL()
	} else {
		// Make a message in app for user: "Not in range" (something like that)
	}
}

func SetDateAndLocation(year, month, day int, Latitude float64, Longitude float64) {
	SetDate(year, month, day)
	SetLocation(Latitude, Longitude)
}

func reloadWEATHERAPIURL() {
	u, err := url.Parse(weatherAPIURL)
	if err != nil {
		log.Fatal(err)
	}
	values := u.Query()
	values.Set("lat", strconv.FormatFloat(latitude, 'f', 2, 64))
	values.Set("lon", strconv.FormatFloat(longitude, 'f', 2, 64))
	values.Set("lat", strconv.FormatFloat(latitude, 'f', 2, 64))
	values.Set("date", date)
	u.RawQuery = values.Encode()
	weatherAPIURL = u.Redacted()
}

func reloadGEOAPIURL() {
	u, err := url.Parse(geoAPIURL)
	if err != nil {
		log.Fatal(err)
	}
	values := u.Query()
	values.Set("q", selectedCity)
	values.Set("appid", geoAPIKey)
	u.RawQuery = values.Encode()
	geoAPIURL = u.Redacted()
}

func ShowWeatherFromTime(day DayWeather, t time.Time) {
	fmt.Println(day.Weather[t.Hour()])
	// Make use of it later
}

var today DayWeather

func main() {
	var minutesRequest int
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	weatherAPIURL = os.Getenv("WAPI_URL")
	geoAPIURL = os.Getenv("GEOAPI_URL")
	geoAPIKey = os.Getenv("GEOAPI_KEY")
	minutesRequest, err = strconv.Atoi(os.Getenv("REQ_AFT_MIN"))
	if err != nil {
		panic(err)
	}
	Cities := make(map[string]City)         // Defines the Variable Cities
	SetLocationByCityName("Berlin", Cities) // Sets Location by the city name for the active Weather Request
	SetDate(2024, 8, 15)                    // Sets Date for the active Weather Request
	RequestWeather()                        // Requests Weather from the API
	ShowWeatherFromTime(today, time.Now())  // Prints weather to terminal
	SaveWeather()                           // Saves weather in weather.json
	RequestWeatherEvery(time.Duration(minutesRequest*int(time.Minute)), ShowWeather)
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

func ShowWeather(time.Time) {
	ShowWeatherFromTime(today, time.Now())
}

func RequestWeatherEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		RequestWeather()
		f(x)
	}
}

func RequestWeather() {
	response, err := http.Get(weatherAPIURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &today)
	if err != nil {
		log.Fatal(err)
	}
}

func SaveWeather() {
	data, err := json.MarshalIndent(today, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = saveJSONFile("resources", "weather.json", data)
	if err != nil {
		log.Fatal(err)
	}
}

// jq commands:
// 2: cat weather.json | jq
// 3.1: cat weather.json | jq '.weather[1]'
// 3.2: cat weather.json | jq '.weather[] | select(.wind_speed>13.3)'
// 3.3: cat weather.json | jq '.weather[] | [.temperature, .wind_speed, .relative_humidity]'
// 3.4: cat weather.json | jq '[.weather[].temperature] | add / length'
// 3.5 Sorting: cat weather.json | jq '.weather[].temperature' | sort -n
// 3.5 Finding: cat weather.json | jq '[.weather[].temperature] | max'
