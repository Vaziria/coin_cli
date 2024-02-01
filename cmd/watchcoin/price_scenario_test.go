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

func TestFindingDataBestPrice(t *testing.T) {
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

	t.Run("test find price normal", func(t *testing.T) {

		best, err := watchcoin.FindBestPriceV2(&hasil)
		assert.Nil(t, err)
		assert.NotEmpty(t, best)
		assert.Equal(t, float64(0.002591), best.Best)
		t.Log(best)

	})

	t.Run("test find price range tidak cukup", func(t *testing.T) {

		ask1 := xeggexlib.Asks{
			Numberprice: hasil.Bids[0].Numberprice + 0.000001,
		}

		asks := append([]xeggexlib.Asks{ask1}, hasil.Asks...)

		datas := xeggexlib.OrderBookRes{
			Bids: hasil.Bids,
			Asks: asks,
		}

		best, err := watchcoin.FindBestPriceV2(&datas)
		assert.Nil(t, err)
		assert.Nil(t, best)

		t.Log(best)
	})
}
