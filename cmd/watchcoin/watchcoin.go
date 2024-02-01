package watchcoin

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Vaziria/coin_cli/lib/commonlib"
	"github.com/Vaziria/coin_cli/lib/xeggexlib"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type ExecutionLoop interface {
	SendErr(err error)
	Notify(alias string, sumbook *ThreesholdResult, marketinfo *xeggexlib.MarketInfoRes, cfg *FilterOrderConfig)
}

func RunLoopClient(
	config *WatchCoinConfig,
	wg *sync.WaitGroup,
	loop ExecutionLoop,
	cred *Accounts,
	marketpair []string,
	sleepdur time.Duration,
) error {
	defer wg.Done()

	report, err := NewQuantityReport("qtyreport.db")
	if err != nil {
		return err
	}

	fHandler := map[string]func(sumdata *SumBooks, totalbook *SumBooks) (bool, error){
		"PairThreeshold": func(sumdata *SumBooks, totalbook *SumBooks) (bool, error) {
			return totalbook.PairSumQuantity >= cred.PairTreeshold, nil
		},

		"PriceChangePercent": func(sumdata *SumBooks, totalbook *SumBooks) (bool, error) {
			return totalbook.Percent <= cred.PriceChangePercent, nil
		},
	}

	client := xeggexlib.NewXeggexClient(cred.Name, cred.ApiKey, cred.SecretKey)

	for {

		books, err := client.GetOrderBook(marketpair...)
		if err != nil {
			loop.SendErr(err)
			continue
		}

		marketinfo, err := client.MarketInfo("VISH", "USDT")
		if err != nil {
			loop.SendErr(err)
			continue
		}

		log.Println(cred.Name, "checking...")

		sumbook, debug, err := GetThreeshold(
			books,
			fHandler,
		)

		if err != nil {
			loop.SendErr(err)
			continue
		}

		if sumbook == nil {
			for key, cek := range debug {
				if !cek {
					message := fmt.Sprintf("[ %s ] orderbook masih kena sampai di batas %s", cred.Name, key)
					log.Println(message)
					break
				}
			}
			time.Sleep(sleepdur)
			continue
		}

		if config.BookDebug {

			msg, _ := sumbook.TotalOrderBook.WaMessage()
			msgdata := ""
			if sumbook.FoundSafe {
				msgdata, _ = sumbook.SafeForOrder.WaMessage()
			}

			log.Println("\n\norder \n\n" + msg + "\n\nsafe order\n\n" + msgdata + "\n---------\n\n")
		}

		if sumbook.FoundSafe {
			loop.Notify(cred.Name, sumbook, &marketinfo, cred.FilterOrderConfig)
			// log.Println(cred.Name, sumbook, &marketinfo, cred.FilterOrderConfig)
		}

		err = report.AddReport(sumbook.TotalOrderBook)

		if err != nil {
			loop.SendErr(err)
			continue
		}

		time.Sleep(sleepdur)
	}

}

func CreateWatchCoinScript() *cli.Command {
	return &cli.Command{

		Name:  "watchcoin",
		Usage: "watch coin and notify",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
			},
			&cli.BoolFlag{
				Name:    "notwatch",
				Aliases: []string{"nw"},
			},
		},

		Action: func(ctx *cli.Context) error {
			// initial service
			debug := ctx.Bool("debug")
			notwatch := ctx.Bool("notwatch")
			dirctx := commonlib.NewBaseLocation(ctx)
			confile := dirctx.Path("watch_coin_configuration.yaml")

			config := WatchCoinConfig{
				fname:     confile,
				BookDebug: debug,
				GroupName: "testing wa",
				Accounts: map[string]*Accounts{
					"default": {
						Name: "default",
						XeggexCredential: &xeggexlib.XeggexCredential{
							ApiKey:    "api_key",
							SecretKey: "secret_key",
						},
						FilterOrderConfig: &FilterOrderConfig{
							PairTreeshold:      1000,
							PriceChangePercent: 0.01,
						},
					},
				},
				MarketPair:  []string{"VISH", "USDT"},
				SleepSecond: 20,
			}

			if _, err := os.Stat(confile); errors.Is(err, os.ErrNotExist) {
				return config.Save()

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

			notifyFunc, err := CreateFuncNotify(dirctx, &config)

			if err != nil {
				return err
			}

			if notifyFunc == nil {
				return errors.New("callback notify nil")
			}

			excreport, cancel := NewExecutionReport(notifyFunc, config.SleepSecond)
			defer cancel()

			var wg sync.WaitGroup

			if notwatch {
				wg.Add(1)
			} else {

				// TODO: perlu direfactor
				for _, cred := range config.Accounts {
					wg.Add(1)
					go RunLoopClient(&config, &wg, excreport, cred, config.MarketPair, time.Second*time.Duration(config.SleepSecond))
				}
			}

			wg.Wait()

			return nil
		},
	}
}
