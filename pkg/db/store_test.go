package db

import (
	"os"
	"testing"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/testutil"
	"github.com/go-pg/pg/v10"
	_ "github.com/lib/pq"
)

func TestDB(t *testing.T) {
	conn := testutil.TestDB(t, dbOpts, "test")
	records := make([]WeatherRecord, 0)
	err := conn.Model(&records).Select(&records)
	if err != nil {
		t.Errorf("failed to get weather records: %v\n", err)
	}

	t.Logf("got weather records: %v\n", records)
}

type WeatherRecord struct {
	Id                         int     `json:"id" pg:"id"`
	Timestamp                  string  `json:"timestamp" pg:"timestamp"`
	SourceID                   int     `json:"source_id" pg:"source_id"`
	Precipitaiton              float64 `json:"precipitation" pg:"precipitation"`
	PressureMSL                float64 `json:"pressure_msl" pg:"pressure_msl"`
	Sunshine                   float64 `json:"sunshine" pg:"sunshine"`
	Temperature                float64 `json:"temperature" pg:"temperature"`
	WindDirection              int     `json:"wind_direction" pg:"wind_direction"`
	WindSpeed                  float64 `json:"wind_speed" pg:"wind_speed"`
	CloudCover                 int     `json:"cloud_cover" pg:"cloud_cover"`
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
	City                       string  `json:"city" pg:"city"`
}

var dbOpts pg.Options

func TestMain(m *testing.M) {
	opts := pg.Options{
		Addr:     "127.0.0.1:5432",
		User:     "forecast",
		Password: "forecast",
		Database: "forecast_test",
	}
	os.Exit(func() int {
		container := testutil.StartPostgresContainer(opts)
		opts.Addr = container.Addr
		dbOpts = opts
		defer container.Shutdown()
		return m.Run()
	}())
}
