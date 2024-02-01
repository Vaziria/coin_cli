package watchcoin_test

import (
	"log"
	"testing"

	"github.com/Vaziria/coin_cli/cmd/watchcoin"
)

type MockRpcOut struct {
}

func (stdout *MockRpcOut) Write(p []byte) (n int, err error) {
	message := string(p)
	log.Println("mock message", message)
	return len(p), err
}

func TestRpc(t *testing.T) {
	watchcoin.CallCommand(nil, &MockRpcOut{}, []string{"rpc", "help"})
}
