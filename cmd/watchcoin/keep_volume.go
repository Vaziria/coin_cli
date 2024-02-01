package watchcoin

import (
	"errors"
	"io"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/xeggexlib"
	"github.com/pdcgo/common_conf/pdc_common"
)

type RandomRange struct {
	Min int
	Max int
}

func (rnd *RandomRange) GetRamdom() int {
	return rand.Intn(rnd.Max-rnd.Min) + rnd.Min
}

type ConfigFake struct {
	OrderInterval  *RandomRange
	SizeMin        float64
	SizeMax        float64
	RangeTolerance float64
	MarketPair     []string
	Accounts       map[string]*xeggexlib.XeggexCredential
	Debug          bool
}

func (config *ConfigFake) GetRandomSize() float64 {
	res := config.SizeMin + rand.Float64()*(config.SizeMax-config.SizeMin)
	return math.Floor(res)
}

var ErrBestBidNotfound error = errors.New("best bid position not found")
var ErrToleranceGreater error = errors.New("tolerance greater")
var ErrLastPrice error = errors.New("lastprice error")

type DataMarket struct {
	*ConfigFake
	log io.Writer
}

func ExecuteOrder(clientPool *PoolAccount, data *DataMarket) error {
	client := clientPool.RandomClient()
	marketInfo, err := client.MarketInfo(data.MarketPair...)
	if err != nil {
		return err
	}

	log.Printf("last price %.8f\n", marketInfo.LastPrice)

	res, err := client.GetOrderBook(data.MarketPair...)
	if err != nil {
		return err
	}

	if marketInfo.LastPrice == 0 {
		return ErrLastPrice
	}

	// best, err := FindBestPrice(&res, data.RangeTolerance, marketInfo.LastPrice)
	best, err := FindBestPriceV2(&res)
	if err != nil {
		return err
	}

	if best == nil {
		return ErrBestBidNotfound
	}

	executedSize := data.GetRandomSize()

	// fmt.Fprintf(data.log, "found position at %.8f with quantity %.8f", best.Best, data.Size)
	log.Printf("[%s] found position at %.8f with quantity %.8f\n", client.Alias, best.Best, executedSize)
	err = client.TemporaryOrder(&xeggexlib.CreateOrderPayload{
		Symbol:   strings.Join(data.MarketPair, "/"),
		Side:     xeggexlib.SellSide,
		Type:     xeggexlib.LimitType,
		Quantity: executedSize,
		Price:    best.Best,
	}, func(order *xeggexlib.CreateOrderRes) error {

		time.Sleep(time.Second * 1)

		// res, err := client.GetOrderBook(data.MarketPair...)
		// if err != nil {
		// 	return err
		// }

		// openask, err := FindBidThreshold(&res, marketInfo.LastPrice, best.Best, data.RangeTolerance, executedSize)
		// if err != nil {
		// 	return err
		// }
		// if openask == nil {
		// 	return ErrToleranceGreater
		// }
		// if !openask.Found {
		// 	return ErrToleranceGreater
		// }

		clientbuy := clientPool.RandomClient()
		log.Printf("[%s] open buy position with quantity %.8f\n", clientbuy.Alias, executedSize)

		buyorder, err := clientbuy.PlaceOrder(&xeggexlib.CreateOrderPayload{
			Symbol:   strings.Join(data.MarketPair, "/"),
			Side:     xeggexlib.BuySide,
			Type:     xeggexlib.MarketType,
			Quantity: executedSize,
		})

		if err != nil {
			pdc_common.ReportError(err)
			return err
		}

		log.Printf("[%s] closed position with quantity %.8f with tolerance %.8f\n", clientbuy.Alias, executedSize, float64(0))

		// sleep for next order
		tInterval := data.OrderInterval.GetRamdom()
		time.Sleep(time.Duration(tInterval) * time.Second)
		log.Printf("sleep for %d seconds..\n", data.OrderInterval)

		// cancel jika nyantol
		_, err = clientbuy.CancelOrder(&xeggexlib.CancelOrderPayload{
			ID: buyorder.ID,
		})

		if err != nil {
			log.Printf("[%s] cancel error %s\n", clientbuy.Alias, err.Error())
			return err
		}

		return nil
	})

	if errors.Is(err, ErrToleranceGreater) {
		log.Println("tolerance greater than", data.RangeTolerance, "USDT")
		return nil
	}

	return err
}

type ErrorChan chan<- error

func FakeVolume(clientPool *PoolAccount, config *ConfigFake, logw io.Writer) (*DataMarket, error) {
	hasil := DataMarket{
		ConfigFake: config,
		log:        logw,
	}

	err := ExecuteOrder(clientPool, &hasil)

	if err != nil {
		return &hasil, err
	}

	return &hasil, nil
}

func ReverseAsks(res *xeggexlib.OrderBookRes) []*xeggexlib.Asks {
	count := len(res.Asks)
	asks := make([]*xeggexlib.Asks, count)

	for i := range res.Asks {
		// asks[count-i-1] = &res.Asks[i]
		asks[i] = &res.Asks[i]
	}

	return asks
}

type BestPrice struct {
	Best         float64
	Tolerance    float64
	PrevPrice    float64
	CurrentPrice float64
}

func FindBestPrice(res *xeggexlib.OrderBookRes, rangeTolerance float64, lastprice float64) (*BestPrice, error) {
	asks := ReverseAsks(res)

	best := BestPrice{
		Best:         0,
		Tolerance:    0,
		PrevPrice:    lastprice,
		CurrentPrice: 0,
	}

	for i, bid := range asks {
		if bid.Numberprice < lastprice {
			log.Printf("skip %d -- %.8f\n", i, bid.Numberprice)
			continue
		}

		best.CurrentPrice = asks[i].Numberprice
		if i > 0 {
			prevPrice := asks[i-1].Numberprice
			log.Printf("prev %.8f\n", best.PrevPrice)

			if prevPrice < lastprice {
				prevPrice = lastprice
			}
			best.PrevPrice = prevPrice
		}

		diff := (best.CurrentPrice - best.PrevPrice) / 2

		if diff > 0.000001 {
			best.Best = best.PrevPrice + 0.000001
			break
		}

		best.Tolerance += bid.PairQuantity

		if best.Tolerance > rangeTolerance {
			return nil, nil
		}

	}

	if best.Best == 0 {
		log.Println("getting best 0")
		return nil, nil
	}

	return &best, nil
}

type AskFound struct {
	Asks      *xeggexlib.Asks
	Tolerance float64
	Quantity  float64
	Found     bool
}

func FindBidThreshold(res *xeggexlib.OrderBookRes, lastprice float64, priceask float64, rangeTolerance float64, quantity float64) (*AskFound, error) {

	asks := ReverseAsks(res)

	foundask := AskFound{
		Tolerance: 0,
		Quantity:  0,
	}

	for _, b := range asks {
		ask := b
		log.Println("thresshold check", ask.Quantity, ask.Numberprice, priceask, quantity)
		if ask.Quantity == quantity && ask.Numberprice == priceask {
			foundask.Found = true
			foundask.Asks = ask
			foundask.Quantity += ask.Quantity
			break
		} else {
			if ask.Numberprice > lastprice {
				foundask.Tolerance += ask.PairQuantity
				foundask.Quantity += ask.Quantity
			}

		}

		if foundask.Tolerance > rangeTolerance {
			break
		}

	}

	return &foundask, nil
}
