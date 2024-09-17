package model

type WeatherRecord struct {
	Hours []HourWeatherRecord `json:"weather"`
}

type HourWeatherRecord struct {
	TimeStamp                  string  `json:"timestamp" pg:"timestamp"`
	SourceID                   int     `json:"source_id" pg:"source_id"`
	Precipitation              float64 `json:"precipitation" pg:"precipitation"`
	PressureMSL                float64 `json:"pressure_msl" pg:"pressure_msl"`
	Sunshine                   float64 `json:"sunshine" pg:"sunshine"`
	Temperature                float64 `json:"temperature" pg:"temperature"`
	WindDirection              int     `json:"wind_direction" pg:"wind_direction"`
	WindSpeed                  float64 `json:"wind_speed" pg:"wind_speed"`
	CloudCover                 float64 `json:"cloud_cover" pg:"cloud_cover"`
	DewPoint                   float64 `json:"dew_point" pg:"dew_point"`
	RelativeHumidity           int     `json:"relative_humidity" pg:"relative_humidity"`
	Visibility                 int     `json:"visibility" pg:"visibility"`
	WindGustDirection          int     `json:"wind_gust_direction" pg:"wind_gust_direction"`
	WindGustSpeed              float64 `json:"wind_gust_speed" pg:"wind_gust_speed"`
	Condition                  string  `json:"condition" pg:"condition"`
	PrecipitationProbability   float64 `json:"precipitation_probability" pg:"precipitation_probability"`
	PrecipitationProbability6h float64 `json:"precipitation_probability_6h" pg:"precipitation_probability_6h"`
	Solar                      float64 `json:"solar" pg:"solar"`
	Icon                       string  `json:"icon" pg:"icon"`
}

type OWCity struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Country   string  `json:"country"`
}

type City struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
