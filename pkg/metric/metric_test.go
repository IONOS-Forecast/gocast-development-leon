package metric_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/metric"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
)

func TestUpdateMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()
	metric.NewMetrics(reg)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	pMux := http.NewServeMux()
	pMux.Handle("/metrics", promHandler)
	go func() {
		log.Fatal(http.ListenAndServe(":8081", pMux))
	}()
	testTimestamp := fmt.Sprintf("2024-10-02T%.2v:00:00+00:00", time.Now().Hour())
	expectedTemperature := 19.6
	expectedWindSpeed := 14.0
	expectedHumidity := 65
	expectedPressure := 996.8
	var hours []model.HourWeatherRecord
	for i := 0; i < 24; i++ {
		hours = append(hours, model.HourWeatherRecord{TimeStamp: testTimestamp, SourceID: 303712, Precipitation: 0, PressureMSL: expectedPressure,
			Sunshine: 60, Temperature: expectedTemperature, WindDirection: 190, WindSpeed: expectedWindSpeed, CloudCover: 100, DewPoint: 12.8, RelativeHumidity: expectedHumidity, Visibility: 61300,
			WindGustDirection: 200, WindGustSpeed: 29.5, Condition: "dry", PrecipitationProbability: 0, PrecipitationProbability6h: 0, Solar: 0, Icon: "cloudy", City: "berlin"})
	}
	record := model.WeatherRecord{Hours: hours}
	metric.UpdateMetrics(record)
	temperature, err := getMetrics("gocast_temperature", testTimestamp)
	if err != nil {
		t.Error(err)
	}
	if temperature != expectedTemperature {
		t.Errorf("metrics returned wrong temperature: got \"%v\" want \"%v\"", temperature, expectedTemperature)
	}
	humidity, err := getMetrics("gocast_humidity", testTimestamp)
	if err != nil {
		t.Error(err)
	}
	if humidity != float64(expectedHumidity) {
		t.Errorf("metrics returned wrong temperature: got \"%v\" want \"%v\"", humidity, expectedHumidity)
	}
	windspeed, err := getMetrics("gocast_wind_speed", testTimestamp)
	if err != nil {
		t.Error(err)
	}
	if windspeed != expectedWindSpeed {
		t.Errorf("metrics returned wrong temperature: got \"%v\" want \"%v\"", windspeed, expectedWindSpeed)
	}
	pressure, err := getMetrics("gocast_pressure", testTimestamp)
	if err != nil {
		t.Error(err)
	}
	if pressure != expectedPressure {
		t.Errorf("metrics returned wrong temperature: got \"%v\" want \"%v\"", pressure, expectedPressure)
	}
	t.Run("checking for unexpected errors", func(t *testing.T) {
		expected := 0.0
		expectedError := "metrics found no values!"
		test, err := getMetrics("gocast_failed", testTimestamp)
		if err == nil || err.Error() != expectedError {
			t.Error(err)
		}
		if test != expected {
			t.Errorf("metrics returned unexpected result: got \"%v\" want \"%v\"", pressure, expected)
		}
		test, err = getMetrics("gocast_test", testTimestamp)
		if err == nil || err.Error() != expectedError {
			t.Error(err)
		}
		if test != expected {
			t.Errorf("metrics returned unexpected result: got \"%v\" want \"%v\"", pressure, expected)
		}
	})
}

func getMetrics(metricName, timestamp string) (float64, error) {
	resp, err := http.Get("http://localhost:8081/metrics")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	reader := bytes.NewReader(body)
	var parser expfmt.TextParser
	metrics, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return 0, err
	}
	if metricFamily, ok := metrics[metricName]; ok {
		for _, metric := range metricFamily.GetMetric() {
			for _, label := range metric.GetLabel() {
				if label.GetName() == "timestamp" && label.GetValue()[:13] == timestamp[:13] {
					return metric.GetGauge().GetValue(), nil
				}
			}
		}
	}
	return 0, fmt.Errorf("metrics found no values!")
}
