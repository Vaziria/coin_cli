package walletcli

import (
	"errors"
)

type SetMiningRes struct {
	RpcRes
	Result map[string]string `json:"result"`
}

func (cli *WalletCli) SetMining() (map[string]string, error) {
	res := SetMiningRes{}

	address, err := cli.GetAddressByLabel("")
	if err != nil {
		return res.Result, err
	}

	for key := range address {

		err := cli.SendReq(&res, "setgenerate", true, 1, key)
		return res.Result, err
	}

	return res.Result, errors.New("address empty")
}
