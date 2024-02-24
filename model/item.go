package model

type Item struct {
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	LastPrice float64 `json:"last_price"`
	Change    float64 `json:"change"`
}
