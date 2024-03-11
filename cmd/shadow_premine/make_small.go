package main

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Vaziria/coin_cli/lib/commonlib"
	"github.com/Vaziria/coin_cli/lib/walletcli"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type MinMax32 struct {
	Min float32
	Max float32
}

func (ra *MinMax32) Get() float32 {
	n := ra.Min + rand.Float32()*(ra.Max-ra.Min)
	return n
}

type MakeSmallConfig struct {
	WalletName string
	SleepTime  int
	Round      int
	AmountItem *MinMax32
	Amount     float32
	Addcount   int
	IsNewAddr  bool
	AddrLabel  string
	DaemonPath string
	DirPath    string
}

type MakeSmall struct {
	Config *MakeSmallConfig
	Wallet *walletcli.WalletCli
}

func (fake *MakeSmall) Run() error {
	amount := fake.Config.Amount
	count := fake.Config.Addcount

	c := 0
	for c < fake.Config.Round {
		err := fake.SendToMany(amount, count)
		if err != nil {
			return err
		}
		c += 1
		time.Sleep(time.Second * time.Duration(fake.Config.SleepTime))
	}

	return nil
}

func (fake *MakeSmall) SendToMany(amount float32, addrcount int) error {
	payload := walletcli.SendManyPayload{}

	amountSend := amount / float32(addrcount)

	addresses, err := fake.GetAddress(addrcount)

	if err != nil {
		return err
	}

	amo := float32(0)

	for _, addr := range addresses {

		if amo >= amount {
			break
		}

		sendaddr := addr
		amountSend := fake.Config.AmountItem.Get()
		amo += amountSend

		payload[sendaddr] = amountSend
	}

	_, err = fake.Wallet.SendMany(payload)
	if err != nil {
		return err
	}

	log.Println("sending to many", amountSend, "total", amount, "to", addrcount)

	return nil
}

func (fake *MakeSmall) GetAddress(count int) ([]string, error) {
	hasil := make([]string, count)

	c := 0
	label := fake.Config.AddrLabel

	if !fake.Config.IsNewAddr {

		addresses, err := fake.Wallet.GetAddressByLabel(label)
		if err != nil {
			return hasil, err
		}
		for key := range addresses {
			if c >= count {
				break
			}

			address := key
			hasil[c] = address
			c += 1

		}
	}

	for c < count {

		address, err := fake.Wallet.NewAddress(label)
		if err != nil {
			return hasil, err
		}

		hasil[c] = address
		c += 1
		log.Println("creating address ", address)
	}

	return hasil, nil

}

func SplitMoney() *cli.Command {
	return &cli.Command{
		Name:    "splitmoney",
		Aliases: []string{"spmoney"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "spfile",
				Aliases: []string{"spf"},
				Value:   "split-money-config.yaml",
			},
		},
		Action: func(ctx *cli.Context) error {

			dirctx := commonlib.NewBaseLocation(ctx)

			config := MakeSmallConfig{
				WalletName: "jj",
				SleepTime:  120,
				Round:      2,
				Amount:     15,
				AmountItem: &MinMax32{
					Min: 1,
					Max: 2,
				},
				Addcount:   3,
				IsNewAddr:  false,
				AddrLabel:  "split",
				DaemonPath: "D:/testunifyroom/unfyd.exe",
				DirPath:    "D:/testunifyroom/datacoin",
			}

			confile := dirctx.Path(ctx.String("spfile"))

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

			wallet := walletcli.WalletCli{
				Host:     "http://127.0.0.1:9998/wallet/" + config.WalletName,
				Username: "virtuoso",
				Password: "virtuoso",
			}

			cancel, err := walletcli.RunServiceDaemon(config.DaemonPath, config.DirPath)
			if err != nil {
				return err
			}

			defer cancel()

			wallet.WaitFullSync(time.Minute * 20)

			makesmall := MakeSmall{
				Config: &config,
				Wallet: &wallet,
			}

			return makesmall.Run()
		},
	}
}
