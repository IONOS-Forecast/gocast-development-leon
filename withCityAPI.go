package main

import (

	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"io"
) 

// start: variables initialization

type WeatherByHour struct {
	TimeStamp                  time.Time  	 `json:"timestamp"`
	SourceId			int      `json:"source_id"`
	Precipitation			float64  `json:"precipitation"`
	PressureMsl			float64  `json:"pressure_msl"`
	Sunshine			float64  `json:"sunshine"`
	Temperature			float64  `json:"temperature"`
	WindDirection			int      `json:"wind_direction"`
	WindSpeed			float64  `json:"wind_speed"`
	CloudCover			int      `json:"cloud_cover"`
	DewPoint			float64  `json:"dew_point"`
	RelativeHumidity		int      `json:"relative_humidity"`
	Visibility			int      `json:"visibility"`
	WindGustDirection		int      `json:"wind_gust_direction"`
	WindGustSpeed			float64  `json:"wind_gust_speed"`
	Condition			string   `json:"condition"`
	PrecipitationProbability	float64 `json:"precipitation_probability"`
	PrecipitationProbability6h	float64 `json:"precipitation_probability_6h"`
	Solar				float64  `json:"solar"`
	Icon				string   `json:"icon"`
}


type WeatherbyDay struct {
	WeatherByHours	[]WeatherByHour		`json:"weather"`
}

type Cityinfo struct {
	Lat				float64	 `json:lat`
	Lon				float64	 `json:lon`
}

var citynumbers []Cityinfo
var weather WeatherbyDay

// end: variables initialization

//function that takes city name as input and saves latitude/longitude in "citynumbers"
func GetLatLong (city string){
	
	url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=2cab1704c3ad14814b44b266c13346a8",city)		//please dont get me banned
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("error with sending the rquest to get the longitude/latitude")
		return
	}
	
	body, err := io.ReadAll(resp.Body)
	
	if err != nil {
		fmt.Println("error with reading response body from longitude/latitude")
		return
	}
	
	err = json.Unmarshal(body, &citynumbers)
}

//function for checking the input. if theyre correct, it sends the variables to the next function
func CheckArguments(year int, month int, day int, hour int){

switch{

case year < 2010 || year > 2024:
fmt.Println("check year input")
return

case month < 1 || month > 12:
fmt.Println("check month input")
return

case day < 1 || day > 31:
fmt.Println("check day input")
return

default:
SendRequest(year, month, day, hour)

}
}

//function that takes input from "CheckArguments" function and prints weather information
func SendRequest(year int, month int, day int, hour int) {
	
	url := fmt.Sprintf("https://api.brightsky.dev/weather?lat=%f&lon=%f&date=%.4d-%.2d-%.2d",citynumbers[0].Lat,citynumbers[0].Lon,year,month,day)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("error with sending the rquest for weather")
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error with reading response body for weather")
		return
	}

	err = json.Unmarshal(body, &weather)

	if err != nil {
		fmt.Println(err)
		return
	}
	
	PrintWeather(year, month, day, hour)

}

//function that takes input from "SendRequest" and prints weather data
func PrintWeather(year int, month int, day int, hour int){

fmt.Printf("time:					%.16v\n",weather.WeatherByHours[hour].TimeStamp)
fmt.Printf("condition:				%s\n",weather.WeatherByHours[hour].Condition)
fmt.Printf("temperature:				%.1f\n",weather.WeatherByHours[hour].Temperature)
fmt.Printf("wind speed:				%.1f\n",weather.WeatherByHours[hour].WindSpeed)
fmt.Printf("wind direction:				%d\n",weather.WeatherByHours[hour].WindDirection)
fmt.Printf("wind gust speed:			%.1f\n",weather.WeatherByHours[hour].WindGustSpeed)
fmt.Printf("wind gust direction:			%d\n",weather.WeatherByHours[hour].WindGustDirection)
fmt.Printf("relative humidity:			%d\n",weather.WeatherByHours[hour].RelativeHumidity)
fmt.Printf("dew point:				%.1f\n",weather.WeatherByHours[hour].DewPoint)
fmt.Printf("precipitation probability:		%.1f\n",weather.WeatherByHours[hour].PrecipitationProbability)
fmt.Printf("precipitation probability 6h:		%.1f\n",weather.WeatherByHours[hour].PrecipitationProbability6h)
fmt.Printf("visibility:				%d\n",weather.WeatherByHours[hour].Visibility)
fmt.Printf("pressure in MSL:			%.1f\n",weather.WeatherByHours[hour].PressureMsl)
fmt.Printf("cloud cover:				%d\n",weather.WeatherByHours[hour].CloudCover)
fmt.Printf("sunshine:				%.0f\n",weather.WeatherByHours[hour].Sunshine)
fmt.Printf("solar:					%.3f\n",weather.WeatherByHours[hour].Solar)
fmt.Printf("general:				%s\n",weather.WeatherByHours[hour].Icon)
fmt.Printf("precipitation:				%.1f\n",weather.WeatherByHours[hour].Precipitation)

}

func main() {

	var city string = "Muenchen"	// city name
	var year int	= 2024		// year
	var month int 	= 8		// month
	var day int 	= 14		// day
	var hour int	= 2		// hour
	
	GetLatLong(city)
	CheckArguments(year, month, day, hour)
}

