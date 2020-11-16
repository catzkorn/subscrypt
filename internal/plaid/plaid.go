package plaid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	Transactions []Transaction `json:"transaction"`
}

type Transaction struct {
	Amount       float32 `json:"amount"`
	Date         string `json:"date"`
	MerchantName string `json:"merchant_name"`
	Name         string `json:"name"`
}

func (p *PlaidAPI) GetTransactions() {
	m := Message{ClientId: "", Secret: "", InstitutionId: "ins_3", InitialProducts: []string{"auth","transactions"},Options: map[string]string{"webhook": "https://www.genericwebhookurl.com/webhook"} }
	b, _ := json.Marshal(m)
	fmt.Println(string(b))
	r , _ := http.Post("https://sandbox.plaid.com/sandbox/public_token/create", "application/json", bytes.NewBuffer(b))
	data, _ := ioutil.ReadAll(r.Body)
	var response Response
	err := json.Unmarshal(data, &response)

	if err != nil {
		println(err)
	}

	fmt.Println(response.PublicToken)

	e := Exchange{ClientId: m.ClientId, Secret: m.Secret, PublicToken: response.PublicToken}
	b, _ = json.Marshal(e)
	r , _ = http.Post("https://sandbox.plaid.com/item/public_token/exchange", "application/json", bytes.NewBuffer(b))
	data, _ = ioutil.ReadAll(r.Body)
	var access Access

	err = json.Unmarshal(data, &access)

	if err != nil {
		println(err)
	}

	fmt.Println(access.AccessToken)
	time.Sleep(10 * time.Second)

	getTransactions(m, access)
}

func getTransactions(m Message, access Access) {
	t := GetTransactions{ClientId: m.ClientId, Secret: m.Secret, AccessToken: access.AccessToken, StartDate: "2018-11-10", EndDate: "2020-11-10"}
	b, _ := json.Marshal(t)
	r, _ := http.Post("https://sandbox.plaid.com/transactions/get", "application/json", bytes.NewBuffer(b))
	data, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("%s\n", data)
	var listOfTransactions ListOfTransactions

	err := json.Unmarshal(data, &listOfTransactions)
	if err != nil {
		println(err)
	}

	fmt.Println(listOfTransactions)
}