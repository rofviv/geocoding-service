package heremaps

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

type HereMaps struct {
	ApiKey string
}

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	Title    string          `json:"title"`
	Address  AddressLabel    `json:"address"`
	Location entity.Location `json:"position"`
}

type AddressLabel struct {
	Label string `json:"label"`
}

func New(key string) *HereMaps {
	return &HereMaps{
		ApiKey: key,
	}
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
		return status.ZERO_RESULTS, nil, errors.New("no results for " + address)
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
		return status.ZERO_RESULTS, nil, errors.New("no results for " + latlng)
	} else {
		address := &entity.Address{
			Name:     items.Items[0].Title,
			Address:  items.Items[0].Address.Label,
			Location: &items.Items[0].Location,
		}
		return status.OK, address, nil
	}
}

// TODO: CREAR UN MODELO PARA LEER LA RESPUESTA DEL SEARCH. DEBE DEVOLER NAME, ADDRESS, LOCATION
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
		return status.ZERO_RESULTS, nil, errors.New("no results for " + latlng)
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
		// address := &entity.Address{
		// 	Address:  items.Items[0].Title,
		// 	Name:     strings.Split(items.Items[0].Title, ",")[0],
		// 	Location: &items.Items[0].Location,
		// }
		return status.OK, list, nil
	}

}

// TODO: ROUTES
