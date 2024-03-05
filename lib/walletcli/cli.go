package walletcli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type WalletCli struct {
	Host     string      `yaml:"host"`
	Username string      `yaml:"username"`
	Password string      `yaml:"password"`
	Client   http.Client `yaml:"-"`
	Debug    bool        `yaml:"debug"`
}

func RunServiceDaemon(daemonName string, datadir string) (func() error, error) {
	// cmd := exec.Cmd{
	// 	Path: daemonName,
	// 	Args: []string{
	// 		daemonName,
	// 		fmt.Sprintf("-datadir=%s", datadir),
	// 		"-printtoconsole",
	// 	},
	// }

	cmd := exec.Command(daemonName, fmt.Sprintf("-datadir=%s", datadir), "-printtoconsole")
	// cmd := exec.Command("ping", "8.8.8.8", "-t")

	log.Println("run daemon", cmd.String())

	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags:    0x10,
		NoInheritHandles: true,
	}

	err := cmd.Start()
	if err != nil {
		return func() error { return nil }, err
	}

	return func() error {
		time.Sleep(time.Minute)
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
			return err
		}

		return nil
	}, nil
}

type RpcPayload struct {
	Jsonrpc    string `json:"jsonrpc"`
	ID         string `json:"id"`
	Method     string `json:"method"`
	Params     []any  `json:"params"`
	Walletname string `json:"walletname"`
}

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RpcRes struct {
	Error *RpcError `json:"error"`
	ID    string    `json:"id"`
}

func (cli *RpcRes) GetError() error {
	if cli.Error != nil {
		return errors.New(cli.Error.Message)
	}
	return nil
}

func (cli *WalletCli) createReq(method string, params ...any) (*http.Request, error) {
	credential := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cli.Username, cli.Password)))

	payload := RpcPayload{
		Jsonrpc: "1.0",
		ID:      uuid.New().String(),
		Method:  method,
		Params:  params,
	}
	// data, _ := json.Marshal(&payload)
	// log.Println("pauload", string(data))

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, cli.Host, &body)

	if err != nil {
		return req, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", credential))

	return req, nil
}

type ErrRes interface {
	GetError() error
}

func (cli *WalletCli) SendReq(hasil ErrRes, method string, params ...any) error {

	req, err := cli.createReq(method, params...)
	if err != nil {
		return err
	}

	res, err := cli.Client.Do(req)
	if err != nil {
		return err
	}

	if cli.Debug {
		data, _ := io.ReadAll(res.Body)
		log.Println(string(data))
		err = json.Unmarshal(data, hasil)
	} else {
		err = json.NewDecoder(res.Body).Decode(hasil)
	}

	if err != nil {
		return err
	}

	return hasil.GetError()
}

type BalanceRes struct {
	RpcRes
	Result float64 `json:"result"`
}

func (cli *WalletCli) Balances() (float64, error) {
	res := BalanceRes{}
	err := cli.SendReq(&res, "getbalance")
	return res.Result, err
}

type Address struct {
	Address       string   `json:"address"`
	Account       string   `json:"account"`
	Amount        float64  `json:"amount"`
	Confirmations int      `json:"confirmations"`
	Label         string   `json:"label"`
	Txids         []string `json:"txids"`
}

type GetAddressesRes struct {
	RpcRes
	Result []*Address `json:"result"`
}

func (cli *WalletCli) GetAddresses() ([]*Address, error) {
	res := GetAddressesRes{}
	err := cli.SendReq(&res, "listreceivedbyaddress", 0, true)

	return res.Result, err
}

type SendRes struct {
	RpcRes
	Result string `json:"result"`
}

type SendManyPayload map[string]float32

func (cli *WalletCli) SendMany(payload SendManyPayload) (string, error) {
	res := SendRes{}
	err := cli.SendReq(&res, "sendmany", "", payload)

	return res.Result, err
}
func (cli *WalletCli) Send(address string, amount float64) (string, error) {
	res := SendRes{}
	err := cli.SendReq(&res, "sendtoaddress", address, amount)

	return res.Result, err
}

func (cli *WalletCli) SendFrom(fromAccount string, address string, amount float64) (string, error) {
	res := SendRes{}

	err := cli.SendReq(&res, "sendfrom", fromAccount, address, amount)
	// if res.Result == "" {
	// 	return res.Result, errors.New("transcation tx empty")
	// }
	return res.Result, err
}

type Unspent struct {
	Txid          string  `json:"txid"`
	Vout          int     `json:"vout"`
	Address       string  `json:"address"`
	Account       string  `json:"account"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Confirmations int     `json:"confirmations"`
	PsRounds      int     `json:"ps_rounds"`
	Spendable     bool    `json:"spendable"`
	Solvable      bool    `json:"solvable"`
}

type ListUnspentRes struct {
	RpcRes
	Result []*Unspent `json:"result"`
}

func (cli *WalletCli) GetUnspent() ([]*Unspent, error) {
	res := ListUnspentRes{}
	err := cli.SendReq(&res, "listunspent")

	return res.Result, err
}

type BlockChainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int     `json:"blocks"`
	Headers              int     `json:"headers"`
	Bestblockhash        string  `json:"bestblockhash"`
	Difficulty           float64 `json:"difficulty"`
	Mediantime           int     `json:"mediantime"`
	Verificationprogress float64 `json:"verificationprogress"`
	Chainwork            string  `json:"chainwork"`
	Pruned               bool    `json:"pruned"`
	// Softforks            []Softforks     `json:"softforks"`
	// Bip9Softforks        []Bip9Softforks `json:"bip9_softforks"`
}
type Enforce struct {
	Status   bool `json:"status"`
	Found    int  `json:"found"`
	Required int  `json:"required"`
	Window   int  `json:"window"`
}
type Reject struct {
	Status   bool `json:"status"`
	Found    int  `json:"found"`
	Required int  `json:"required"`
	Window   int  `json:"window"`
}
type Softforks struct {
	ID      string  `json:"id"`
	Version int     `json:"version"`
	Enforce Enforce `json:"enforce"`
	Reject  Reject  `json:"reject"`
}
type Bip9Softforks struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type BlockChainInfoRes struct {
	RpcRes
	Result BlockChainInfo `json:"result"`
}

func (cli *WalletCli) GetBlockchainInfo() (BlockChainInfo, error) {
	res := BlockChainInfoRes{}
	err := cli.SendReq(&res, "getblockchaininfo")

	return res.Result, err
}

func (cli *WalletCli) WaitFullSync(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	time.Sleep(time.Second * 2)
Parent:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:

			info, err := cli.GetBlockchainInfo()
			if err != nil {

				switch err.Error() {
				case "Verifying wallet...":
					time.Sleep(time.Second * 2)
					continue Parent
				case "Loading wallet...":
					time.Sleep(time.Second * 2)
					continue Parent
				case "Loading block index...":
					time.Sleep(time.Second * 2)
					continue Parent
				case "Loading fulfilled requests cache...":
					time.Sleep(time.Second * 2)
					continue Parent
				default:
					log.Println("sync error", err)
					return err
				}

			}

			log.Println("service sync block", info.Blocks, info.Headers)
			if info.Blocks == 0 {
				time.Sleep(time.Second * 2)
				continue Parent
			}

			if info.Blocks == info.Headers {
				return nil
			}

			time.Sleep(time.Second * 2)
		}

	}

}
