package googlemaps

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"maps.patio.com/entity"
)

type GoogleMaps struct {
	ApiKey string
}

type Results struct {
	Results []ResultItem `json:"results"`
	Status  string       `json:"status"`
}

type ResultItem struct {
	ResultItem Geometry `json:"geometry"`
	Address    string   `json:"formatted_address"`
}

type Geometry struct {
	Location entity.Location `json:"location"`
}

func New(key string) *GoogleMaps {
	return &GoogleMaps{
		ApiKey: key,
	}
}

func (g *GoogleMaps) Geocoding(address string) (string, *entity.Location, error) {

	params := url.Values{}
	params.Add("address", address)
	params.Add("key", g.ApiKey)

	var uri string = fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?%s", params.Encode())

	resp, err := http.Get(uri)
	if err != nil {
		return "FAILED", nil, err
	}

	defer resp.Body.Close()

	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return "FAILED", nil, err
	}

	var results Results
	errUnmarshal := json.Unmarshal(bytes, &results)
	if errUnmarshal != nil {
		return results.Status, nil, err
	}
	if len(results.Results) == 0 {
		return results.Status, nil, errors.New("no results for " + address)
	} else {
		return results.Status, &results.Results[0].ResultItem.Location, nil
	}
}

func (g *GoogleMaps) ReverseGeocoding(location *entity.Location) (string, *entity.Address, error) {
	latlng := fmt.Sprintf("%f,%f", location.Lat, location.Lng)
	params := url.Values{}
	params.Add("latlng", latlng)
	params.Add("key", g.ApiKey)

	var uri string = fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?%s", params.Encode())

	resp, err := http.Get(uri)
	if err != nil {
		return "FAILED", nil, err
	}

	defer resp.Body.Close()
	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return "FAILED", nil, err
	}

	var results Results
	errUnmarshal := json.Unmarshal(bytes, &results)
	if errUnmarshal != nil {
		return results.Status, nil, err
	}
	if len(results.Results) == 0 {
		return results.Status, nil, errors.New("no results for " + latlng)
	} else {
		address := &entity.Address{
			Address: results.Results[0].Address,
		}
		return "OK", address, nil
	}
}

func (g *GoogleMaps) Search(address string, location *entity.Location) (string, []*entity.Place, error) {
	// latlng := fmt.Sprintf("%f,%f", location.Lat, location.Lng)
	params := url.Values{}
	params.Add("keyword", address)
	// params.Add("location", latlng)
	params.Add("key", g.ApiKey)
	// params.Add("radius", "1500")
	params.Add("libraries", "places")
	// params.Add("type", "restaurant")

	var uri string = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?%s", params.Encode())
	resp, err := http.Get(uri)
	if err != nil {
		return "FAILED", nil, err
	}

	defer resp.Body.Close()
	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return "FAILED", nil, err
	}

	fmt.Println(string(bytes))
	return "", nil, nil
}
