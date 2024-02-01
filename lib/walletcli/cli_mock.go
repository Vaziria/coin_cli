package walletcli

import (
	"sync"
	"testing"
	"time"
)

var serviceRunning sync.Once

func RunMockDaemonService(t *testing.T) *WalletCli {
	serviceRunning.Do(func() {

		RunServiceDaemon("D:/testvish/vishaid.exe", "D:/testvish")

		// assert.Nil(t, err)
	})

	cli := WalletCli{
		Host:     "http://localhost:14277",
		Username: "virtuoso",
		Password: "virtuoso",
	}

	cli.WaitFullSync(time.Minute * 10)

	return &cli
}
