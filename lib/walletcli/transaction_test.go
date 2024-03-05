package walletcli_test

import (
	"testing"

	"github.com/Vaziria/coin_cli/lib/walletcli"
	"github.com/stretchr/testify/assert"
)

func TestTransactionBasic(t *testing.T) {
	cli := walletcli.RunMockDaemonService(t)

	res, err := cli.MasternodeOutputs()

	assert.Nil(t, err)
	assert.NotEmpty(t, res)

}
