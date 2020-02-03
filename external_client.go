package JumboLoyaltyClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type BalanceResponse struct {
	Id        string `json:"id"`
	Status    int32  `json:"status"`
	ExpiresOn string `json:"expiresOn"`
	Accounts  []BalanceAccount
}

type BalanceAccount struct {
	Id      int64  `json:"id"`
	Balance int32  `json:"balance"`
	Updated string `json:"updateOn"`
}

type externalTransaction struct {
	Amount int32 `json:"amount"`
}

type externalClient struct {
	balanceTable     string
	transactionTable string
	baseUrl          string
	key              string
	secret           string
	httpClient       *http.Client
}

func (e externalClient) Spend(cardNumber string, amount int32) error {
	uri := fmt.Sprintf("%v/customer/card/%v/redeem", e.baseUrl, cardNumber)

	marshalledBody, err := json.Marshal(&externalTransaction{Amount: amount})

	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", uri, strings.NewReader(string(marshalledBody)))

	if err != nil {
		return err
	}

	request.Header.Add("x-ibm-client-id", e.key)
	request.Header.Add("x-ibm-client-secret", e.secret)

	result, err := e.httpClient.Do(request)

	if err != nil {
		return err
	}

	if result.StatusCode >= 400 {
		return errors.New("there was a problem retrieving your card details")
	}

	return nil
}

func (e externalClient) GetExternalBalance(cardNumber string, accountNumber int64) (*BalanceAccount, error) {
	uri := fmt.Sprintf("%v/customer/card/%v/overview", e.baseUrl, cardNumber)
	request, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("x-ibm-client-id", e.key)
	request.Header.Add("x-ibm-client-secret", e.secret)

	result, err := e.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	if result.StatusCode >= 400 {
		return nil, errors.New("there was a problem retrieving your card details")
	}

	response := &BalanceResponse{}
	err = json.NewDecoder(result.Body).Decode(&response)

	if err != nil {
		return nil, err
	}

	var balance BalanceAccount
	for _, val := range response.Accounts {

		if val.Id == accountNumber {
			balance = val
			return &balance, nil
		}
	}

	return nil, nil
}

func NewExternalClient(
	balanceTable string,
	transactionTable string,
	baseUrl string,
	key string,
	secret string,
) *externalClient {
	externalClient := new(externalClient)

	externalClient.balanceTable = balanceTable
	externalClient.transactionTable = transactionTable
	externalClient.baseUrl = baseUrl
	externalClient.httpClient = &http.Client{}
	externalClient.key = key
	externalClient.secret = secret

	return externalClient
}
