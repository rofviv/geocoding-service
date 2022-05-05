package repository

import (
	"fmt"

	"maps.patio.com/configuration"
	"maps.patio.com/entity"
	"maps.patio.com/repository/googlemaps"
	"maps.patio.com/repository/heremaps"
)

type Repository interface {
	Geocoding(address string) (status string, location *entity.Location, err error)
	ReverseGeocoding(location *entity.Location) (status string, address *entity.Address, err error)
	Search(address string, location *entity.Location) (status string, places []*entity.Place, err error)
}

func New(config *configuration.Configuration) (Repository, error) {

	var repo Repository
	var err error

	switch config.MAPS.Provider {
	case "google_maps":
		repo = googlemaps.New(config.MAPS.ApiKey)
	case "here_maps":
		repo = heremaps.New(config.MAPS.ApiKey)
	default:
		err = fmt.Errorf("invalid engine %v", config.MAPS.Provider)
	}

	return repo, err
}