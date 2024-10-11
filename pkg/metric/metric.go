package metric

import (
	"fmt"
	"log"
	"time"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
)

var m *model.Metrics

func NewMetrics(reg prometheus.Registerer) *model.Metrics {
	m = &model.Metrics{
		Temperature: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "gocast",
			Name:      "temperature",
			Help:      "Temperature of each city",
		}, []string{"location", "timestamp"}),
		Humidity: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "gocast",
			Name:      "humidity",
			Help:      "Humidity of each city",
		}, []string{"location", "timestamp"}),
		Windspeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "gocast",
			Name:      "wind_speed",
			Help:      "Wind speed of each city",
		}, []string{"location", "timestamp"}),
		Pressure: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "gocast",
			Name:      "pressure",
			Help:      "Pressure of each city",
		}, []string{"location", "timestamp"}),
	}
	reg.MustRegister(m.Temperature)
	reg.MustRegister(m.Humidity)
	reg.MustRegister(m.Windspeed)
	reg.MustRegister(m.Pressure)
	return m
}

func UpdateMetrics(record model.WeatherRecord, hour int) {
	_date := record.Hours[0].TimeStamp[:10]
	now := time.Now()
	min := now.Minute()
	sec := now.Second()
	date, err := time.Parse(time.RFC3339, fmt.Sprintf("%vT%.2v:%.2v:%.2v+02:00", _date, hour, min, sec))
	if err != nil {
		log.Print("updating metrics failed: ", err)
	}
	timestamp := date.Format(time.RFC3339)
	m.Temperature.With(prometheus.Labels{"location": record.Hours[hour].City, "timestamp": timestamp}).Set(record.Hours[hour].Temperature)
	m.Humidity.With(prometheus.Labels{"location": record.Hours[hour].City, "timestamp": timestamp}).Set(float64(record.Hours[hour].RelativeHumidity))
	m.Windspeed.With(prometheus.Labels{"location": record.Hours[hour].City, "timestamp": timestamp}).Set(record.Hours[hour].WindSpeed)
	m.Pressure.With(prometheus.Labels{"location": record.Hours[hour].City, "timestamp": timestamp}).Set(record.Hours[hour].PressureMSL)
}

func UpdateMetricsDay(record model.WeatherRecord) {
	for i := 0; i < 24; i++ {
		UpdateMetrics(record, i)
	}
}

func UpdateMetricsNow(record model.WeatherRecord) {
	now := time.Now()
	timestamp := now.Format(time.RFC3339)
	m.Temperature.With(prometheus.Labels{"location": record.Hours[now.Hour()].City, "timestamp": timestamp}).Set(record.Hours[now.Hour()].Temperature)
	m.Humidity.With(prometheus.Labels{"location": record.Hours[now.Hour()].City, "timestamp": timestamp}).Set(float64(record.Hours[now.Hour()].RelativeHumidity))
	m.Windspeed.With(prometheus.Labels{"location": record.Hours[now.Hour()].City, "timestamp": timestamp}).Set(record.Hours[now.Hour()].WindSpeed)
	m.Pressure.With(prometheus.Labels{"location": record.Hours[now.Hour()].City, "timestamp": timestamp}).Set(record.Hours[now.Hour()].PressureMSL)
}
