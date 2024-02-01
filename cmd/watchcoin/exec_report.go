package watchcoin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Vaziria/coin_cli/lib/xeggexlib"
	"github.com/dustin/go-humanize"
	"github.com/pdcgo/common_conf/pdc_common"
)

type ExecutionReport struct {
	sync.Mutex
	Sleep          int
	SendThreeshold int64
	LastSend       map[string]int64
	NotifyCallback func(message string) error
}

func NewExecutionReport(callback func(message string) error, sleep int) (*ExecutionReport, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.TODO())

	execr := ExecutionReport{
		LastSend:       map[string]int64{},
		SendThreeshold: 3 * 60,
		Sleep:          sleep,
		NotifyCallback: callback,
	}

	go func() {
	Parent:
		for {
			select {
			case <-ctx.Done():
				break Parent

			default:

				func() {
					execr.Lock()
					defer execr.Unlock()

					now := time.Now().Unix()
					for key, lastsend := range execr.LastSend {
						difftime := now - lastsend
						if difftime > execr.SendThreeshold {
							execr.LastSend[key] = 0
						}

					}
				}()

			}
		}

		time.Sleep(time.Second)
	}()

	return &execr, cancel

}

func (rep *ExecutionReport) SendErr(err error) {
	pdc_common.ReportError(err)
	time.Sleep(time.Second * time.Duration(rep.Sleep))
}

func (rep *ExecutionReport) Notify(alias string, sumbook *ThreesholdResult, marketinfo *xeggexlib.MarketInfoRes, cfg *FilterOrderConfig) {
	rep.Lock()
	defer rep.Unlock()

	if rep.LastSend[alias] != 0 {
		fmt.Printf("Waiting %d Second for sending to whatsapp\n", rep.SendThreeshold)
		return
	}

	// sending to wa
	// strdata, err := yaml.Marshal(sumbook)

	// if err != nil {
	// 	rep.SendErr(err)
	// 	return
	// }

	// filtercfg, err := yaml.Marshal(cfg)

	// if err != nil {
	// 	rep.SendErr(err)
	// 	return
	// }

	msgsum, _ := sumbook.TotalOrderBook.WaMessage()
	msgsum = fmt.Sprintf("*Order Book* [%s] :\n\nLast Price: *%.8f*\n", alias, marketinfo.LastPrice) + msgsum
	err := rep.NotifyCallback(msgsum)
	if err != nil {
		rep.SendErr(err)
		return
	}

	if sumbook.FoundSafe {
		msglimit, _ := sumbook.SafeForOrder.WaMessage()

		msgdlimit := fmt.Sprintf(
			"[%s] jika *USDT* lebih *%s* dan *price impact* kurang *%.2f%%*\n\n%s",
			alias,
			humanize.Comma(int64(cfg.PairTreeshold)),
			cfg.PriceChangePercent*100,
			msglimit,
		)

		err = rep.NotifyCallback(msgdlimit)
		if err != nil {
			rep.SendErr(err)
			return
		}
	}

	rep.LastSend[alias] = time.Now().Unix()
}
