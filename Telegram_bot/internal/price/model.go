package price

import "go.mod/internal/coin"

type Price struct {
	ID      string    `json:"id"`
	HowMany float64   `json:"how_many"`
	Coin    coin.Coin `json:"coin"`
}
