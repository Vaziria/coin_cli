package xeggexlib

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/pdcgo/common_conf/pdc_common"
	"github.com/sethvargo/go-retry"
)

type XeggexCredential struct {
	ApiKey    string
	SecretKey string
}

type XeggexClient struct {
	*XeggexCredential
	Alias  string
	client *http.Client
	Debug  bool
}

func NewXeggexClient(alias string, apikey string, secretkey string) *XeggexClient {
	return &XeggexClient{
		Alias: alias,
		XeggexCredential: &XeggexCredential{
			ApiKey:    apikey,
			SecretKey: secretkey,
		},
		client: CreateHttpClient(),
	}
}

func (client *XeggexClient) createReq(method string, path string, body io.Reader) (*http.Request, error) {

	uri := "https://api.xeggex.com/api/v2" + path
	req, err := http.NewRequest(method, uri, body)
	req.SetBasicAuth(client.ApiKey, client.SecretKey)

	headers := map[string]string{
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Content-Type":              "application/json",
		"Accept-Encoding":           "gzip",
		"Accept-Language":           "en,en-US;q=0.9,id-ID;q=0.8,id;q=0.7",
		"Cache-Control":             "no-cache",
		"Pragma":                    "no-cache",
		"Sec-Ch-Ua":                 `"Not A(Brand";v="99", "Google Chrome";v="121", "Chromium";v="121"`,
		"Sec-Ch-Ua-Mobile":          "?0",
		"Sec-Ch-Ua-Platform":        `"Windows"`,
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
		"Upgrade-Insecure-Requests": "1",
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	return req, err
}

type BalanceRes struct {
	Asset     string `json:"asset"`
	Name      string `json:"name"`
	Available string `json:"available"`
	Pending   string `json:"pending"`
	Held      string `json:"held"`
	Assetid   string `json:"assetid"`
}

type ListBalanceRes []BalanceRes

func (balances ListBalanceRes) FindBalance(ticker string) *BalanceRes {
	for _, balance := range balances {
		if balance.Asset == ticker {
			return &balance
		}
	}
	return nil
}

func (client *XeggexClient) GetBalance() (ListBalanceRes, error) {
	hasil := ListBalanceRes{}

	err := client.sendReq(&hasil, func() (*http.Request, error) {
		return client.createReq(http.MethodGet, "/balances", nil)
	})
	return hasil, err
}

func (client *XeggexClient) sendReq(hasil any, createreq func() (*http.Request, error)) error {

	b := retry.NewFibonacci(1 * time.Second)

	// Ensure the maximum total retry time is 5s.
	b = retry.WithMaxRetries(5, b)

	err := retry.Do(context.Background(), b, func(ctx context.Context) error {
		req, err := createreq()
		if err != nil {
			return err
		}

		res, err := client.client.Do(req)

		if err != nil {
			if client.Debug {
				pdc_common.ReportError(err)
			} else {
				log.Println("retry xeggex client..", err)
			}

			return retry.RetryableError(err)
		}

		switch res.StatusCode {
		case 400:
			return retry.RetryableError(errors.New("unknown error"))
		case 500:
			return retry.RetryableError(errors.New("bad server"))
		case 401:
			return retry.RetryableError(errors.New("account bermasalah"))
		}

		// Decompress the response body
		reader, err := gzip.NewReader(res.Body)
		if err != nil {
			data, _ := io.ReadAll(res.Body)
			log.Println(string(data), res.StatusCode)
			pdc_common.ReportError(err)
			return retry.RetryableError(err)
		}
		defer reader.Close()
		// defer res.Body.Close()

		if client.Debug {
			data, _ := io.ReadAll(reader)
			log.Println(string(data))
			err = json.Unmarshal(data, hasil)
		} else {
			err = json.NewDecoder(reader).Decode(hasil)
		}

		if err != nil {
			pdc_common.ReportError(err)
			return retry.RetryableError(err)
		}
		return nil
	})

	return err

}
