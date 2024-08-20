package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// start: variables initialization

type WeatherByHour struct {
	TimeStamp                  time.Time `json:"timestamp"`
	SourceId                   int       `json:"source_id"`
	Precipitation              float64   `json:"precipitation"`
	PressureMsl                float64   `json:"pressure_msl"`
	Sunshine                   float64   `json:"sunshine"`
	Temperature                float64   `json:"temperature"`
	WindDirection              int       `json:"wind_direction"`
	WindSpeed                  float64   `json:"wind_speed"`
	CloudCover                 int       `json:"cloud_cover"`
	DewPoint                   float64   `json:"dew_point"`
	RelativeHumidity           int       `json:"relative_humidity"`
	Visibility                 int       `json:"visibility"`
	WindGustDirection          int       `json:"wind_gust_direction"`
	WindGustSpeed              float64   `json:"wind_gust_speed"`
	Condition                  string    `json:"condition"`
	PrecipitationProbability   float64   `json:"precipitation_probability"`
	PrecipitationProbability6h float64   `json:"precipitation_probability_6h"`
	Solar                      float64   `json:"solar"`
	Icon                       string    `json:"icon"`
}

type WeatherbyDay struct {
	WeatherByHours []WeatherByHour `json:"weather"`
}

type Cityinfo struct {
	Lat float64 `json:lat`
	Lon float64 `json:lon`
}

var citynumbers []Cityinfo
var weather WeatherbyDay

const CityAPIKey string = "2cab1704c3ad14814b44b266c13346a8"

var CityAPIUrl string = "http://api.openweathermap.org/geo/1.0/direct"
var WeatherAPIUrl string = "https://api.brightsky.dev/weather"

// end: variables initialization

// function that takes city as input, use it in API to get lat/long of the city and saves it in the variable citynumbers
func GetLatLong(city string) {

	// start trying

	result, err := url.Parse(CityAPIUrl)
	ErrorHandling(err)

	values := result.Query()
	values.Set("q", city)
	values.Set("appid", CityAPIKey)
	values.Set("limit", "1")

	result.RawQuery = values.Encode()

	CityAPIUrl = result.String()

	resp, err := http.Get(CityAPIUrl)
	ErrorHandling(err)

	body, err := io.ReadAll(resp.Body)
	ErrorHandling(err)

	err = json.Unmarshal(body, &citynumbers)
	ErrorHandling(err)

}

// a simple function for error handling
func ErrorHandling(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// function for checking date input (yyyy-mm-dd) and returns a formatted date
func CheckDate(year int, month int, day int) string {

	date, err := time.Parse("2006-01-02", fmt.Sprintf("%.4d-%.2d-%.2d", year, month, day))
	ErrorHandling(err)

	dateString := date.Format("2006-01-02")

	if year < 2010 {
		log.Fatal("year input year should be 2010 or more")
	}

	return dateString

}

// function that takes (year, month, day, hour) as input, uses the lat/long from function GetLatLong, uses an API to get the weather data for that time and then calls function PrintWeather to print it
func SendRequest(year int, month int, day int, hour int) {

	result, err := url.Parse(WeatherAPIUrl)
	ErrorHandling(err)

	queries := result.Query()
	queries.Add("lat", strconv.FormatFloat(citynumbers[0].Lat, 'f', -1, 64))
	queries.Add("lon", strconv.FormatFloat(citynumbers[0].Lon, 'f', -1, 64))
	queries.Add("date", CheckDate(year, month, day))

	result.RawQuery = queries.Encode()

	WeatherAPIUrl = result.String()

	resp, err := http.Get(WeatherAPIUrl)
	ErrorHandling(err)

	body, err := io.ReadAll(resp.Body)
	ErrorHandling(err)

	err = json.Unmarshal(body, &weather)
	ErrorHandling(err)

	//check if the weather data are available for that day/ hour. Because the weather API predits exactly 10 days ahead
	if len(weather.WeatherByHours)-1 < hour {
		log.Fatal("weather data hasn't been imported from API")
	}

	PrintWeather(hour)
}

// function that takes input from "SendRequest" and prints weather data
func PrintWeather(hour int) {

	fmt.Printf("time:					%.16v\n", weather.WeatherByHours[hour].TimeStamp)
	fmt.Printf("condition:				%s\n", weather.WeatherByHours[hour].Condition)
	fmt.Printf("temperature:				%.1f\n", weather.WeatherByHours[hour].Temperature)
	fmt.Printf("wind speed:				%.1f\n", weather.WeatherByHours[hour].WindSpeed)
	fmt.Printf("wind direction:				%d\n", weather.WeatherByHours[hour].WindDirection)
	fmt.Printf("wind gust speed:			%.1f\n", weather.WeatherByHours[hour].WindGustSpeed)
	fmt.Printf("wind gust direction:			%d\n", weather.WeatherByHours[hour].WindGustDirection)
	fmt.Printf("relative humidity:			%d\n", weather.WeatherByHours[hour].RelativeHumidity)
	fmt.Printf("dew point:				%.1f\n", weather.WeatherByHours[hour].DewPoint)
	fmt.Printf("precipitation probability:		%.1f\n", weather.WeatherByHours[hour].PrecipitationProbability)
	fmt.Printf("precipitation probability 6h:		%.1f\n", weather.WeatherByHours[hour].PrecipitationProbability6h)
	fmt.Printf("visibility:				%d\n", weather.WeatherByHours[hour].Visibility)
	fmt.Printf("pressure in MSL:			%.1f\n", weather.WeatherByHours[hour].PressureMsl)
	fmt.Printf("cloud cover:				%d\n", weather.WeatherByHours[hour].CloudCover)
	fmt.Printf("sunshine:				%.0f\n", weather.WeatherByHours[hour].Sunshine)
	fmt.Printf("solar:					%.3f\n", weather.WeatherByHours[hour].Solar)
	fmt.Printf("general:				%s\n", weather.WeatherByHours[hour].Icon)
	fmt.Printf("precipitation:				%.1f\n", weather.WeatherByHours[hour].Precipitation)

}

func main() {

	var city string = "Muenchen" // city name
	var year int = 03750         // year
	var month int = 8            // month
	var day int = 25             // day
	var hour int = 15            // hour

	GetLatLong(city)
	SendRequest(year, month, day, hour)
}
