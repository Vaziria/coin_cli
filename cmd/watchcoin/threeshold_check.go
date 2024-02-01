package watchcoin

import (
	"errors"
	"fmt"
	"math"

	"github.com/Vaziria/coin_cli/lib/xeggexlib"
	"github.com/dustin/go-humanize"
)

type SumBooks struct {
	Quantity        float64
	PairSumQuantity float64
	FirstPrice      float64
	Percent         float64
	LastBook        *xeggexlib.Bids
}

func (data *SumBooks) WaMessage() (string, error) {
	// log.Println(data.LastBook.Price)

	message := `Quantity : *%s* Vish @%f
Nominal : *%s* Usdt
Impact : *-%.2f%%* at %f`

	hasil := fmt.Sprintf(
		message,
		humanize.Comma(int64(data.Quantity)),
		data.FirstPrice,
		humanize.Comma(int64(data.PairSumQuantity)),
		data.Percent*100,
		data.LastBook.Numberprice,
	)

	return hasil, nil
}

type ThreesholdResult struct {
	FoundSafe      bool
	SafeForOrder   *SumBooks
	TotalOrderBook *SumBooks
}

func GetThreeshold(databook xeggexlib.OrderBookRes, handlers map[string]func(sumdata *SumBooks, totalbook *SumBooks) (bool, error)) (*ThreesholdResult, map[string]bool, error) {

	sumBook := SumBooks{

		PairSumQuantity: 0,
		FirstPrice:      0,
		Percent:         0,
	}

	totalBook := SumBooks{

		PairSumQuantity: 0,
		FirstPrice:      0,
		Percent:         0,
	}

	hasil := ThreesholdResult{
		FoundSafe:      false,
		SafeForOrder:   &sumBook,
		TotalOrderBook: &totalBook,
	}

	checkres := map[string]bool{}

	// sort.Slice(databook.Bids, func(i, j int) bool {
	// 	return databook.Bids[i].Numberprice < databook.Bids[j].Numberprice
	// })

	for ind, b := range databook.Bids {
		book := b
		checkres = map[string]bool{}

		totalBook.PairSumQuantity += book.PairQuantity
		totalBook.Quantity += book.Quantity
		totalBook.Percent = math.Abs((book.Numberprice / totalBook.FirstPrice) - 1)
		totalBook.LastBook = &book

		if ind == 0 {
			sumBook.PairSumQuantity = book.PairQuantity
			sumBook.Quantity = book.Quantity
			sumBook.FirstPrice = book.Numberprice
			totalBook.FirstPrice = book.Numberprice

		} else {

			if !hasil.FoundSafe {
				sumBook.PairSumQuantity += book.PairQuantity
				sumBook.Quantity += book.Quantity
				sumBook.Percent = math.Abs((book.Numberprice / sumBook.FirstPrice) - 1)
			}
		}

		if sumBook.FirstPrice < book.Numberprice {
			return nil, checkres, errors.New("ask/ bid not in order")
		}

		useThis := true

		for key, handle := range handlers {
			cek, err := handle(&sumBook, &totalBook)
			if err != nil {
				return nil, checkres, err
			}

			checkres[key] = cek
			useThis = cek && useThis

		}

		if useThis && !hasil.FoundSafe {
			hasil.FoundSafe = true
			hasil.SafeForOrder = &sumBook
			sumBook.LastBook = &book
		}

	}

	return &hasil, checkres, nil

}
