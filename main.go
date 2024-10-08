package main

import (
	"log"
	"net/http"
	"time"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/db"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/metric"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	fdbUser,
	fdbPass,
	fdbDB,
	fdbAddress string
)

func main() {
	reg := prometheus.NewRegistry()
	metric.NewMetrics(reg)
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
	record, err := database.GetWeatherRecord(cityName, date)
	if err != nil {
		log.Print(err)
	}
	metric.RegisterMetrics(record, 0)
	cityName, err = database.SetLocationByCityName("München")
	if err != nil {
		log.Print(err)
	}
	record, err = database.GetWeatherRecord(cityName, date)
	if err != nil {
		log.Print(err)
	}
	metric.RegisterMetrics(record, 0)
	cityName, err = database.SetLocationByCityName("Hamburg")
	if err != nil {
		log.Print(err)
	}
	record, err = database.GetWeatherRecord(cityName, date)
	if err != nil {
		log.Print(err)
	}
	metric.RegisterMetrics(record, 0)
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
	*/
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	var handler metric.Handler
	mux := http.NewServeMux()
	mux.HandleFunc("/GET", handler.Get)
	mux.HandleFunc("/error", handler.Error)

	pMux := http.NewServeMux()
	pMux.Handle("/metrics", promHandler)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", mux))
	}()
	go func() {
		log.Fatal(http.ListenAndServe(":8081", pMux))
	}()
	select {}
}
