package JumboLoyaltyClient

import "github.com/google/uuid"

type reservation struct {
	ReservationId string  `json:"reservation_id"`
	ExternalId    string  `json:"external_id"`
	SessionId     string  `json:"session_id"`
	Amount        float32 `json:"amount"`
	Expires       string  `json:"expires"`
}

func NewReservation(externalId string, sessionId string, amount float32) *reservation {
	reservation := new(reservation)

	reservation.ReservationId = uuid.New().String()
	reservation.ExternalId = externalId
	reservation.SessionId = sessionId
	reservation.Amount = amount

	return reservation
}