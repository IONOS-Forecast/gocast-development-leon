package main

import (
	"log"
	"net/http"
	"time"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/db"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/metric"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/utils"
)

var (
	fdbUser,
	fdbPass,
	fdbDB,
	fdbAddress string
)

func main() {
	utils.MakeCities()
	utils.DotEnvironment()
	utils.SetWeatherAPIURL(utils.Options.WeatherAPIURL)
	utils.SetGeocodingAPIURL(utils.Options.GeoAPIURL)
	utils.SetGeocodingAPIKey(utils.Options.GeoAPIKey)
	now := time.Now()
	date, err := utils.SetDate(now.Year(), int(now.Month()), now.Day())
	if err != nil {
		log.Print(err)
	}
	cityName := utils.GetCityName()
	database, err := db.NewPG(utils.FdbUser, utils.FdbPass, utils.FdbDB, utils.FdbAddress)
	if err != nil {
		log.Print(err)
	}
	defer database.Close()
	// Getting Default Cities if needed
	cityName, err = database.SetLocationByCityName("Berlin")
	if err != nil {
		log.Print(err)
	}
	_, err = database.GetWeatherRecord(cityName, date)
	if err != nil {
		log.Print(err)
	}
	cityName, err = database.SetLocationByCityName("München")
	if err != nil {
		log.Print(err)
	}
	_, err = database.GetWeatherRecord(cityName, date)
	if err != nil {
		log.Print(err)
	}
	cityName, err = database.SetLocationByCityName("Hamburg")
	if err != nil {
		log.Print(err)
	}
	_, err = database.GetWeatherRecord(cityName, date)
	if err != nil {
		log.Print(err)
	}
	/* TODO: ADD OR REMOVE LATER
	Error Examples and MinutesRequest
	fmt.Println(":----------------------------------------------------------------------------:")
	fmt.Println("\t\t\tERROR EXAMPLES BEGIN HERE")
	fmt.Println(":----------------------------------------------------------------------------:")
	date, err = utils.SetDate(2024, 2, 30)
	if err != nil {
		log.Print(err)
	}
	cityName, err = utils.SetLocationByCityName("Afrika", cities)
	if err != nil {
		log.Print(err)
	}
	fmt.Println(":----------------------------------------------------------------------------:")
	fmt.Println("\t\t\t ERROR EXAMPLES END HERE")
	fmt.Println(":----------------------------------------------------------------------------:")
	minutesRequest, err := strconv.Atoi(opts.MinutesRequest)
	if err != nil {
		log.Fatal(err)
	}
	requestWeatherEvery(time.Duration(minutesRequest*int(time.Minute)), showWeather)*/
	var handler metric.Handler
	http.HandleFunc("/GET", handler.Get)
	http.ListenAndServe(":3333", nil)
}
