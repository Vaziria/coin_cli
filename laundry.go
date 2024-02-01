package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Vaziria/coin_cli/lib/commonlib"
	"github.com/Vaziria/coin_cli/lib/walletcli"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type LaundryConfig struct {
	DaemonPath  string                     `yaml:"daemon_path"`
	DirPath     string                     `yaml:"dir_path"`
	Cli         walletcli.WalletCli        `yaml:"cli"`
	Distribute  walletcli.DistributeConfig `yaml:"distribute"`
	WalletCount int                        `yaml:"wallet_count"`
}

func CreateLaundryScript() *cli.Command {

	return &cli.Command{
		Name:  "laundry",
		Usage: "generate script tool",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "laundry-config.yaml",
			},
			&cli.StringFlag{
				Name:    "scenario",
				Aliases: []string{"sc"},
				Value:   "",
			},
		},

		Action: func(ctx *cli.Context) error {
			dirctx := commonlib.NewBaseLocation(ctx)

			config := LaundryConfig{
				WalletCount: 500,
				DaemonPath:  "D:/testvish/vishaid.exe",
				DirPath:     "D:/testvish2",
				Cli: walletcli.WalletCli{
					Host:     "http://localhost:14277",
					Username: "virtuoso",
					Password: "virtuoso",
				},
				Distribute: walletcli.DistributeConfig{
					UnspentThreeshold: 200,
					SendPercent:       0.6,
					SendSleep:         0.1,
				},
			}

			confile := dirctx.Path(ctx.String("file"))

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

			cancel, err := walletcli.RunServiceDaemon(config.DaemonPath, config.DirPath)
			if err != nil {
				return err
			}

			defer cancel()

			cli := config.Cli
			cli.WaitFullSync(time.Minute * 20)

			dis := walletcli.DistributeToWallets{
				Config:      &config.Distribute,
				WalletCount: config.WalletCount,
				Cli:         &cli,
			}

			scenario := ctx.String("scenario")
			switch scenario {
			case "allwallet":
				err = dis.DistributeAllWallet(context.Background())
			default:
				err = dis.Distribute()
			}

			if err != nil {
				log.Println(err)
				return err
			}

			return nil
		},
	}
}
