package metric

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/db"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/utils"
)

type Rest interface {
	Get(w http.ResponseWriter, r *http.Request)
	Error(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	db     db.DBI
	cityDB db.CityDBI
}

func NewHandler(db db.DBI, cityDB db.CityDBI) Rest {
	return &handler{db: db, cityDB: cityDB}
}

func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	city := r.URL.Query().Get("city")
	if date != "" && city != "" {
		year, month, day, err := utils.SplitDate(date)
		if err != nil {
			log.Print(err)
			http.Redirect(w, r, "/error", http.StatusBadRequest)
			return
		}
		city, err = utils.SetDateAndLocationByCityName(year, month, day, city, utils.GetCities())
		if err != nil {
			log.Print(err)
			http.Redirect(w, r, "/error", http.StatusBadRequest)
			return
		}
		/*database, err := db.NewPG(utils.FdbUser, utils.FdbPass, utils.FdbDB, utils.FdbAddress)
		if err != nil {
			log.Print(err)
			http.Redirect(w, r, "/error", http.StatusBadRequest)
			return
		}
		defer database.Close()*/
		h.db.QueryCitiesDatabase(&city, "name", city)
		if city == "" {
			http.Redirect(w, r, "/error", http.StatusBadRequest)
			return
		}
		record, err := h.db.GetWeatherRecord(city, date)
		if err != nil {
			log.Print(err)
			http.Redirect(w, r, "/error", http.StatusBadRequest)
			return
		}
		if record.Hours == nil {
			w.WriteHeader(http.StatusOK)
			return
		}
		data, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			log.Print(err)
			http.Redirect(w, r, "/error", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		RegisterMetrics(record, 0)
		// Method for whole day
		// RegisterDayMetrics(record)
	} else {
		http.Redirect(w, r, "/error", http.StatusBadRequest)
		return
	}
}

func (h handler) Error(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "400 - Bad Request")
}
