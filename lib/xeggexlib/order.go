package xeggexlib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type OrderBookRes struct {
	Marketid  string `json:"marketid"`
	Symbol    string `json:"symbol"`
	Timestamp int64  `json:"timestamp"`
	Bids      []Bids `json:"bids"`
	Asks      []Asks `json:"asks"`
}

type Bids struct {
	Price        string  `json:"price"`
	Numberprice  float64 `json:"numberprice"`
	Quantity     float64 `json:"quantity"`
	PairQuantity float64
}

func (b *Bids) UnmarshalJSON(data []byte) error {

	type Alias Bids
	aux := &struct {
		Quantity string `json:"quantity"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	qtx, err := strconv.ParseFloat(aux.Quantity, 64)
	b.Quantity = qtx
	b.PairQuantity = b.Quantity * b.Numberprice

	return err

}

type Asks struct {
	Price        string  `json:"price"`
	Numberprice  float64 `json:"numberprice"`
	Quantity     float64 `json:"quantity"`
	PairQuantity float64
}

func (b *Asks) UnmarshalJSON(data []byte) error {

	type Alias Asks
	aux := &struct {
		Quantity string `json:"quantity"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	qtx, err := strconv.ParseFloat(aux.Quantity, 64)
	b.Quantity = qtx
	b.PairQuantity = b.Quantity * b.Numberprice

	return err

}

func (client *XeggexClient) GetOrderBook(pair ...string) (OrderBookRes, error) {

	marketpair := strings.Join(pair, "_")

	uri := fmt.Sprintf("/market/getorderbookbysymbol/%s", marketpair)

	hasil := OrderBookRes{}
	req, err := client.createReq(http.MethodGet, uri, nil)

	if err != nil {
		return hasil, err
	}

	err = client.sendReq(&hasil, func() (*http.Request, error) {
		return req, nil
	})
	return hasil, err
}

// func (client *XeggexClient) GetOrderSnapshot(pair ...string) (OrderBookRes, error) {
// 	marketpair := strings.Join(pair, "_")

// 	hasil := OrderBookRes{}
// 	req, err := client.createReq(http.MethodGet, "/orders/snapshot", nil)

// 	if err != nil {
// 		return hasil, err
// 	}

// 	err = client.sendReq(&hasil, req)
// 	return hasil, err
// }
