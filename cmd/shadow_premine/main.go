package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Vaziria/coin_cli/lib/commonlib"
	"github.com/Vaziria/coin_cli/lib/walletcli"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type MinMax struct {
	Min int
	Max int
}

func (ra *MinMax) Get() int {
	n := ra.Min + rand.Intn(ra.Max-ra.Min+1)
	return n
}

type FakeMineConfig struct {
	AddressCount     *MinMax
	MiningPoolWallet string
	MinerWallet      string
	DaemonPath       string
	DirPath          string
	PaymentTime      int
	SetMining        bool
	MinPayout        float64
}

type FakeMineTransaction struct {
	Config        *FakeMineConfig
	Addresses     []string
	MinerCli      *walletcli.WalletCli
	MiningPoolCli *walletcli.WalletCli
}

func (fake *FakeMineTransaction) SendToAddresses() error {
	balances, err := fake.MiningPoolCli.ListAddressBalances()

	if err != nil {
		return err
	}

	var totalBalance float64 = 0

	for _, balance := range balances {
		totalBalance += balance
	}

	countaddr := len(fake.Addresses)

	totalBalance = 0.9 * totalBalance
	perbalance := totalBalance / float64(countaddr)
	if perbalance < fake.Config.MinPayout {
		log.Println("min payout tidak cukup, count addr=", countaddr, " per balance=", perbalance, " total balance=", totalBalance)
		return nil
	}
	log.Println("sending transaction, count addr=", countaddr, " per balance=", perbalance, " total balance=", totalBalance)

	payload := walletcli.SendManyPayload{}

	for _, item := range fake.Addresses {
		addr := item
		payload[addr] = float32(perbalance)
	}

	res, err := fake.MiningPoolCli.SendMany(payload)
	if err != nil {
		return err
	}

	log.Println("send to many", res)

	return nil
}

func (fake *FakeMineTransaction) InitializeAddress() error {
	label := "miner"
	addresses, err := fake.MinerCli.GetAddressByLabel(label)
	if err != nil {
		return err
	}

	c := 0
	rangec := fake.Config.AddressCount.Get()
	fake.Addresses = make([]string, rangec)

	for key := range addresses {
		if c >= rangec {
			break
		}

		address := key
		fake.Addresses[c] = address
		c += 1

	}

	for c < rangec {

		address, err := fake.MinerCli.NewAddress(label)
		if err != nil {
			return err
		}

		fake.Addresses[c] = address
		c += 1
		log.Println("miner creating address ", address)
	}

	return nil

}

func (fake *FakeMineTransaction) Run(ctx context.Context) error {

	tick := time.NewTicker(time.Duration(fake.Config.PaymentTime) * time.Second)
	defer tick.Stop()

	if fake.Config.SetMining {
		_, err := fake.MiningPoolCli.SetMining()

		if err != nil {
			return err
		}
		log.Println("running mining")
	}

Parent:
	for {
		select {
		case <-tick.C:
			err := fake.InitializeAddress()
			if err != nil {
				return err
			}

			err = fake.SendToAddresses()

			if err != nil {
				log.Println("error send many", err)
			}

		case <-ctx.Done():
			break Parent
		}

	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	app := &cli.App{
		Name:     "command line fake scenario",
		Usage:    "fake scenario untuk premine",
		Commands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "shadow-premine-config.yaml",
			},
		},
		Action: func(ctx *cli.Context) error {
			dirctx := commonlib.NewBaseLocation(ctx)

			config := FakeMineConfig{
				AddressCount: &MinMax{
					Min: 50,
					Max: 100,
				},
				MiningPoolWallet: "miningpool",
				MinerWallet:      "miner",
				DaemonPath:       "D:/testunifyroom/unfyd.exe",
				DirPath:          "D:/testunifyroom/datacoin",
				PaymentTime:      3600,
				SetMining:        false,
				MinPayout:        20,
			}

			confile := dirctx.Path(ctx.String("file"))

			// load configuration
			if _, err := os.Stat(confile); errors.Is(err, os.ErrNotExist) {
				file, err := os.OpenFile(confile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					return err
				}
				defer file.Close()
				return yaml.NewEncoder(file).Encode(&config)
			} else {
				rawdata, err := os.ReadFile(confile)
				if err != nil {
					return err
				}

				err = yaml.Unmarshal(rawdata, &config)
				if err != nil {
					return err
				}
			}

			// initialize fake transaction configuration
			miningpool := walletcli.WalletCli{
				Host:     "http://127.0.0.1:9998/wallet/miningpool",
				Username: "virtuoso",
				Password: "virtuoso",
			}

			miner := walletcli.WalletCli{
				Host:     "http://127.0.0.1:9998/wallet/miner",
				Username: "virtuoso",
				Password: "virtuoso",
			}

			runner := FakeMineTransaction{
				Config:        &config,
				MinerCli:      &miner,
				MiningPoolCli: &miningpool,
			}

			cancel, err := walletcli.RunServiceDaemon(config.DaemonPath, config.DirPath)
			if err != nil {
				return err
			}

			defer cancel()

			miningpool.WaitFullSync(time.Minute * 20)

			err = runner.Run(context.Background())

			return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		color.Red(err.Error())
	}
}
