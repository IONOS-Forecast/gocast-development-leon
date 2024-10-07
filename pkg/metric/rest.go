package metric

import (
	"fmt"
	"log"
	"net/http"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/db"
	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/utils"
)

type Handler struct{}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	city := r.URL.Query().Get("city")
	year, month, day, err := utils.SplitDate(date)
	if err != nil {
		log.Print(err)
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	utils.SetDate(year, month, day)
	database, err := db.NewPG(utils.FdbUser, utils.FdbPass, utils.FdbDB, utils.FdbAddress)
	if err != nil {
		log.Print(err)
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	defer database.Close()
	database.QueryCitiesDatabase(&city, "name", city)
	if city == "" {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	weather, err := database.GetWeatherRecord(city, date)
	if err != nil {
		log.Print(err)
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	if weather.Hours == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusOK)
	for i, v := range weather.Hours {
		v.City = city
		fmt.Fprintf(w, fmt.Sprintf("Weather Hour %v: %v\n", i, v))
	}
	fmt.Fprintf(w, "City:"+city)
}

func (h Handler) Error(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "404 - Bad Request")
}
