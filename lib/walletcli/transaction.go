package walletcli

type Transaction struct {
	Amount              float64    `json:"amount"`
	Confirmations       int        `json:"confirmations"`
	Instantlock         bool       `json:"instantlock"`
	InstantlockInternal bool       `json:"instantlock_internal"`
	Chainlock           bool       `json:"chainlock"`
	Generated           bool       `json:"generated"`
	Blockhash           string     `json:"blockhash"`
	Blockheight         int        `json:"blockheight"`
	Blockindex          int        `json:"blockindex"`
	Blocktime           int        `json:"blocktime"`
	Txid                string     `json:"txid"`
	Time                int        `json:"time"`
	Timereceived        int        `json:"timereceived"`
	Details             []*Details `json:"details"`
	Hex                 string     `json:"hex"`
}
type Details struct {
	Address  string  `json:"address"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Label    string  `json:"label"`
	Vout     int     `json:"vout"`
}

type TransactionRes struct {
	RpcRes
	Result *Transaction `json:"result"`
}

func (cli *WalletCli) GetTransaction(addr string) (*TransactionRes, error) {
	res := TransactionRes{}
	err := cli.SendReq(&res, "gettransaction", addr)
	return &res, err
}
