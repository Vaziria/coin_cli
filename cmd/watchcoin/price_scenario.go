package watchcoin

import "github.com/Vaziria/bitcoin_development_env/coin_cli/lib/xeggexlib"

func FindBestPriceV2(res *xeggexlib.OrderBookRes) (*BestPrice, error) {
	highbid := res.Bids[0]
	lowask := res.Asks[0]

	hasil := BestPrice{
		PrevPrice:    highbid.Numberprice,
		CurrentPrice: lowask.Numberprice,
	}

	diff := (lowask.Numberprice - highbid.Numberprice) / 2

	if diff > 0.000001 {
		hasil.Best = highbid.Numberprice + 0.000001
		return &hasil, nil
	}

	return nil, nil

}
