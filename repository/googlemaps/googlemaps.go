package googlemaps

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/twpayne/go-polyline"
	"maps.patio.com/entity"
	status "maps.patio.com/responses"
)

type GoogleMaps struct {
	ApiKey string
}

type Results struct {
	Results      []ResultItem `json:"results"`
	Status       string       `json:"status"`
	ErrorMessage string       `json:"error_message"`
}

type ResultItem struct {
	ResultItem Geometry `json:"geometry"`
	Address    string   `json:"formatted_address"`
}

type Geometry struct {
	Location entity.Location `json:"location"`
}

type Response struct {
	Rows         []Row  `json:"rows"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

type Row struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Distance       ValueFloat `json:"distance"`
	Duration       ValueFloat `json:"duration"`
	StatusDistance string     `json:"Status"`
}

type ResponseRoute struct {
	Routes []Route `json:"routes"`
	Status string  `json:"status"`
}

type Route struct {
	Legs             []Leg            `json:"legs"`
	OverviewPolyline OverviewPolyline `json:"overview_polyline"`
}

type Leg struct {
	Distance ValueFloat `json:"distance"`
	Duration ValueFloat `json:"duration"`
}

type OverviewPolyline struct {
	Points string `json:"points"`
}

type ValueFloat struct {
	Value float64 `json:"value"`
}

func New(key string) *GoogleMaps {
	return &GoogleMaps{
		ApiKey: key,
	}
}

func (g *GoogleMaps) Provider() string {
	return "GOOGLE MAPS"
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
		if results.ErrorMessage != "" {
			return results.Status, nil, errors.New(results.ErrorMessage)
		}
		return results.Status, nil, errors.New("No results for " + address)
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
		if results.ErrorMessage != "" {
			return results.Status, nil, errors.New(results.ErrorMessage)
		}
		return results.Status, nil, errors.New("No results for " + latlng)
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
		if results.ErrorMessage != "" {
			return results.Status, nil, errors.New(results.ErrorMessage)
		}
		return results.Status, nil, errors.New("No results for " + address)
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

func (g *GoogleMaps) Distance(origin *entity.Location, destination *entity.Location) (string, *entity.Summary, error) {
	from := fmt.Sprintf("%f,%f", origin.Lat, origin.Lng)
	to := fmt.Sprintf("%f,%f", destination.Lat, destination.Lng)
	params := url.Values{}
	params.Add("origins", from)
	params.Add("destinations", to)
	params.Add("mode", "driving")
	params.Add("key", g.ApiKey)

	var uri string = fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?%s", params.Encode())
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

	if len(response.Rows) <= 0 {
		if response.ErrorMessage != "" {
			return response.Status, nil, errors.New(response.ErrorMessage)
		}
		return response.Status, nil, errors.New("Distance for origin or destination invalid")
	}

	var summary = &entity.Summary{
		Duration: response.Rows[0].Elements[0].Duration.Value,
		Distance: response.Rows[0].Elements[0].Distance.Value,
	}

	if response.Rows[0].Elements[0].StatusDistance != "" && response.Rows[0].Elements[0].StatusDistance == status.ZERO_RESULTS {
		return status.ZERO_RESULTS, nil, errors.New("failed to calculate distance")
	}

	return status.OK, summary, nil
}

// TODO: ROUTES
func (g *GoogleMaps) Route(origin *entity.Location, destination *entity.Location) (string, *entity.Route, error) {
	from := fmt.Sprintf("%f,%f", origin.Lat, origin.Lng)
	to := fmt.Sprintf("%f,%f", destination.Lat, destination.Lng)
	params := url.Values{}
	params.Add("origin", from)
	params.Add("destination", to)
	params.Add("mode", "driving")
	params.Add("key", g.ApiKey)

	var uri string = fmt.Sprintf("https://maps.googleapis.com/maps/api/directions/json?%s", params.Encode())
	resp, err := http.Get(uri)
	if err != nil {
		return status.FAILED, nil, err
	}

	defer resp.Body.Close()
	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return status.FAILED, nil, err
	}

	var responseRoute ResponseRoute
	errUnmarshal := json.Unmarshal(bytes, &responseRoute)
	if errUnmarshal != nil {
		return status.FAILED, nil, err
	}

	if len(responseRoute.Routes) <= 0 {
		return status.ZERO_RESULTS, nil, errors.New("Route for origin or destination invalid")
	}
	buf := []byte(responseRoute.Routes[0].OverviewPolyline.Points)
	coords, _, err := polyline.DecodeCoords(buf)
	if err != nil {
		return status.FAILED, nil, err
	}

	var list []*entity.Location

	for _, v := range coords {
		var locationTmp = entity.Location{Lat: v[0], Lng: v[1]}
		list = append(list, &locationTmp)
	}

	summaryTmp := entity.Summary{
		Duration: responseRoute.Routes[0].Legs[0].Duration.Value,
		Distance: responseRoute.Routes[0].Legs[0].Distance.Value,
	}

	var route = &entity.Route{
		Summary:  summaryTmp,
		Polyline: list,
	}

	return status.OK, route, nil
}
