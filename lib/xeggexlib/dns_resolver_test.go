package xeggexlib_test

import (
	"net/http"
	"testing"

	"github.com/Vaziria/coin_cli/lib/xeggexlib"
	"github.com/stretchr/testify/assert"
)

func TestXeggexRequestApi(t *testing.T) {

	client := xeggexlib.CreateHttpClient()

	res, err := client.Get("https://api.xeggex.com/api/v2/market/getorderbookbysymbol/VISH_USDT")
	// res, err := client.Get("https://facebook.com/")

	assert.Nil(t, err)

	assert.Equal(t, res.StatusCode, http.StatusOK)

	// data, _ := io.ReadAll(res.Body)
	// t.Log(string(data))
}
