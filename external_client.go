package JumboLoyaltyClient

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type externalClient struct {
	balanceTable      string
	transactionTable  string
	baseUrl           string
	key               string
	secret            string
	httpClient        *http.Client
}

func (e externalClient) GetExternalBalance(cardNumber string, accountNumber int64) (*BalanceAccount, error) {
	uri := fmt.Sprintf( "%v/customer/card/%v/overview", e.baseUrl, cardNumber)
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

	externalClient.balanceTable     = balanceTable
	externalClient.transactionTable = transactionTable
	externalClient.baseUrl          = baseUrl
	externalClient.httpClient       = &http.Client{}
	externalClient.key              = key
	externalClient.secret           = secret

	return externalClient
}