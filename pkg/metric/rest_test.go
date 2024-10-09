package metric_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/metric"
)

/*
func TestGet(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/direct?date=2024-10-08&city=Berlin", nil)
	rr := httptest.NewRecorder()
	var h metric.Handler
	h.Get(rr, req)
	//handler := http.HandlerFunc(h.Get)
	//handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}*/

func TestGetError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/direct?date=2024-10-20&city=Berlin", nil)
	rr := httptest.NewRecorder()
	var h metric.Handler
	h.Get(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost:8080/error", nil)
	rr := httptest.NewRecorder()
	var h metric.Handler
	h.Get(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
