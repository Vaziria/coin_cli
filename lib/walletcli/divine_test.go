package walletcli_test

import (
	"testing"

	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/walletcli"
	"github.com/stretchr/testify/assert"
)

func TestDivi(t *testing.T) {
	t.Skip("skip")
	dat := walletcli.GetDivine(0.1, 0.3, 500)

	dd := float64(0)
	for _, d := range dat {
		dd += d
	}

	t.Log(dat)

	assert.LessOrEqual(t, dd, float64(1))
	assert.GreaterOrEqual(t, dd, float64(0.899999999))
}
