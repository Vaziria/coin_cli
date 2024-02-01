package xeggexlib

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	goquery "github.com/google/go-querystring/query"
)

type OrderStatus string

const (
	OrderActive    OrderStatus = "active"
	OrderFilled    OrderStatus = "filled"
	OrderCancelled OrderStatus = "cancelled"
)

type AccountOrderQuery struct {
	Pair   string      `url:"pair"`
	Status OrderStatus `url:"status"`
	Limit  int         `url:"limit"`
	Skip   int         `url:"skip"`
}

type Order struct {
	ID                 string  `json:"id"`
	UserProvidedID     string  `json:"userProvidedId"`
	Market             *Market `json:"market"`
	Side               string  `json:"side"`
	Type               string  `json:"type"`
	Price              string  `json:"price"`
	Quantity           string  `json:"quantity"`
	ExecutedQuantity   string  `json:"executedQuantity"`
	RemainQuantity     string  `json:"remainQuantity"`
	RemainTotal        string  `json:"remainTotal"`
	RemainTotalWithFee string  `json:"remainTotalWithFee"`
	LastTradeAt        int     `json:"lastTradeAt"`
	Status             string  `json:"status"`
	IsActive           bool    `json:"isActive"`
	CreatedAt          int64   `json:"createdAt"`
	UpdatedAt          int64   `json:"updatedAt"`
}
type Market struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
}

func (client *XeggexClient) GetAccountOrder(query *AccountOrderQuery) ([]*Order, error) {

	hasil := []*Order{}

	v, err := goquery.Values(query)

	if err != nil {
		return nil, err
	}

	req, err := client.createReq(http.MethodGet, "/getorders", nil)
	req.URL.RawQuery = v.Encode()

	if err != nil {
		return hasil, err
	}

	err = client.sendReq(&hasil, func() (*http.Request, error) {
		return req, nil
	})
	return hasil, err
}

type Side string

const (
	SellSide Side = "sell"
	BuySide  Side = "buy"
)

type CreateOrderPayload struct {
	UserProvidedID string    `json:"userProvidedId"`
	Symbol         string    `json:"symbol"`
	Side           Side      `json:"side"`
	Type           OrderType `json:"type"`
	Quantity       float64   `json:"quantity"`
	Price          float64   `json:"price"`
	StrictValidate bool      `json:"strictValidate"`
}

func (pay *CreateOrderPayload) MarshalJSON() ([]byte, error) {
	type Alias CreateOrderPayload

	qty := strconv.FormatFloat(pay.Quantity, 'f', 6, 64)
	price := strconv.FormatFloat(pay.Price, 'f', 6, 64)

	return json.Marshal(&struct {
		Quantity string `json:"quantity"`
		Price    string `json:"price"`
		*Alias
	}{
		Quantity: qty,
		Price:    price,
		Alias:    (*Alias)(pay),
	})
}

// // Read implements io.Reader.
// func (pay *CreateOrderPayload) Read(p []byte) (n int, err error) {
// 	data, err := json.Marshal(pay)
// 	if err != nil {
// 		return 0, err
// 	}

// 	ndata := copy(p, data)
// 	return ndata, nil
// }

type OrderType string

const (
	MarketType OrderType = "market"
	LimitType  OrderType = "limit"
)

type CreateOrderRes struct {
	ID                 string `json:"id"`
	UserProvidedID     string `json:"userProvidedId"`
	Market             string `json:"market"`
	Side               string `json:"side"`
	Type               string `json:"type"`
	Price              string `json:"price"`
	Quantity           string `json:"quantity"`
	ExecutedQuantity   string `json:"executedQuantity"`
	RemainQuantity     string `json:"remainQuantity"`
	RemainTotal        string `json:"remainTotal"`
	RemainTotalWithFee string `json:"remainTotalWithFee"`
	LastTradeAt        int    `json:"lastTradeAt"`
	Status             string `json:"status"`
	IsActive           bool   `json:"isActive"`
	CreatedAt          int    `json:"createdAt"`
	UpdatedAt          int    `json:"updatedAt"`
}

func (client *XeggexClient) TemporaryOrder(payload *CreateOrderPayload, handler func(order *CreateOrderRes) error) error {
	order, err := client.PlaceOrder(payload)

	if err != nil {
		return err
	}

	err = handler(order)

	// log.Printf("[%s] cancel order %s\n", client.Alias, order.ID)
	_, errcancel := client.CancelOrder(&CancelOrderPayload{
		ID: order.ID,
	})
	if errcancel != nil {
		return errcancel
	}
	return err

}

func (client *XeggexClient) PlaceOrder(payload *CreateOrderPayload) (*CreateOrderRes, error) {

	hasil := CreateOrderRes{}

	data, err := json.Marshal(payload)
	if err != nil {
		return &hasil, err
	}

	err = client.sendReq(&hasil, func() (*http.Request, error) {
		return client.createReq(http.MethodPost, "/createorder", bytes.NewBuffer(data))

	})
	return &hasil, err
}

type CancelOrderPayload struct {
	ID string `json:"id"`
}

type CancelOrderRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

func (client *XeggexClient) CancelOrder(payload *CancelOrderPayload) (*CancelOrderRes, error) {

	hasil := CancelOrderRes{}

	data, err := json.Marshal(payload)
	if err != nil {
		return &hasil, err
	}

	err = client.sendReq(&hasil, func() (*http.Request, error) {
		return client.createReq(http.MethodPost, "/cancelorder", bytes.NewBuffer(data))
	})
	return &hasil, err
}

type CancelAllOrderPayload struct {
	Symbol string `json:"symbol"`
	Side   string `json:"side"`
}

type CancelAllOrderRes struct {
	Success bool     `json:"success"`
	Ids     []string `json:"ids"`
}
