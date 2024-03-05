package walletcli

import (
	"sync"
	"testing"
	"time"
)

var serviceRunning sync.Once

func RunMockDaemonService(t *testing.T) *WalletCli {
	serviceRunning.Do(func() {

		RunServiceDaemon("D:/crypto/new_deployer/workspace/dist/bin/dashd.exe", "D:/sampah/dash")

		// assert.Nil(t, err)
	})

	cli := WalletCli{
		Host:     "http://localhost:14277/wallet/dev",
		Username: "virtuoso",
		Password: "virtuoso",
	}

	cli.WaitFullSync(time.Minute * 10)

	return &cli
}
