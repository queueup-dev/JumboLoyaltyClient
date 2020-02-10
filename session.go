package JumboLoyaltyClient

type Session struct {
	SessionId  string `json:"session_id"`
	ExternalId string `json:"card_number"`
	Email      string `json:"email"`
	Expires    int64  `json:"expires"`
	FirstName  string `json:"first_name"`
	PostalCode string `json:"postal_code"`
}
