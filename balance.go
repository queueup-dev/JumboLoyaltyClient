package JumboLoyaltyClient

type balance struct {
	ExternalId string  `json:"external_id"`
	Balance    float32 `json:"balance"`
	Reserved   float32 `json:"reserved"`
	Pending    float32 `json:"pending"`
	Credit     float32 `json:"credit"`
}
