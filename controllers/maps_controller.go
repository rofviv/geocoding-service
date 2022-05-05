package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"maps.patio.com/entity"
	"maps.patio.com/repository"
	status "maps.patio.com/responses"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var mMap repository.Repository

func New(repo repository.Repository) {
	mMap = repo
}

func IndexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Maps ROUTE")
}

func Geocoding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result := Response{}

	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result.Status = status.FAILED
		result.Message = status.FAILED_MESSAGE
	} else if body["address"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		result.Status = status.MISSING_PARAMS
		result.Message = status.MISSING_PARAMS_MESSAGE
	} else {
		addr := fmt.Sprint(body["address"])
		if len(strings.TrimSpace(addr)) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			result.Status = status.MISSING_PARAMS
			result.Message = status.EMPTY_FIELD_MESSAGE
		} else {
			statusMaps, location, err := mMap.Geocoding(addr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				result.Status = statusMaps
				result.Message = err.Error()
			} else {
				w.WriteHeader(http.StatusOK)
				result.Status = statusMaps
				result.Message = status.OK_MESSAGE
				result.Data = location
			}
		}
	}
	json.NewEncoder(w).Encode(result)
}

func ReverseGeocoding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result := Response{}

	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result.Status = status.FAILED
		result.Message = status.FAILED_MESSAGE
	} else if body["lat"] == nil || body["lng"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		result.Status = status.MISSING_PARAMS
		result.Message = status.MISSING_PARAMS_MESSAGE
	} else {
		lat, errLat := strconv.ParseFloat(fmt.Sprint(body["lat"]), 64)
		if errLat != nil {
			w.WriteHeader(http.StatusBadRequest)
			result.Status = status.INVALID_DATA
			result.Message = status.INVALID_DATA_MESSAGE + " 'lat'"
		} else {
			lng, errLng := strconv.ParseFloat(fmt.Sprint(body["lng"]), 64)
			if errLng != nil {
				w.WriteHeader(http.StatusBadRequest)
				result.Status = status.INVALID_DATA
				result.Message = status.INVALID_DATA_MESSAGE + " 'lng'"
			} else {
				location := &entity.Location{
					Lat: lat,
					Lng: lng,
				}
				statusMaps, address, err := mMap.ReverseGeocoding(location)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					result.Status = statusMaps
					result.Message = err.Error()
				} else {
					w.WriteHeader(http.StatusOK)
					result.Status = statusMaps
					result.Message = status.OK_MESSAGE
					result.Data = address
				}
			}
		}
	}
	json.NewEncoder(w).Encode(result)
}

func Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result := Response{}

	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result.Status = status.FAILED
		result.Message = status.FAILED_MESSAGE
	} else if body["address"] == nil || body["lat"] == nil || body["lng"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		result.Status = status.MISSING_PARAMS
		result.Message = status.MISSING_PARAMS_MESSAGE
	} else {
		lat, errLat := strconv.ParseFloat(fmt.Sprint(body["lat"]), 64)
		if errLat != nil {
			w.WriteHeader(http.StatusBadRequest)
			result.Status = status.INVALID_DATA
			result.Message = status.INVALID_DATA_MESSAGE + " 'lat'"
		} else {
			lng, errLng := strconv.ParseFloat(fmt.Sprint(body["lng"]), 64)
			if errLng != nil {
				w.WriteHeader(http.StatusBadRequest)
				result.Status = status.INVALID_DATA
				result.Message = status.INVALID_DATA_MESSAGE + " 'lng'"
			} else {
				addr := fmt.Sprint(body["address"])
				if len(strings.TrimSpace(addr)) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					result.Status = status.MISSING_PARAMS
					result.Message = status.EMPTY_FIELD_MESSAGE
				} else {
					location := &entity.Location{
						Lat: lat,
						Lng: lng,
					}
					statusMaps, places, err := mMap.Search(addr, location)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						result.Status = statusMaps
						result.Message = err.Error()
					} else {
						w.WriteHeader(http.StatusOK)
						result.Status = statusMaps
						result.Message = status.OK_MESSAGE
						result.Data = places
					}
				}
			}
		}
	}
	json.NewEncoder(w).Encode(result)
}
