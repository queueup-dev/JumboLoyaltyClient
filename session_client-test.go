package JumboLoyaltyClient

import "fmt"

var (
	initializedSessionClient = NewSessionClient("securequeueup-securequeueup-jumbonl-test-Sessions")
	initializedLoyaltyClient = NewJumboLoyaltyClient(
		"colleqtdirect-securequeueup-jumbonl-test-Balance",
		"colleqtdirect-securequeueup-jumbonl-test-Reservations",
		"colleqtdirect-securequeueup-jumbonl-test-ConfirmedReservations",
	)
)

func main() {
	result, err := initializedSessionClient.GetSession("0ae67e87-8ffa-4e72-b149-1513cfbf3743")

	if err != nil {
		panic(err)
	}

	fmt.Print(result)

	reservation, err := initializedLoyaltyClient.GetReservation("0031a636-b236-4fdc-aca5-fa0180a60d78")

	if err != nil {
		panic(err)
	}

	fmt.Print(reservation)

	reservations, err := initializedLoyaltyClient.ListReservations("db4fa44c-2098-49cc-9c50-873dbb1afeb3")

	if err != nil {
		panic(err)
	}

	fmt.Print(reservations)
}