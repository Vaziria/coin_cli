package xeggexlib_test

import (
	"testing"

	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/commonlib"
	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/xeggexlib"
	"github.com/stretchr/testify/assert"
)

func TestMarket(t *testing.T) {

	base := commonlib.MockBaseLocation()
	client := xeggexlib.GetXeggexTestClient(t, base)
	client.Debug = true

	t.Run("testing market info", func(t *testing.T) {
		hasil, err := client.MarketInfo("VISH", "USDT")

		assert.Nil(t, err)
		assert.NotEmpty(t, hasil)

		t.Log(hasil)
	})
}
