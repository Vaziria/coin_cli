package watchcoin_test

import (
	"testing"

	"github.com/Vaziria/bitcoin_development_env/coin_cli/cmd/watchcoin"
	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/commonlib"
	"github.com/stretchr/testify/assert"
)

func TestGettingNotifyCallback(t *testing.T) {
	t.Skip("skip")
	base := commonlib.MockBaseLocation("callback")
	funccallback, err := watchcoin.CreateFuncNotify(base, &watchcoin.WatchCoinConfig{
		GroupName: "testing wa",
	})

	assert.Nil(t, err)

	err = funccallback("testing creating manager function")

	assert.Nil(t, err)
}
