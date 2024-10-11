package metric

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/db"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/utils"
)

type Rest interface {
	Get(w http.ResponseWriter, r *http.Request)
	Error(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	db           db.DBI
	cityDB       db.CityDBI
	weatherMapDB db.WeatherMapDB
}

func NewHandler(db db.DBI, cityDB db.CityDBI, weatherMap db.WeatherMapDB) Rest {
	return &handler{db: db, cityDB: cityDB, weatherMapDB: weatherMap}
}

func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	city := r.URL.Query().Get("city")
	if date != "" && city != "" {
		year, month, day, err := utils.SplitDate(date)
		if err != nil {
			log.Print(err)
			http.Redirect(w, r, "/error", http.StatusFound)
			return
		}
		date, err = utils.SetDate(year, month, day)
		if err != nil {
			log.Print(err)
			http.Redirect(w, r, "/error", http.StatusBadRequest)
			return
		}
		cityExists, err := h.cityDB.ContainsCity(city)
		if err != nil {
			log.Print(err)
			http.Redirect(w, r, "/error", http.StatusBadRequest)
			return
		}
		if cityExists {
			city, err := h.cityDB.GetCity(city)
			if err != nil {
				log.Print(err)
				http.Redirect(w, r, "/error", http.StatusBadRequest)
				return
			}
			utils.SetCity(city.Name)
			utils.SetLocation(city.Lat, city.Lon)
		} else {
			city, err = h.cityDB.SetLocationByCityName(city)
			if err != nil {
				log.Print(err)
				http.Redirect(w, r, "/error", http.StatusBadRequest)
				return
			}
		}
		h.db.QueryCitiesDatabase(&city, "name", city)
		if city == "" {
			http.Redirect(w, r, "/error", http.StatusFound)
			return
		}
		record, err := h.db.GetWeatherRecord(city, date)
		if err != nil {
			log.Print(err)
			http.Redirect(w, r, "/error", http.StatusFound)
			return
		}
		data, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			log.Print(err)
			http.Redirect(w, r, "/error", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		if record.Hours != nil && len(record.Hours) == 24 {
			UpdateMetrics(record, time.Now().Hour())
			// UpdateMetricsDay(record) // Method for whole day
			// UpdateMetricsNow(record) // Method for timestamp now
		}
	} else {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
}

func (h handler) Error(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "400 - Bad Request")
}
