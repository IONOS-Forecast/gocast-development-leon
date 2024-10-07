package model

type WeatherRecord struct {
	Hours []HourWeatherRecord `json:"weather"`
}

type HourWeatherRecord struct {
	tableName                  struct{} `pg:"weather_records"`
	TimeStamp                  string   `json:"timestamp" pg:"timestamp"`
	SourceID                   int      `json:"source_id" pg:"source_id"`
	Precipitation              float64  `json:"precipitation" pg:"precipitation"`
	PressureMSL                float64  `json:"pressure_msl" pg:"pressure_msl"`
	Sunshine                   float64  `json:"sunshine" pg:"sunshine"`
	Temperature                float64  `json:"temperature" pg:"temperature"`
	WindDirection              int      `json:"wind_direction" pg:"wind_direction"`
	WindSpeed                  float64  `json:"wind_speed" pg:"wind_speed"`
	CloudCover                 float64  `json:"cloud_cover" pg:"cloud_cover"`
	DewPoint                   float64  `json:"dew_point" pg:"dew_point"`
	RelativeHumidity           int      `json:"relative_humidity" pg:"relative_humidity"`
	Visibility                 int      `json:"visibility" pg:"visibility"`
	WindGustDirection          int      `json:"wind_gust_direction" pg:"wind_gust_direction"`
	WindGustSpeed              float64  `json:"wind_gust_speed" pg:"wind_gust_speed"`
	Condition                  string   `json:"condition" pg:"condition"`
	PrecipitationProbability   float64  `json:"precipitation_probability" pg:"precipitation_probability"`
	PrecipitationProbability6h float64  `json:"precipitation_probability_6h" pg:"precipitation_probability_6h"`
	Solar                      float64  `json:"solar" pg:"solar"`
	Icon                       string   `json:"icon" pg:"icon"`
	City                       string   `pg:"city"`
}

type OWCity struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Country   string  `json:"country"`
}

type City struct {
	tableName struct{} `pg:"cities"`
	Name      string   `pg:"name"`
	Lat       float64  `json:"lat"`
	Lon       float64  `json:"lon"`
}

type Options struct {
	WeatherAPIURL  string `v:"wapiurl" long:"weather-api-url" env:"WEATHER_API_URL" description:"URL to interact with Weather provider"`
	GeoAPIURL      string `v:"gapiurl" long:"geo-api-url" env:"GEO_API_URL" description:"URL to interact with GEO provider"`
	GeoAPIKey      string `v:"gapikey" long:"geo-api-key" env:"GEO_API_KEY" description:"KEY to interact with GEO provider"`
	MinutesRequest string `v:"reqmin" long:"req-aft-min" env:"REQ_AFT_MIN" description:"Minutes until the next request to the Weather provicer is made"`
	FdbUser        string `v:"fdb-u" long:"fdb-user" env:"FDB_USER" description:"The user that connects to the forecast-database"`
	FdbPassword    string `v:"fdb-p" long:"fdb_password" env:"FDB_PASSWORD" description:"The password to the user for connecting to the forecast-database"`
	FdbDatabase    string `v:"fdb-db" long:"fdb_database" env:"FDB_DATABASE" description:"The database that the user connects to"`
	FdbAddress     string `v:"fdb-addr" long:"fdb_address" env:"FDB_ADDRESS" description:"The address of the database"`
}
