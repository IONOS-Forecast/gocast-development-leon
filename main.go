package main

import (
	"log"
	"time"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/db"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/utils"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
)

var (
	fdbUser,
	fdbPass,
	fdbDB,
	fdbAddress string
)
var opts options

type options struct {
	WeatherAPIURL  string `v:"wapiurl" long:"weather-api-url" env:"WEATHER_API_URL" description:"URL to interact with Weather provider"`
	GeoAPIURL      string `v:"gapiurl" long:"geo-api-url" env:"GEO_API_URL" description:"URL to interact with GEO provider"`
	GeoAPIKey      string `v:"gapikey" long:"geo-api-key" env:"GEO_API_KEY" description:"KEY to interact with GEO provider"`
	MinutesRequest string `v:"reqmin" long:"req-aft-min" env:"REQ_AFT_MIN" description:"Minutes until the next request to the Weather provicer is made"`
	FdbUser        string `v:"fdb-u" long:"fdb-user" env:"FDB_USER" description:"The user that connects to the forecast-database"`
	FdbPassword    string `v:"fdb-p" long:"fdb_password" env:"FDB_PASSWORD" description:"The password to the user for connecting to the forecast-database"`
	FdbDatabase    string `v:"fdb-db" long:"fdb_database" env:"FDB_DATABASE" description:"The database that the user connects to"`
	FdbAddress     string `v:"fdb-addr" long:"fdb_address" env:"FDB_ADDRESS" description:"The address of the database"`
}

func main() {
	utils.MakeCities()
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	err = godotenv.Load(".env.dev")
	if err != nil {
		log.Fatal(err)
	}
	_, err = flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	fdbUser = opts.FdbUser
	fdbPass = opts.FdbPassword
	fdbDB = opts.FdbDatabase
	fdbAddress = opts.FdbAddress
	utils.SetWeatherAPIURL(opts.WeatherAPIURL)
	utils.SetGeocodingAPIURL(opts.GeoAPIURL)
	utils.SetGeocodingAPIKey(opts.GeoAPIKey)
	now := time.Now()
	date, err := utils.SetDate(now.Year(), int(now.Month()), now.Day())
	if err != nil {
		log.Print(err)
	}
	cityName := utils.GetCityName()
	database, err := db.NewPG(fdbUser, fdbPass, fdbDB, fdbAddress)
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

}
