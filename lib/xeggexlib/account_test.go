package xeggexlib_test

import (
	"testing"

	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/commonlib"
	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/xeggexlib"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	base := commonlib.MockBaseLocation()
	client := xeggexlib.GetXeggexTestClient(t, base)

	client.Debug = true

	t.Run("test create order", func(t *testing.T) {
		data, err := client.PlaceOrder(&xeggexlib.CreateOrderPayload{
			Symbol:   "VISH/USDT",
			Side:     xeggexlib.SellSide,
			Type:     xeggexlib.LimitType,
			Quantity: 1.3,
			Price:    0.1,
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, data)

		t.Run("test akun order", func(t *testing.T) {

			_, err := client.GetAccountOrder(&xeggexlib.AccountOrderQuery{
				Pair:   "VISH/USDT",
				Status: xeggexlib.OrderActive,
				Limit:  10,
				Skip:   0,
			})

			assert.Nil(t, err)
			assert.NotEmpty(t, data)

		})

		t.Run("test cancel order", func(t *testing.T) {

			ord, err := client.CancelOrder(
				&xeggexlib.CancelOrderPayload{
					ID: data.ID,
				},
			)

			assert.Nil(t, err)
			assert.NotEmpty(t, ord)

		})
	})

	// t.Run("test create order market", func(t *testing.T) {
	// 	data, err := client.PlaceOrder(&xeggexlib.CreateOrderPayload{
	// 		Symbol:   "VISH/USDT",
	// 		Side:     xeggexlib.SellSide,
	// 		Type:     "market",
	// 		Quantity: "1",
	// 		// Price:    "0.1",
	// 	})

	// 	assert.Nil(t, err)
	// 	assert.NotEmpty(t, data)

	// })
}
