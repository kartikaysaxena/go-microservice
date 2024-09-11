package types

type RateRequest struct {
	RateLimit float32 `json:"rate_limit"`
	Minutes   float32 `json:"minutes"`
}

type RateResponse struct {
	Count float32 `json:"count"`
}
