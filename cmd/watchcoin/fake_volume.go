package watchcoin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Vaziria/coin_cli/lib/commonlib"
	"github.com/Vaziria/coin_cli/lib/xeggexlib"
	"github.com/pdcgo/common_conf/pdc_common"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var FakeVolumeRunning bool = false
var FakeVolumeCancel context.CancelFunc
var FakeVolumeLock sync.Mutex

func CreateFakeVolumeScript(writer io.Writer, wait bool) *cli.Command {

	errChan := make(chan error, 200)

	go func() {
		for err := range errChan {
			if errors.Is(err, ErrBestBidNotfound) {
				log.Println(err.Error())
				continue
			}

			if err != nil {
				pdc_common.ReportError(err)
			}

		}
	}()

	return &cli.Command{

		Name:  "fakevolume",
		Usage: "Fake Volume Trading",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "status",
				Aliases: []string{"s"},
			},

			&cli.BoolFlag{
				Name:    "run",
				Aliases: []string{"r"},
			},

			&cli.BoolFlag{
				Name:    "stop",
				Aliases: []string{"st"},
			},
		},

		Action: func(ctx *cli.Context) error {
			fname := "fake_volume_config.yaml"
			dirctx := commonlib.NewBaseLocation(ctx)
			confile := dirctx.Path(fname)
			config := ConfigFake{
				OrderInterval: &RandomRange{
					Min: 30,
					Max: 120,
				},
				MarketPair:     []string{"VISH", "USDT"},
				SizeMin:        1000,
				SizeMax:        2000,
				RangeTolerance: 1,
				Accounts: map[string]*xeggexlib.XeggexCredential{
					"akun_satu": {
						ApiKey:    "apikey",
						SecretKey: "secretkey",
					},
					"akun_dua": {
						ApiKey:    "apikey",
						SecretKey: "secretkey",
					},
				},
			}

			status := ctx.Bool("status")
			run := ctx.Bool("run")
			stop := ctx.Bool("stop")

			// initial service

			if _, err := os.Stat(confile); errors.Is(err, os.ErrNotExist) {
				return dirctx.SaveYaml(&config, "fake_volume_config.yaml")

			} else {
				// TODO: refactor load taruh di base location
				config.Accounts = map[string]*xeggexlib.XeggexCredential{}
				rawdata, err := os.ReadFile(confile)
				if err != nil {
					return err
				}

				err = yaml.Unmarshal(rawdata, &config)
				if err != nil {
					return err
				}
			}

			if status {
				_, err := fmt.Fprintf(writer, "fake volume running is %v\n", FakeVolumeRunning)
				return err
			}

			if run {
				FakeVolumeLock.Lock()
				defer FakeVolumeLock.Unlock()

				if FakeVolumeCancel != nil {
					FakeVolumeCancel()
				}
				FakeVolumeCancel = nil

				ctx, cancel := context.WithCancel(context.Background())

				go runLoopFakeVolume(ctx, writer, &config, errChan)

				FakeVolumeCancel = cancel

				FakeVolumeRunning = true

				_, err := fmt.Fprintf(writer, "fake volume running is %v\n", FakeVolumeRunning)

				if wait {
					for {
						time.Sleep(time.Hour)
					}
				}
				return err
			}

			if stop {
				func() {
					FakeVolumeLock.Lock()
					defer FakeVolumeLock.Unlock()

					if FakeVolumeCancel != nil {
						FakeVolumeCancel()
					}

					FakeVolumeCancel = nil
					FakeVolumeRunning = false
				}()

				_, err := fmt.Fprintf(writer, "fake volume running is %v\n", FakeVolumeRunning)
				return err
			}

			return nil
		},
	}
}

func runLoopFakeVolume(ctx context.Context, logw io.Writer, config *ConfigFake, errChan ErrorChan) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	clientPool := NewPoolAccount()
	for key, cred := range config.Accounts {
		clientPool.AddClient(key, cred)
	}

Parent:
	for {
		select {
		case <-ctx.Done():
			break Parent
		case <-ticker.C:
			_, err := FakeVolume(clientPool, config, logw)
			if err != nil {
				errChan <- err
			}

		}
	}
}
