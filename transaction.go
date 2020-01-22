package JumboLoyaltyClient

import "github.com/google/uuid"

type transaction struct {
	TransactionId  string   `json:"confirmation_id"`
	ExternalId     string   `json:"external_id"`
	Amount         float32  `json:"amount"`
	Announced      bool     `json:"announced"`
}

func NewTransaction(externalId string, amount float32) *transaction {
	transaction := new(transaction)

	transaction.TransactionId = uuid.New().String()
	transaction.ExternalId    = externalId
	transaction.Amount        = amount
	transaction.Announced     = false

	return transaction
}