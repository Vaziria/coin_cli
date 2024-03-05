package walletcli

import "strings"

type ListAddressBalancesRes struct {
	RpcRes
	Result map[string]float64 `json:"result"`
}

func (cli *WalletCli) ListAddressBalances() (map[string]float64, error) {
	res := ListAddressBalancesRes{}
	err := cli.SendReq(&res, "listaddressbalances")

	return res.Result, err
}

type AddrLabel struct {
	Purpose string `json:"purpose"`
}

type GetAddressByLabelRes struct {
	RpcRes
	Result map[string]*AddrLabel `json:"result"`
}

func (cli *WalletCli) GetAddressByLabel(label string) (map[string]*AddrLabel, error) {
	res := GetAddressByLabelRes{
		Result: map[string]*AddrLabel{},
		RpcRes: RpcRes{
			Error: &RpcError{
				Message: "",
				Code:    0,
			},
		},
	}

	err := cli.SendReq(&res, "getaddressesbylabel", label)
	if err != nil {
		if strings.Contains(err.Error(), "No addresses with label") {
			return res.Result, nil
		}
	}

	return res.Result, err
}

type NewAddress struct {
	RpcRes
	Result string `json:"result"`
}

func (cli *WalletCli) NewAddress(alias string) (string, error) {
	res := NewAddress{}
	err := cli.SendReq(&res, "getnewaddress", alias)

	return res.Result, err
}

// 11:34:26
// listaddressbalances

// 11:34:26
// {
//   "UYmHdfLjmf6NYx9dx3jRzFz6ZAco7rVAuL": 540.00000000
// }

// 11:35:29
// getaddressesbylabel "miner"

// 11:35:29
// {
//   "UUmFA4uEhdNzjHzjugBjkX3yFxszW3woNo": {
//     "purpose": "receive"
//   },
//   "UVDwQGdVrGaZQ2bovD25UDPdp7KzDgqAMw": {
//     "purpose": "receive"
//   },
//   "Uaz8P1MjZKxtnuoz4LxQkuBawmMv1yuLue": {
//     "purpose": "receive"
//   },
//   "UcTAa6TXVfN1FEsjVEadY8jwdTMcPf4eWo": {
//     "purpose": "receive"
//   }
// }
