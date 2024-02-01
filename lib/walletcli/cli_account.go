package walletcli

type AccountRes struct {
	RpcRes
	Result string `json:"result"`
}

func (cli *WalletCli) GetAccount(addr string) (string, error) {
	res := AccountRes{}
	err := cli.SendReq(&res, "getaccount", addr)

	return res.Result, err
}

func (cli *WalletCli) GetAccountAddress(walletname string) (string, error) {
	res := AccountRes{}
	err := cli.SendReq(&res, "getaccountaddress", walletname)

	return res.Result, err
}

func (cli *WalletCli) SetAccount(addr string, alias string) error {
	res := AccountRes{}
	err := cli.SendReq(&res, "setaccount", addr, alias)

	return err
}
