package entity

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Address struct {
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	Location *Location `json:"location"`
}

type Route struct {
	Duration float64 `json:"duration"`
	Distance float64 `json:"distance"`
	Polyline string  `json:"polyline"`
}
