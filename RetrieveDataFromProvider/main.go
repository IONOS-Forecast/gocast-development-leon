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
	"time"
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

var URL = "https://api.brightsky.dev/weather?lat=52&lon=7.6&date=2020-04-21"
var _url = "https://api.brightsky.dev/weather?"
var date = "2020-04-21"
var latitude float64 = 52
var longitude float64 = 7.6

// Date Format year-month-day
func SetDate(year, month, day int) {
	if month <= 12 && month >= 1 && year >= 2010 && year <= time.Now().Year() && day >= 0 && day <= 31 {
		_day := fmt.Sprintf("%02d", day)
		_month := fmt.Sprintf("%02d", month)
		fmt.Println(_day)
		fmt.Println(_month)
		fmt.Println(strconv.Itoa(year))
		date = strconv.Itoa(year) + "-" + _month + "-" + _day
		fmt.Println(date)
		reloadURL()
	} else {
		date = strconv.Itoa(time.Now().Year()) + "-" + fmt.Sprintf("%02d", int(time.Now().Month())) + "-" + strconv.Itoa(time.Now().Day())
		fmt.Println(date)
		// Make a message in app for user: "Incorrect Date" (something like that)
	}
}

func SetLocation(Latitude float64, Longitude float64) {
	if Latitude <= 54 && Latitude >= 48 && Longitude <= 14 && Longitude >= 6 {
		latitude = Latitude
		longitude = Longitude
		reloadURL()
	} else {
		// Make a message in app for user: "Not in range" (something like that)
	}
}

func SetDateAndLocation(year, month, day int, Latitude float64, Longitude float64) {
	SetDate(year, month, day)
	SetLocation(Latitude, Longitude)
}

func reloadURL() {
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}
	values := u.Query()
	values.Set("lat", strconv.FormatFloat(latitude, 'f', 2, 64))
	values.Set("lon", strconv.FormatFloat(longitude, 'f', 2, 64))
	values.Set("date", date)
	fmt.Println(u.Query())
	fmt.Println(u.Redacted())
	u.RawQuery = values.Encode()
	URL = u.Redacted()
	fmt.Println(u.Query())
	fmt.Println(u.Redacted())
}

func ShowWeatherFromTime(day DayWeather, t time.Time) {
	fmt.Println(day.Weather[t.Hour()])
}

func main() {
	SetDateAndLocation(2025, 10, 5, 54, 14)

	var today DayWeather
	response, err := http.Get(URL)
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
	/*mdata, err := json.MarshalIndent(today, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = saveJSONFile("weather.json", mdata)
	if err != nil {
		log.Fatal(err)
	}*/
}

func saveJSONFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	fmt.Println("File created!")
	return nil
}

// jq commands:
// 2: cat weather.json | jq
// 3.1: cat weather.json | jq '.weather[1]'
// 3.2: cat weather.json | jq '.weather[] | select(.wind_speed>13.3)'
// 3.3: cat weather.json | jq '.weather[] | [.temperature, .wind_speed, .relative_humidity]'
// 3.4: cat weather.json | jq '[.weather[].temperature] | add / length'
// 3.5 Sorting: cat weather.json | jq '.weather[].temperature' | sort -n
// 3.5 Finding: cat weather.json | jq '[.weather[].temperature] | max'
