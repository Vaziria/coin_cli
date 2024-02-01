package watchcoin_test

import (
	"testing"

	"github.com/Vaziria/coin_cli/cmd/watchcoin"
	"github.com/Vaziria/coin_cli/lib/commonlib"
	"github.com/Vaziria/coin_cli/lib/xeggexlib"
	"github.com/stretchr/testify/assert"
)

func TestQuantity(t *testing.T) {
	base := commonlib.MockBaseLocation()
	report, err := watchcoin.NewQuantityReport(base.Path("reportorder.db"))

	assert.Nil(t, err)

	err = report.AddReport(&watchcoin.SumBooks{
		Quantity:        10,
		PairSumQuantity: 10,
		FirstPrice:      10,
		Percent:         10,
		LastBook: &xeggexlib.Bids{
			Price:        "12",
			Numberprice:  12,
			Quantity:     12,
			PairQuantity: 12,
		},
	})

	assert.Nil(t, err)

}
