package price

type CreatePriceDTO struct {
	HowMany float64 `json:"how_many"`
	CoinID  int     `json:"coin_id"`
}
