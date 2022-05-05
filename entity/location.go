package entity

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Address struct {
	Address string `json:"address"`
}

type Place struct {
	Location *Location `json:"location"`
	Address  *Address  `json:"address"`
}
