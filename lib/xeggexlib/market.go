package xeggexlib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MarketInfoRes struct {
	ID               string  `json:"id"`
	Symbol           string  `json:"symbol"`
	PrimaryAsset     string  `json:"primaryAsset"`
	SecondaryAsset   string  `json:"secondaryAsset"`
	LastPrice        float64 `json:"lastPrice"`
	HighPrice        string  `json:"highPrice"`
	LowPrice         string  `json:"lowPrice"`
	Volume           string  `json:"volume"`
	LineChart        string  `json:"lineChart"`
	LastTradeAt      int     `json:"lastTradeAt"`
	PriceDecimals    int     `json:"priceDecimals"`
	QuantityDecimals int     `json:"quantityDecimals"`
	IsActive         bool    `json:"isActive"`
}

func (b *MarketInfoRes) UnmarshalJSON(data []byte) error {

	type Alias MarketInfoRes
	aux := &struct {
		LastPrice string `json:"lastPrice"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	lastprice, err := strconv.ParseFloat(aux.LastPrice, 64)
	b.LastPrice = lastprice

	return err

}

func (client *XeggexClient) MarketInfo(pair ...string) (MarketInfoRes, error) {

	marketpair := strings.Join(pair, "_")

	uri := fmt.Sprintf("/market/getbysymbol/%s", marketpair)

	hasil := MarketInfoRes{}
	req, err := client.createReq(http.MethodGet, uri, nil)

	if err != nil {
		return hasil, err
	}

	err = client.sendReq(&hasil, func() (*http.Request, error) {
		return req, nil
	})
	return hasil, err
}
