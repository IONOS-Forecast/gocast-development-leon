package main

import (

	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"io"
) 

type WeatherByHour struct {
	TimeStamp                  time.Time  	 `json:"timestamp"`
	SourceId			int      `json:"source_id"`
	Precipitation			float64  `json:"precipitation"`
	PressureMsl			float64  `json:"pressure_msl"`
	SunShine			float64  `json:"sunshine"`
	Temprature			int      `json:"tempreture"`
	WindDirection			int      `json:"wind_direction"`
	WindSpeed			float64  `json:"wind_speed"`
	CloudCover			int      `json:"cloud_cover"`
	DewPoint			float64  `json:"dew_point"`
	RelaiveHumidity			int      `json:"relative_humidity"`
	Visibility			int      `json:"visibility"`
	WindGustDirection		int      `json:"wind_gust_direction"`
	WindGustSpeed			float64  `json:"wind_gust_speed"`
	Condition			string   `json:"condition"`
	PrecipitationProbability	*float64 `json:"precipitation_probability"`
	PrecipitationProbability6h	*float64 `json:"precipitation_probability_6h"`
	Solar				float64  `json:"solar"`
	Icon				string   `json:"icon"`
}

type Sources struct {
	Id				int       `json:"id"`
	DwdStationId			string    `json:"dwd_station_id"`
	ObservationType			string    `json:"observation_type"`
	Lat				float64   `json:"lat"`
	Lon				float64   `json:"lon"`
	Height				float64   `json:"height"`
	StationName			string    `json:"station_name"`
	WmoStationId			string    `json:"wmo_station_id"`
	FirstRecord			time.Time `json:"first_record"`
	LastRecord			time.Time `json:"last_record"`
	Distance			float64   `json:"distance"`
}

type WeatherbyDay struct {
	WeatherByHours	[]WeatherByHour		`json:"weather"`
	Source		[]Sources		`json:"sources"`
}

type Cityinfo struct {
	Lat				float64	 `json:lat`
	Lon				float64	 `json:lon`
}


func CheckArguments(year int, month int, day int, hour int, city string){

switch{

case year < 2010 || year > 2024:

fmt.Println("check year input")
return

case month < 1 || month > 12:
fmt.Println("check month input")
return

case day < 1 || day > 31:
fmt.Println("check day input")

//case lat < -90 || lat > 90:
//fmt.Println("check latitude input")

//case lon < -180 || lat > 180:
//fmt.Println("check longitude input")

default:
SendRequest(year, month, day, hour, city)
}
}

func SendRequest(year int, month int, day int, hour int, city string) {


	var citynumbers []Cityinfo
	
	url1 := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=2cab1704c3ad14814b44b266c13346a8",city)
	resp1, err1 := http.Get(url1)
	
	if err1 != nil {
		fmt.Println("error with sending the rquest")
		return
	}
	
	body1, err1 := io.ReadAll(resp1.Body)
	
	if err1 != nil {
		fmt.Println("error with reading response body")
		return
	}
	
	err1 = json.Unmarshal(body1, &citynumbers)

	var weather WeatherbyDay
	
	url := fmt.Sprintf("https://api.brightsky.dev/weather?lat=%f&lon=%f&date=%.4d-%.2d-%.2d",citynumbers[0].Lat,citynumbers[0].Lon,year,month,day)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("error with sending the rquest")
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error with reading response body")
		return
	}

	err = json.Unmarshal(body, &weather)

	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println(weather.WeatherByHours[hour])

}

func main() {

	//year ,month, day, latitude, longitude, hour. example: 2020,4,21,52,7.6,0
	CheckArguments(2024,8,14,2, "koethen")
}

