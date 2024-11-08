package main

import (
	"log"
	"net/http"
	"strconv"
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
	reqAftMin := utils.Options.MinutesRequest
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
	cityDB := db.NewCityDB(database, "")
	cities, err := cityDB.GetCities()
	if err != nil {
		log.Print(err)
	}
	weatherMapDB := db.NewWeatherMapDB(cities, utils.Options.GeoAPIURL, utils.Options.GeoAPIKey)
	// Getting Default Cities if needed
	cityName, err = database.SetLocationByCityName("Berlin")
	if err != nil {
		log.Print(err)
	}
	record, err := database.GetWeatherRecords(cityName, date)
	if err != nil {
		log.Print(err)
	}
	metric.UpdateMetrics(record, time.Now().Hour())
	cityName, err = database.SetLocationByCityName("München")
	if err != nil {
		log.Print(err)
	}
	record, err = database.GetWeatherRecords(cityName, date)
	if err != nil {
		log.Print(err)
	}
	metric.UpdateMetrics(record, time.Now().Hour())
	cityName, err = database.SetLocationByCityName("Hamburg")
	if err != nil {
		log.Print(err)
	}
	record, err = database.GetWeatherRecords(cityName, date)
	if err != nil {
		log.Print(err)
	}
	metric.UpdateMetrics(record, time.Now().Hour())
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
	handler := metric.NewHandler(database, cityDB, weatherMapDB)
	mux := http.NewServeMux()
	mux.HandleFunc("/direct", handler.Get)
	mux.HandleFunc("/error", handler.Error)

	cityName, err = database.SetLocationByCityName("Berlin")
	if err != nil {
		log.Print(err)
	}
	record, err = database.GetWeatherRecord(cityName, time.Now().Format("2006-01-02"))
	metric.UpdateMetricsNow(record)

	pMux := http.NewServeMux()
	pMux.Handle("/metrics", promHandler)
	go func() {
		minutes, err := strconv.Atoi(reqAftMin)
		if err != nil {
			log.Fatal(err)
		}
		c := time.NewTicker(time.Duration(minutes) * time.Minute)
		for {
			select {
			case <-c.C:
				record, err = database.GetWeatherRecord(cityName, time.Now().Format("2006-01-02"))
				if err != nil {
					log.Print(err)
				}
				metric.UpdateMetricsNow(record)
			}
		}
	}()
	go func() {
		log.Fatal(http.ListenAndServe(":8080", mux))
	}()
	go func() {
		log.Fatal(http.ListenAndServe(":8081", pMux))
	}()
	select {}
}
