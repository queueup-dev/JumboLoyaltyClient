package JumboLoyaltyClient

import "fmt"

var (
	initializedSessionClient = NewSessionClient("securequeueup-securequeueup-jumbonl-test-Sessions")
	initializedLoyaltyClient = NewJumboLoyaltyClient(
		"colleqtdirect-securequeueup-jumbonl-test-Balance",
		"colleqtdirect-securequeueup-jumbonl-test-Reservations",
		"colleqtdirect-securequeueup-jumbonl-test-ConfirmedReservations",
	)
	initializedExternalClient = NewExternalClient(
		"colleqtdirect-securequeueup-jumbonl-test-Balance",
		"colleqtdirect-securequeueup-jumbonl-test-ConfirmedReservations",
		"https://api-gw-acc.jumbo.com/ext/core/loyalty/v2/",
		"c80775c455194854ca4b118d159ded63",
		"455a20ea53e3c77ca04bac66f96c0290",
	)
)

func main() {

	externalBalance, err := initializedExternalClient.GetExternalBalance("2621922000008", 1000)

	if err != nil {
		panic(err)
	}

	balance, err := initializedLoyaltyClient.SaveBalance("1234567890123", float32(externalBalance.Balance))

	fmt.Print(balance)
}