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

func TestThreeshold(t *testing.T) {

	base := commonlib.MockBaseLocation()
	dataraw, err := os.ReadFile(base.Path("order_book.json"))

	assert.Nil(t, err)

	orderbooks := xeggexlib.OrderBookRes{}

	err = json.Unmarshal(dataraw, &orderbooks)
	assert.Nil(t, err)

	t.Run("test price Under", func(t *testing.T) {
		sumbook, _, err := watchcoin.GetThreeshold(
			orderbooks,
			map[string]func(sumdata *watchcoin.SumBooks, totalbook *watchcoin.SumBooks) (bool, error){
				"FilterThreeshold": func(sumdata *watchcoin.SumBooks, totalbook *watchcoin.SumBooks) (bool, error) {
					// t.Log(sumdata, sumdata.PairSumQuantity)

					if sumdata.PairSumQuantity >= 2000000 {
						return true, nil
					}

					return false, nil
				},

				"FilterPriceChangeUnder": func(sumdata *watchcoin.SumBooks, totalbook *watchcoin.SumBooks) (bool, error) {

					if sumdata.Percent <= 0.11 {
						return true, nil
					}

					return false, nil
				},
			},
		)

		assert.Nil(t, err)
		assert.False(t, sumbook.FoundSafe)

	})

	t.Run("test price Under", func(t *testing.T) {
		sumbook, _, err := watchcoin.GetThreeshold(
			orderbooks,
			map[string]func(sumdata *watchcoin.SumBooks, totalbook *watchcoin.SumBooks) (bool, error){
				"FilterThreeshold": func(sumdata *watchcoin.SumBooks, totalbook *watchcoin.SumBooks) (bool, error) {
					// t.Log(sumdata, sumdata.PairSumQuantity)

					if sumdata.PairSumQuantity >= 85 {
						return true, nil
					}

					return false, nil
				},

				"FilterPriceChangeUnder": func(sumdata *watchcoin.SumBooks, totalbook *watchcoin.SumBooks) (bool, error) {
					// t.Log("percent", sumdata.Percent, sumdata.PairSumQuantity)

					if sumdata.Percent <= 0.14 {
						return true, nil
					}

					return false, nil
				},
			},
		)

		assert.Nil(t, err)
		assert.NotEmpty(t, sumbook)
		assert.True(t, sumbook.FoundSafe)
		assert.GreaterOrEqual(t, sumbook.SafeForOrder.PairSumQuantity, float64(85))
		assert.LessOrEqual(t, sumbook.SafeForOrder.Percent, float64(0.14))

		msg, err := sumbook.SafeForOrder.WaMessage()
		assert.Nil(t, err)
		t.Log(msg)
	})

	t.Run("testing last price bener", func(t *testing.T) {

		book, _, err := watchcoin.GetThreeshold(xeggexlib.OrderBookRes{
			Bids: []xeggexlib.Bids{
				{
					Numberprice:  1,
					Quantity:     100,
					PairQuantity: 100,
				}, {
					Numberprice:  0.8,
					Quantity:     200,
					PairQuantity: 160,
				},
				{
					Numberprice:  0.01,
					Quantity:     1000,
					PairQuantity: 10,
				},
			},
		}, map[string]func(sumdata *watchcoin.SumBooks, totalbook *watchcoin.SumBooks) (bool, error){
			"FilterThreeshold": func(sumdata *watchcoin.SumBooks, totalbook *watchcoin.SumBooks) (bool, error) {
				return totalbook.PairSumQuantity >= 200, nil
			},

			"FilterPriceChangeUnder": func(sumdata *watchcoin.SumBooks, totalbook *watchcoin.SumBooks) (bool, error) {
				return totalbook.Percent <= 0.50, nil
			},
		})

		data, _ := json.MarshalIndent(book, "", "\t")
		t.Log(string(data))

		assert.Nil(t, err)
		assert.NotEmpty(t, book)
		assert.True(t, book.FoundSafe)
		assert.Equal(t, float64(260), book.SafeForOrder.PairSumQuantity)
		assert.Equal(t, float64(300), book.SafeForOrder.Quantity)

		assert.LessOrEqual(t, book.SafeForOrder.Percent, float64(0.5))
		assert.Equal(t, float64(200), book.SafeForOrder.LastBook.Quantity)

		assert.Equal(t, book.TotalOrderBook.PairSumQuantity, float64(270))

	})
}
