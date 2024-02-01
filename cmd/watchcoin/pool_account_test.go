package watchcoin_test

import (
	"testing"

	"github.com/Vaziria/coin_cli/cmd/watchcoin"
	"github.com/Vaziria/coin_cli/lib/xeggexlib"
	"github.com/stretchr/testify/assert"
)

func TestPoolClient(t *testing.T) {
	pool := watchcoin.NewPoolAccount()

	t.Run("testing add client", func(t *testing.T) {
		err := pool.AddClient("default", &xeggexlib.XeggexCredential{})
		assert.Nil(t, err)

		err = pool.AddClient("default", &xeggexlib.XeggexCredential{})
		assert.ErrorIs(t, err, watchcoin.ErrClientExist)
	})

	t.Run("testing get random client", func(t *testing.T) {

		for range [5]int{} {
			pool.RandomClient()
		}
	})

}
