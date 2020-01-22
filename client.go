package JumboLoyaltyClient

import "github.com/pkg/errors"

type client struct {
	balanceTable     string
	reservationTable string
	transactionTable string
}

func (c client) Reserve(sessionId string, externalId string, amount float32) (*reservation, error) {
	ok := c.hasEnoughBalance(externalId, amount)

	if !ok {
		return nil, errors.New("Invalid balance supplied or not enough balance available")
	}

	reservation := NewReservation(externalId, sessionId, amount)

	_, err := c.processReservation(reservation)

	return reservation, err
}

func (c client) Release(reservationId string) error {
	reservation, err := c.GetReservation(reservationId)

	if err != nil {
		return err
	}

	_, err = c.processRelease(reservation)

	return err
}

func (c client) Sell(reservationId string) (*transaction, error) {
	reservation, err := c.GetReservation(reservationId)

	if err != nil {
		return nil, err
	}

	return c.processSale(reservation)
}

func (c client) GetReservation(reservationId string) (*reservation, error) {
	result, err := Dynamo.getItem(c.reservationTable, "reservation_id", reservationId, &reservation{})

	if err != nil {
		return nil, err
	}

	reservation := result.(*reservation)
	return reservation, nil
}

func (c client) GetBalance(externalId string) (*balance, error) {
	result, err := Dynamo.getItem(c.balanceTable, "external_id", externalId, &balance{})

	if err != nil {
		return nil, err
	}

	balance := result.(*balance)
	return balance, nil
}

func (c client) ListReservations(sessionId string) (*[]reservation, error) {

	reservations := &[]reservation{}
	result, err := Dynamo.listItems(c.reservationTable, "session_id-index", []QueryCondition{
		{
			Key:       "session_id",
			Value:     sessionId,
			Operation: "EQ",
		},
	}, reservations, 10)

	if err != nil {
		return nil, err
	}

	reservations = result.(*[]reservation)
	return reservations, nil
}

func (c client) processSale(reservation *reservation) (*transaction, error) {

	transaction := NewTransaction(reservation.ExternalId, reservation.Amount)

	err := Dynamo.saveItem(c.transactionTable, transaction)

	if err != nil {
		return nil, err
	}

	_, err = c.processRelease(reservation)

	return transaction, err
}

func (c client) processReservation(reservation *reservation) (*balance, error) {
	balance, err := c.GetBalance(reservation.ExternalId)

	if err != nil {
		return nil, err
	}

	balance.Reserved += reservation.Amount

	err = Dynamo.saveItem(c.balanceTable, balance)

	if err != nil {
		return nil, err
	}

	err = Dynamo.saveItem(c.reservationTable, reservation)

	return balance, err
}

func (c client) processRelease(reservation *reservation) (*balance, error) {
	balance, err := c.GetBalance(reservation.ExternalId)

	if err != nil {
		return nil, err
	}

	balance.Reserved -= reservation.Amount

	err = Dynamo.saveItem(c.balanceTable, balance)

	if err != nil {
		return nil, err
	}

	err = Dynamo.deleteItem(c.reservationTable, "reservation_id", reservation.ReservationId)

	return balance, err
}

func (c client) hasEnoughBalance(externalId string, amount float32) bool {

	balance, err := c.GetNumericBalance(externalId)

	if err != nil {
		return false
	}

	if amount <= 0 {
		return false
	}

	return balance >= amount
}

func (c client) GetNumericBalance(externalId string) (float32, error) {
	balance, err := c.GetBalance(externalId)

	if err != nil {
		return 0, err
	}

	return balance.Balance - (balance.Reserved + balance.Pending), nil
}

func NewJumboLoyaltyClient(
	balanceTable string,
	reservationTable string,
	transactionTable string,
) *client {
	client := new(client)

	client.reservationTable = reservationTable
	client.transactionTable = transactionTable
	client.balanceTable = balanceTable

	return client
}
