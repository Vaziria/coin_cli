package xeggexlib_test

import (
	"testing"

	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/commonlib"
	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/xeggexlib"
	"github.com/stretchr/testify/assert"
)

func TestClientApi(t *testing.T) {

	base := commonlib.MockBaseLocation()
	client := xeggexlib.GetXeggexTestClient(t, base)
	client.Debug = true

	t.Run("testing getting balance", func(t *testing.T) {
		hasil, err := client.GetBalance()

		assert.Nil(t, err)

		osin := hasil.FindBalance("OSN")
		assert.NotNil(t, osin)

		t.Log(osin)
	})
}
