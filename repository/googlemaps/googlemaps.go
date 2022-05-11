package googlemaps

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

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

func (g *GoogleMaps) Geocoding(address string) (string, *entity.Location, error) {

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
			Address: results.Results[0].Address,
		}
		return status.OK, address, nil
	}
}

// TODO: CREAR UN MODELO PARA LEER LA RESPUESTA DEL SEARCH. DEBE DEVOLER NAME, ADDRESS, LOCATION
func (g *GoogleMaps) Search(address string, location *entity.Location) (string, []*entity.Place, error) {
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
		return results.Status, nil, errors.New("no results for " + latlng)
	} else {
		// TODO: RECORRER ARRAY RESULTS
		fmt.Println(results.Results[0])
		return status.OK, nil, nil
	}

}

// TODO: ROUTES
