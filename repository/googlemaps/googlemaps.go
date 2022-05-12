package googlemaps

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"maps.patio.com/entity"
	status "maps.patio.com/responses"
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

func (g *GoogleMaps) Geocoding(address string) (string, *entity.Address, error) {

	params := url.Values{}
	params.Add("address", address)
	params.Add("key", g.ApiKey)

	var uri string = fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?%s", params.Encode())

	resp, err := http.Get(uri)
	if err != nil {
		return status.FAILED, nil, err
	}

	defer resp.Body.Close()

	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return status.FAILED, nil, err
	}

	var results Results
	errUnmarshal := json.Unmarshal(bytes, &results)
	if errUnmarshal != nil {
		return results.Status, nil, err
	}
	if len(results.Results) == 0 {
		return results.Status, nil, errors.New("no results for " + address)
	} else {
		newAddress := &entity.Address{
			Name:     strings.Split(results.Results[0].Address, ",")[0],
			Address:  results.Results[0].Address,
			Location: &results.Results[0].ResultItem.Location,
		}
		return results.Status, newAddress, nil
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
		return status.FAILED, nil, err
	}

	defer resp.Body.Close()
	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return status.FAILED, nil, err
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
			Address:  results.Results[0].Address,
			Name:     strings.Split(results.Results[0].Address, ",")[0],
			Location: &results.Results[0].ResultItem.Location,
		}
		return status.OK, address, nil
	}
}

func (g *GoogleMaps) Search(address string, location *entity.Location) (string, []*entity.Address, error) {
	latlng := fmt.Sprintf("%f,%f", location.Lat, location.Lng)
	params := url.Values{}
	params.Add("query", address)
	params.Add("location", latlng)
	params.Add("key", g.ApiKey)

	var uri string = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/textsearch/json?%s", params.Encode())
	resp, err := http.Get(uri)
	if err != nil {
		return status.FAILED, nil, err
	}

	defer resp.Body.Close()
	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return status.FAILED, nil, err
	}

	var results Results
	errUnmarshal := json.Unmarshal(bytes, &results)
	if errUnmarshal != nil {
		return results.Status, nil, err
	}
	if len(results.Results) == 0 {
		return results.Status, nil, errors.New("no results for " + address)
	} else {
		list := []*entity.Address{}
		for _, v := range results.Results {
			locationTmp := &entity.Location{
				Lat: v.ResultItem.Location.Lat,
				Lng: v.ResultItem.Location.Lng,
			}
			placeTmp := &entity.Address{
				Name:     strings.Split(v.Address, ",")[0],
				Address:  v.Address,
				Location: locationTmp,
			}
			list = append(list, placeTmp)
		}
		return status.OK, list, nil
	}

}

// TODO: ROUTES
func (h *GoogleMaps) Route(origin *entity.Location, destination *entity.Location) (string, *entity.Route, error) {

	return status.OK, nil, nil
}
