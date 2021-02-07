package plaid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type PlaidAPI struct {
}

type GetPublicToken struct {
	ClientId        string            `json:"client_id"`
	Secret          string            `json:"secret"`
	InstitutionId   string            `json:"institution_id"`
	InitialProducts []string          `json:"initial_products"`
	Options         map[string]string `json:"options"`
}

type PublicToken struct {
	Token string `json:"public_token"`
}

type GetAccessToken struct {
	ClientId    string `json:"client_id"`
	Secret      string `json:"secret"`
	PublicToken string `json:"public_token"`
}

type AccessToken struct {
	Token string `json:"access_token"`
}

type GetTransactions struct {
	ClientId    string `json:"client_id"`
	Secret      string `json:"secret"`
	AccessToken string `json:"access_token"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type TransactionList struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Amount float32 `json:"amount"`
	Date   string  `json:"date"`
	Name   string  `json:"name"`
}

func (p *PlaidAPI) GetTransactions() ([]Transaction, error) {
	response := getPublicToken()
	time.Sleep(2 * time.Second)
	access := getAccessToken(response)
	time.Sleep(2 * time.Second)
	TransactionList := getTransactions(access)

	return TransactionList.Transactions, nil
}

func getPublicToken() PublicToken {
	m := GetPublicToken{ClientId: os.Getenv("CLIENT_ID"), Secret: os.Getenv("SECRET"), InstitutionId: "ins_3", InitialProducts: []string{"auth", "transactions"}, Options: map[string]string{"webhook": "https://www.genericwebhookurl.com/webhook"}}

	b, err := json.Marshal(m)
	if err != nil {
		_ = fmt.Errorf("unexpected error: %w", err)
	}

	r, err := http.Post("https://sandbox.plaid.com/sandbox/public_token/create", "application/json", bytes.NewBuffer(b))
	if err != nil {
		_ = fmt.Errorf("unexpected error: %w", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_ = fmt.Errorf("unexpected error: %w", err)
	}

	var response PublicToken
	err = json.Unmarshal(data, &response)
	if err != nil {
		_ = fmt.Errorf("unexpected error: %w", err)
	}

	return response
}

func getAccessToken(response PublicToken) AccessToken {
	e := GetAccessToken{ClientId: os.Getenv("CLIENT_ID"), Secret: os.Getenv("SECRET"), PublicToken: response.Token}

	b, err := json.Marshal(e)
	if err != nil {
		_ = fmt.Errorf("unexpected error: %w", err)
	}

	r, err := http.Post("https://sandbox.plaid.com/item/public_token/exchange", "application/json", bytes.NewBuffer(b))
	if err != nil {
		_ = fmt.Errorf("unexpected error: %w", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_ = fmt.Errorf("unexpected error: %w", err)
	}

	var access AccessToken

	err = json.Unmarshal(data, &access)
	if err != nil {
		_ = fmt.Errorf("unexpected error: %w", err)
	}

	return access
}

func getTransactions(access AccessToken) TransactionList {
	t := GetTransactions{ClientId: os.Getenv("CLIENT_ID"), Secret: os.Getenv("SECRET"), AccessToken: access.Token, StartDate: "2018-11-10", EndDate: "2020-11-10"}

	b, err := json.Marshal(t)
	if err != nil {
		_ = fmt.Errorf("unexpected error: %w", err)
	}

	r, err := http.Post("https://sandbox.plaid.com/transactions/get", "application/json", bytes.NewBuffer(b))
	if err != nil {
		_ = fmt.Errorf("unexpected error: %w", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_ = fmt.Errorf("unexpected error: %w", err)
	}

	var listOfTransactions TransactionList

	err = json.Unmarshal(data, &listOfTransactions)
	if err != nil {
		_ = fmt.Errorf("unexpected error: %w", err)
	}
	return listOfTransactions
}
