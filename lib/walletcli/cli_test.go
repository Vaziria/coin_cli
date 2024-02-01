package walletcli_test

import (
	"testing"

	"github.com/Vaziria/coin_cli/lib/walletcli"
	"github.com/stretchr/testify/assert"
)

func TestCli(t *testing.T) {
	t.Skip("skip")
	cli := walletcli.RunMockDaemonService(t)

	t.Run("test getting balance", func(t *testing.T) {
		balance, err := cli.Balances()

		t.Log(balance)

		assert.Nil(t, err)
		assert.Greater(t, balance, float64(0))

	})

	// var addr string
	// t.Run("test getting new address", func(t *testing.T) {
	// 	var err error
	// 	addr, err = cli.NewAddress()
	// 	t.Log(addr)
	// 	assert.Nil(t, err)
	// 	assert.NotEmpty(t, addr)

	// })

	t.Run("testing getting list address", func(t *testing.T) {
		t.Skip("skip")

		hasil, err := cli.GetAddresses()
		// t.Log(hasil)

		assert.Nil(t, err)
		assert.NotEmpty(t, hasil)

		assert.Contains(t, hasil, "uqxrxtAhct9Pq7xZRvQ7LcefZXwuo6Ehwg")

	})

	t.Run("getting unspent", func(t *testing.T) {
		t.Skip("skip")

		hasil, err := cli.GetUnspent()
		// t.Log(len(hasil))

		assert.Nil(t, err)
		assert.NotEmpty(t, hasil)

		// for _, unspent := range hasil {
		// 	if unspent.Address == "uqxrxtAhct9Pq7xZRvQ7LcefZXwuo6Ehwg" {
		// 		t.Log(unspent.Amount)
		// 	}
		// }
	})
}
