package xeggexlib_test

import (
	"encoding/json"
	"testing"

	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/commonlib"
	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/xeggexlib"
	"github.com/stretchr/testify/assert"
)

func TestParsingData(t *testing.T) {
	t.Run("test parse ask", func(t *testing.T) {
		askstr := `{"price":"0.004700","numberprice":0.0047,"quantity":"100.00000000"}`

		hasil := xeggexlib.Asks{}

		err := json.Unmarshal([]byte(askstr), &hasil)
		assert.Nil(t, err)
		assert.NotEmpty(t, hasil)
		assert.NotEmpty(t, hasil.Numberprice)
		assert.NotEmpty(t, hasil.PairQuantity)
		assert.Equal(t, hasil.Quantity, float64(100))

	})

	t.Run("test parse bid", func(t *testing.T) {
		askstr := `{"price":"0.004700","numberprice":0.0047,"quantity":"100.00000000"}`

		hasil := xeggexlib.Bids{}

		err := json.Unmarshal([]byte(askstr), &hasil)
		assert.Nil(t, err)
		assert.NotEmpty(t, hasil)
		assert.NotEmpty(t, hasil.Numberprice)
		assert.NotEmpty(t, hasil.PairQuantity)
		assert.Equal(t, hasil.Quantity, float64(100))

	})
}

func TestGetOrder(t *testing.T) {
	base := commonlib.MockBaseLocation()
	client := xeggexlib.GetXeggexTestClient(t, base)

	// client.Debug = true
	for i := range [50]int{} {
		_, err := client.GetOrderBook("OSN", "USDT")
		assert.Nil(t, err)

		t.Log(i)
	}

}
