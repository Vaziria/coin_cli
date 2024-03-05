package walletcli

type MasternodeOutputsRes struct {
	RpcRes
	Result []string `json:"result"`
}

func (cli *WalletCli) MasternodeOutputs() (*MasternodeOutputsRes, error) {
	res := MasternodeOutputsRes{}
	err := cli.SendReq(&res, "masternode", "outputs")
	return &res, err
}

type BlsData struct {
	Secret string `json:"secret"`
	Public string `json:"public"`
	Scheme string `json:"scheme"`
}

type BlsGenerateRes struct {
	RpcRes
	Result *BlsData `json:"result"`
}

func (cli *WalletCli) BlsGenerate() (*BlsGenerateRes, error) {
	res := BlsGenerateRes{}
	err := cli.SendReq(&res, "bls", "generate")
	return &res, err
}
