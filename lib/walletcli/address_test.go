package walletcli_test

import (
	"testing"

	"github.com/Vaziria/coin_cli/lib/walletcli"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {

	cli := walletcli.WalletCli{
		Host:     "http://127.0.0.1:9998/wallet/miningpool",
		Username: "virtuoso",
		Password: "virtuoso",
	}

	t.Run("test getting address list", func(t *testing.T) {

		hasil, err := cli.ListAddressBalances()
		t.Log(hasil)

		assert.Nil(t, err)
		assert.NotEmpty(t, hasil)

	})

	t.Run("get address by label", func(t *testing.T) {
		hasil, err := cli.GetAddressByLabel("")
		t.Log(hasil)

		assert.Nil(t, err)
		assert.NotEmpty(t, hasil)
	})
}
