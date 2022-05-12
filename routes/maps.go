package routes

import (
	"github.com/gorilla/mux"
	ctrl "maps.patio.com/controllers"
	"maps.patio.com/repository"
)

func Maps(repo repository.Repository) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	ctrl.New(repo)

	router.HandleFunc("/", ctrl.IndexRoute)
	router.HandleFunc("/geocoding", ctrl.Geocoding).Methods("POST")
	router.HandleFunc("/reverse-geocoding", ctrl.ReverseGeocoding).Methods("POST")
	router.HandleFunc("/search", ctrl.Search).Methods("POST")
	router.HandleFunc("/distance", ctrl.Distance).Methods("POST")
	router.HandleFunc("/route", ctrl.Route).Methods("POST")

	return router
}
