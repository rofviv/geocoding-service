package heremaps

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"maps.patio.com/entity"
)

type HereMaps struct {
	ApiKey string
}

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	Title    string          `json:"title"`
	Location entity.Location `json:"position"`
}

func New(key string) *HereMaps {
	return &HereMaps{
		ApiKey: key,
	}
}

func (h *HereMaps) Geocoding(address string) (string, *entity.Location, error) {
	params := url.Values{}
	params.Add("q", address)
	params.Add("apikey", h.ApiKey)

	var uri string = fmt.Sprintf("https://geocode.search.hereapi.com/v1/geocode?%s", params.Encode())

	resp, err := http.Get(uri)
	if err != nil {
		return "FAILED", nil, err
	}
	defer resp.Body.Close()

	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return "FAILED", nil, err
	}
	var items Items
	errUnmarshal := json.Unmarshal(bytes, &items)
	if errUnmarshal != nil {
		return "FAILED", nil, err
	}

	if len(items.Items) == 0 {
		return "ZERO_RESULTS", nil, errors.New("no results for " + address)
	} else {
		return "OK", &items.Items[0].Location, nil
	}
}

func (h *HereMaps) ReverseGeocoding(location *entity.Location) (string, *entity.Address, error) {
	latlng := fmt.Sprintf("%f,%f", location.Lat, location.Lng)
	params := url.Values{}
	params.Add("at", latlng)
	params.Add("apikey", h.ApiKey)
	params.Add("lang", "en-US")

	var uri string = fmt.Sprintf("https://revgeocode.search.hereapi.com/v1/revgeocode?%s", params.Encode())

	resp, err := http.Get(uri)
	if err != nil {
		return "FAILED", nil, err
	}

	defer resp.Body.Close()
	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return "FAILED", nil, err
	}

	var items Items
	errUnmarshal := json.Unmarshal(bytes, &items)
	if errUnmarshal != nil {
		return "FAILED", nil, err
	}

	if len(items.Items) == 0 {
		return "ZERO_RESULTS", nil, errors.New("no results for " + latlng)
	} else {
		address := &entity.Address{
			Address: items.Items[0].Title,
		}
		return "OK", address, nil
	}
}

func (h *HereMaps) Search(address string, location *entity.Location) (string, []*entity.Place, error) {
	return "", nil, nil
}
