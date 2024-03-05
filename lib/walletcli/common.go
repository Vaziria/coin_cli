package walletcli

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sort"
	"time"
)

type DistributeToWallets struct {
	Config      *DistributeConfig
	WalletCount int
	Cli         *WalletCli
}

// bakalan deprecated
func (dis *DistributeToWallets) InitiateAddress() ([]*Address, error) {
	addresses, err := dis.Cli.GetAddresses()
	addrCount := len(addresses)
	if err != nil {
		return addresses, err
	}

	log.Println("address count", addrCount)

	if addrCount < dis.WalletCount {
		for c := addrCount; c < dis.WalletCount; c++ {
			addr, err := dis.Cli.NewAddress("")
			if err != nil {
				return addresses, err
			}
			log.Println("creating address", addr)
		}

		return dis.InitiateAddress()
	}

	return addresses, nil
}

func (dis *DistributeToWallets) DistributeAllWallet(ctx context.Context) error {
	rand.Seed(time.Now().Unix())
	// initiate wallet
	randAddr, err := NewAddrList(dis.Cli, dis.WalletCount)
	if err != nil {
		return err
	}

	balance, err := dis.Cli.Balances()
	if err != nil {
		return err
	}
	log.Println("balance unlocked", balance)
	if balance == 0 {
		return errors.New("balance unlocked 0")
	}

	ratios := GetDivine(0.1, 0.3, dis.WalletCount)

	for c := 0; c < dis.WalletCount; c++ {
		addr, err := randAddr.Get(c)
		if err != nil {
			return err
		}

		ratio := ratios[c]
		sendAmount := ratio * balance
		sendAmount = float64(int(sendAmount))
		tx, err := dis.Cli.Send(addr.Addr, sendAmount)

		if err != nil {
			return err
		}

		log.Println("send", sendAmount, "to", addr, "tx", tx)
		time.Sleep(time.Second * time.Duration(dis.Config.SendSleep))
	}

	time.Sleep(time.Minute * 10)

	return nil

}

func (dis *DistributeToWallets) Distribute() error {
	rand.Seed(time.Now().Unix())

	// ratios := GetDivine(0.1, 0.3, dis.WalletCount)

	for {

		log.Println("reloading address list")
		addrList, err := NewAddrList(dis.Cli, dis.WalletCount)
		if err != nil {
			return err
		}

		balance := addrList.Balance()
		ratio := (float64(1) + rand.Float64()) / float64(dis.WalletCount)

		bigaddrs := addrList.GetBigUnspent(balance * ratio)
		// log.Println("balance", balance, "big address list count", len(bigaddrs))

		for _, bigaddr := range bigaddrs {
			log.Println("big addr am", bigaddr.Amount, ratio)
			sratio := (0.2 + rand.Float64()*(0.3-0.2))
			sendAmount := bigaddr.Amount * sratio
			sendAmount = float64(int(sendAmount) - 1)

			hitlimit := false
			limitamount := float64(1000000)
			csend := 1
			if sendAmount > limitamount {
				sendAmount = 50000
				hitlimit = true
				csend = 5
			}

			for c := 0; c < csend; c++ {
				sendaddr, err := addrList.GetRandom()
				if err != nil {
					return err
				}

				accountName, err := bigaddr.AccountName(dis.Cli)
				if err != nil {
					return err
				}

				if bigaddr.Addr == sendaddr.Addr {
					continue
				}

				// tx, err := dis.Cli.SendFrom(accountName, sendaddr.Addr, sendAmount)
				tx, err := dis.Cli.Send(sendaddr.Addr, sendAmount)

				if err != nil {
					switch err.Error() {
					case "Account has insufficient funds":
						log.Println("err -----------", err)
						// time.Sleep(time.Second * 5)
					case "transcation tx empty":
						log.Println("err -----------", err)
					default:
						return err
					}
				} else {
					log.Println(sendAmount, "from", accountName, "to", sendaddr.Addr, "tx", tx)
				}
			}

			time.Sleep(time.Second * time.Duration(dis.Config.SendSleep))

			if hitlimit {
				break
			}
		}

		time.Sleep(time.Minute)

	}

}

func (dis *DistributeToWallets) GetBigUnspent(handler func(utxo *Unspent) error) error {

	utxos, err := dis.Cli.GetUnspent()
	if err != nil {
		return err
	}

	log.Println("getting utxo", len(utxos))

	sort.Slice(utxos[:], func(i, j int) bool {
		return utxos[i].Amount > utxos[j].Amount
	})

	// maxAmount := float64(0)

	for _, utxo := range utxos {
		if !utxo.Spendable {
			continue
		}
		if utxo.Amount > dis.Config.UnspentThreeshold {
			err := handler(utxo)
			if err != nil {
				return err
			}

		}
	}

	return nil

}

type DistributeConfig struct {
	UnspentThreeshold float64 `yaml:"unspent_threeshold"`
	SendPercent       float64 `yaml:"send_percent"`
	SendSleep         float32 `yaml:"send_sleep"`
}
