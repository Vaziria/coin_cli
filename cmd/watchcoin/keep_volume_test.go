package watchcoin_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/Vaziria/bitcoin_development_env/coin_cli/cmd/watchcoin"
	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/commonlib"
	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/xeggexlib"
	"github.com/stretchr/testify/assert"
)

func TestKeepVolume(t *testing.T) {
	base := commonlib.MockBaseLocation()

	// load data

	fname := base.Path("test_keep_volume.json")

	getOrderbook := func() xeggexlib.OrderBookRes {
		hasil := xeggexlib.OrderBookRes{}

		data, err := os.ReadFile(fname)
		assert.Nil(t, err)

		err = json.Unmarshal(data, &hasil)

		assert.Nil(t, err)

		return hasil
	}

	hasil := getOrderbook()

	t.Run("test find bid threshold", func(t *testing.T) {
		t.Run("masih tolerance skip 5 bid di bawah", func(t *testing.T) {
			qty, err := watchcoin.FindBidThreshold(&hasil, 0.002600, 0.002689, 4000, 1607.1029)

			qtyTolerance := float64(0)
			bids := watchcoin.ReverseAsks(&hasil)
			for i, bid := range bids[:6] {
				if i > 5 {
					break
				}

				qtyTolerance += bid.Quantity
				t.Log(i, qtyTolerance, bid.Quantity)

			}

			assert.Nil(t, err)
			assert.NotNil(t, qty.Asks)
			assert.Equal(t, qty.Tolerance, float64(69.2556356889))
			assert.Less(t, qty.Tolerance, float64(4000))
			assert.Equal(t, qtyTolerance, qty.Quantity)
			assert.True(t, qty.Found)

			t.Run("skip 5 bid dibawah dengan last price dibawah current", func(t *testing.T) {
				qty, err := watchcoin.FindBidThreshold(&hasil, 0.002688, 0.002689, 4000, 1607.1029)

				assert.Nil(t, err)
				assert.NotNil(t, qty.Asks)
				assert.Equal(t, qty.Tolerance, float64(0))
				assert.Equal(t, 1607.1029, qty.Quantity)
				assert.True(t, qty.Found)
			})
		})

		t.Run("tidak tolerance", func(t *testing.T) {
			qty, err := watchcoin.FindBidThreshold(&hasil, 0.002600, 0.002689, 1, 1607.1029)

			t.Log(qty.Tolerance)

			assert.Nil(t, err)
			assert.Nil(t, qty.Asks)
			assert.False(t, qty.Found)
		})

	})

	t.Run("test find best price", func(t *testing.T) {
		t.Run("find best true", func(t *testing.T) {
			best, err := watchcoin.FindBestPrice(&hasil, 0, 0.001798)

			assert.Nil(t, err)
			assert.NotEmpty(t, best)
			assert.Equal(t, best.Tolerance, float64(0))
			assert.NotEqual(t, float64(0), best.Best)

			t.Log(best.Best)

		})

		t.Run("find best false", func(t *testing.T) {
			bid := hasil.Bids[0]
			bid.Numberprice = 0.001799
			bid.Quantity = 1000000
			bid.PairQuantity = 0.001799 * bid.Quantity

			bids := append(hasil.Bids, bid)

			tamper := xeggexlib.OrderBookRes{
				Bids: bids,
			}

			best, err := watchcoin.FindBestPrice(&tamper, 0, 0.001799)

			assert.Nil(t, err)
			assert.Nil(t, best)
			// assert.Equal(t, best.Tolerance, float64(0))
		})

		t.Run("find best true ada bid price range under last price", func(t *testing.T) {
			ask1 := xeggexlib.Asks{
				Numberprice:  0.00001,
				Quantity:     100,
				PairQuantity: 100 * 0.00001,
			}

			ask2 := xeggexlib.Asks{
				Numberprice:  0.00008,
				Quantity:     100,
				PairQuantity: 100 * 0.00008,
			}
			asks := []xeggexlib.Asks{ask1, ask2}
			asks = append(asks, hasil.Asks...)

			tamper := xeggexlib.OrderBookRes{
				Asks: asks,
			}

			best, err := watchcoin.FindBestPrice(&tamper, 4000, 0.002509)

			assert.Nil(t, err)
			assert.NotEmpty(t, best)
			assert.Equal(t, best.Tolerance, float64(0))
			assert.Greater(t, best.Best, float64(0.001798))
			assert.NotEqual(t, float64(0), best.Best)
			t.Log(best.Best)

		})

		t.Run("find best false dengan ada bid price range under last price", func(t *testing.T) {
			bid1 := xeggexlib.Bids{
				Numberprice:  0.00001,
				Quantity:     100,
				PairQuantity: 100 * 0.00001,
			}

			bid2 := xeggexlib.Bids{
				Numberprice:  0.00008,
				Quantity:     100,
				PairQuantity: 100 * 0.00008,
			}

			bids := append(hasil.Bids, bid2, bid1)

			tamper := xeggexlib.OrderBookRes{
				Bids: bids,
			}

			best, err := watchcoin.FindBestPrice(&tamper, 1, 0.001799)
			assert.Nil(t, err)
			assert.Nil(t, best)
		})

	})

}
