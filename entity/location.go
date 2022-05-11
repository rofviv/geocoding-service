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

// type Place struct {
// 	Location *Location `json:"location"`
// 	Address  *Address  `json:"address"`
// }
