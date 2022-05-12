package heremaps

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/heremaps/flexible-polyline/golang/flexpolyline"
	"maps.patio.com/entity"
	status "maps.patio.com/responses"
)

type HereMaps struct {
	ApiKey string
}

type Items struct {
	Items            []Item `json:"items"`
	ErrorDescription string `json:"error_description"`
}

type Item struct {
	Title    string          `json:"title"`
	Address  AddressLabel    `json:"address"`
	Location entity.Location `json:"position"`
}

type AddressLabel struct {
	Label string `json:"label"`
}

type Response struct {
	Routes           []Route `json:"routes"`
	ErrorDescription string  `json:"error_description"`
}
type Route struct {
	Sections []Section `json:"sections"`
}

type Section struct {
	Summary  Summary `json:"summary"`
	Polyline string  `json:"polyline"`
}

type Summary struct {
	Duration float64 `json:"duration"`
	Distance float64 `json:"length"`
}

func New(key string) *HereMaps {
	return &HereMaps{
		ApiKey: key,
	}
}
func (h *HereMaps) Provider() string {
	return "HERE MAPS"
}

func (h *HereMaps) Geocoding(address string) (string, *entity.Address, error) {
	params := url.Values{}
	params.Add("q", address)
	params.Add("apikey", h.ApiKey)

	var uri string = fmt.Sprintf("https://geocode.search.hereapi.com/v1/geocode?%s", params.Encode())

	resp, err := http.Get(uri)
	if err != nil {
		return status.FAILED, nil, err
	}
	defer resp.Body.Close()

	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return status.FAILED, nil, err
	}
	var items Items
	errUnmarshal := json.Unmarshal(bytes, &items)
	if errUnmarshal != nil {
		return status.FAILED, nil, err
	}

	if len(items.Items) == 0 {
		if items.ErrorDescription != "" {
			return status.ZERO_RESULTS, nil, errors.New(items.ErrorDescription)
		}
		return status.ZERO_RESULTS, nil, errors.New("No results for " + address)
	} else {
		newAddress := &entity.Address{
			Name:     items.Items[0].Title,
			Address:  items.Items[0].Address.Label,
			Location: &items.Items[0].Location,
		}
		return status.OK, newAddress, nil
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
		return status.FAILED, nil, err
	}

	defer resp.Body.Close()
	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return status.FAILED, nil, err
	}

	var items Items
	errUnmarshal := json.Unmarshal(bytes, &items)
	if errUnmarshal != nil {
		return status.FAILED, nil, err
	}

	if len(items.Items) == 0 {
		if items.ErrorDescription != "" {
			return status.ZERO_RESULTS, nil, errors.New(items.ErrorDescription)
		}
		return status.ZERO_RESULTS, nil, errors.New("No results for " + latlng)
	} else {
		address := &entity.Address{
			Name:     items.Items[0].Title,
			Address:  items.Items[0].Address.Label,
			Location: &items.Items[0].Location,
		}
		return status.OK, address, nil
	}
}

func (h *HereMaps) Search(address string, location *entity.Location) (string, []*entity.Address, error) {

	latlng := fmt.Sprintf("%f,%f", location.Lat, location.Lng)
	params := url.Values{}
	params.Add("at", latlng)
	params.Add("q", address)
	params.Add("apikey", h.ApiKey)
	params.Add("lang", "en-US")

	var uri string = fmt.Sprintf("https://autosuggest.search.hereapi.com/v1/autosuggest?%s", params.Encode())
	resp, err := http.Get(uri)
	if err != nil {
		return status.FAILED, nil, err
	}

	defer resp.Body.Close()
	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return status.FAILED, nil, err
	}

	var items Items
	errUnmarshal := json.Unmarshal(bytes, &items)
	if errUnmarshal != nil {
		return status.FAILED, nil, err
	}

	if len(items.Items) == 0 {
		if items.ErrorDescription != "" {
			return status.ZERO_RESULTS, nil, errors.New(items.ErrorDescription)
		}
		return status.ZERO_RESULTS, nil, errors.New("No results for " + address)
	} else {
		list := []*entity.Address{}
		for _, v := range items.Items {
			locationTmp := &entity.Location{
				Lat: v.Location.Lat,
				Lng: v.Location.Lng,
			}
			placeTmp := &entity.Address{
				Name:     v.Title,
				Address:  v.Address.Label,
				Location: locationTmp,
			}
			list = append(list, placeTmp)
		}
		return status.OK, list, nil
	}

}

func (h *HereMaps) Distance(origin *entity.Location, destination *entity.Location) (string, *entity.Summary, error) {
	from := fmt.Sprintf("%f,%f", origin.Lat, origin.Lng)
	to := fmt.Sprintf("%f,%f", destination.Lat, destination.Lng)
	params := url.Values{}
	params.Add("origin", from)
	params.Add("destination", to)
	params.Add("transportMode", "bicycle")
	params.Add("return", "summary")
	params.Add("apikey", h.ApiKey)

	var uri string = fmt.Sprintf("https://router.hereapi.com/v8/routes?%s", params.Encode())
	resp, err := http.Get(uri)
	if err != nil {
		return status.FAILED, nil, err
	}

	defer resp.Body.Close()
	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return status.FAILED, nil, err
	}

	var response Response
	errUnmarshal := json.Unmarshal(bytes, &response)
	if errUnmarshal != nil {
		return status.FAILED, nil, err
	}

	if len(response.Routes) <= 0 {
		if response.ErrorDescription != "" {
			return status.ZERO_RESULTS, nil, errors.New(response.ErrorDescription)
		}
		return status.ZERO_RESULTS, nil, errors.New("Distance for origin or destination invalid")
	}

	var summary = &entity.Summary{
		Duration: response.Routes[0].Sections[0].Summary.Duration,
		Distance: response.Routes[0].Sections[0].Summary.Distance,
	}

	return status.OK, summary, nil
}
func (h *HereMaps) Route(origin *entity.Location, destination *entity.Location) (string, *entity.Route, error) {
	from := fmt.Sprintf("%f,%f", origin.Lat, origin.Lng)
	to := fmt.Sprintf("%f,%f", destination.Lat, destination.Lng)
	params := url.Values{}
	params.Add("origin", from)
	params.Add("destination", to)
	params.Add("transportMode", "bicycle")
	params.Add("return", "polyline,summary")
	params.Add("apikey", h.ApiKey)

	var uri string = fmt.Sprintf("https://router.hereapi.com/v8/routes?%s", params.Encode())
	resp, err := http.Get(uri)
	if err != nil {
		return status.FAILED, nil, err
	}

	defer resp.Body.Close()
	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return status.FAILED, nil, err
	}

	var response Response
	errUnmarshal := json.Unmarshal(bytes, &response)
	if errUnmarshal != nil {
		return status.FAILED, nil, err
	}

	if len(response.Routes) <= 0 {
		if response.ErrorDescription != "" {
			return status.ZERO_RESULTS, nil, errors.New(response.ErrorDescription)
		}
		return status.ZERO_RESULTS, nil, errors.New("Route for origin or destination invalid")
	}

	poly, err := flexpolyline.Decode(response.Routes[0].Sections[0].Polyline)
	if err != nil {
		return status.FAILED, nil, err
	}

	// TODO: VER LA MANERA DE ENCRIPTAR LA LIST A UN STRING CORTO PARA NO MANDAR MUCHOS DATOS POR LA RED
	var list []*entity.Location

	for _, v := range poly.Coordinates() {
		var locationTmp = entity.Location{Lat: v.Lat, Lng: v.Lng}
		list = append(list, &locationTmp)
	}

	summaryTmp := entity.Summary{
		Duration: response.Routes[0].Sections[0].Summary.Duration,
		Distance: response.Routes[0].Sections[0].Summary.Distance,
	}

	var route = &entity.Route{
		Summary:  summaryTmp,
		Polyline: list,
	}

	return status.OK, route, nil
}
