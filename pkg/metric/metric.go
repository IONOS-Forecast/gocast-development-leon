package metric

import (
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

func RegisterMetrics(record model.WeatherRecord, hour int) {
	if len(record.Hours) != 0 {
		m.Temperature.With(prometheus.Labels{"location": record.Hours[hour].City, "timestamp": record.Hours[hour].TimeStamp}).Set(record.Hours[hour].Temperature)
		m.Humidity.With(prometheus.Labels{"location": record.Hours[hour].City, "timestamp": record.Hours[hour].TimeStamp}).Set(float64(record.Hours[hour].RelativeHumidity))
		m.Windspeed.With(prometheus.Labels{"location": record.Hours[hour].City, "timestamp": record.Hours[hour].TimeStamp}).Set(record.Hours[hour].WindSpeed)
		m.Pressure.With(prometheus.Labels{"location": record.Hours[hour].City, "timestamp": record.Hours[hour].TimeStamp}).Set(record.Hours[hour].PressureMSL)
	}
}

func getMetrics() *model.Metrics {
	if m != nil {
		return m
	}
	return nil
}
