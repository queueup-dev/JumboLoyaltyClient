package JumboLoyaltyClient

import (
	"cirello.io/dynamolock"
	"errors"
	"time"
)

var (
	lockClient *dynamolock.Client
)

type voucherClient struct {
	voucherTable string
}

type voucher struct {
	Code string `json:"voucher"`
	Used int8   `json:"used"`
}

func (v voucherClient) ReserveVouchers(amount int) (*[]voucher, error) {

	lock, err := lockClient.AcquireLock("reserveVouchers")

	if err != nil {
		panic(err)
	}

	defer lockClient.ReleaseLock(lock)

	vouchers := &[]voucher{}

	result, err := Dynamo.findItems(v.voucherTable, "used = :used", map[string]string{
		":used": "0",
	}, vouchers, int64(amount))

	if err != nil {
		return nil, err
	}

	vouchers = result.(*[]voucher)

	if len(*vouchers) != amount {
		return nil, errors.New("not enough vouchers available for this reservation")
	}

	for _, voucher := range *vouchers {
		voucher.Used = 1
		err = Dynamo.saveItem(v.voucherTable, voucher)
	}

	return vouchers, nil
}

func NewVoucherClient(lockTable string, voucherTable string) *voucherClient {

	client, err := dynamolock.New(
		Dynamo.db,
		lockTable,
		dynamolock.WithLeaseDuration(2*time.Second),
		dynamolock.WithHeartbeatPeriod(3*time.Millisecond),
	)

	if err != nil {
		panic(err)
	}

	lockClient = client
	voucherClient := new(voucherClient)
	voucherClient.voucherTable = voucherTable

	return voucherClient
}
