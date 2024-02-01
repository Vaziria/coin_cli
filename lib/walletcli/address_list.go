package walletcli

import (
	"log"
	"math/rand"
	"sort"
	"sync"
)

var AccountPrefix = "wallet_"

type AddrItem struct {
	Addr   string
	Amount float64
	Alias  string
}

func (addr *AddrItem) AccountName(cli *WalletCli) (string, error) {
	if addr.Alias != "" {
		return addr.Addr, nil
	}

	alias, err := cli.GetAccount(addr.Addr)

	if err != nil {
		return "", err
	}

	if alias == "" {
		newalias := AccountPrefix + addr.Addr
		err := cli.SetAccount(addr.Addr, newalias)
		if err != nil {
			return "", err
		}

		alias = newalias

	}

	addr.Alias = alias
	return alias, nil
}

type AddrList struct {
	sync.Mutex
	cli         *WalletCli
	Data        []*AddrItem
	Index       map[string]*AddrItem
	WalletCount int
}

func NewAddrList(cli *WalletCli, count int) (*AddrList, error) {

	raddr := &AddrList{
		cli:   cli,
		Index: map[string]*AddrItem{},
	}
	addresses, err := cli.GetAddresses()
	addrlen := len(addresses)
	log.Println("address count", addrlen)
	if err != nil {
		return raddr, err
	}

	if count == -1 || count < addrlen {
		count = addrlen

	}

	raddr.Data = make([]*AddrItem, count)
	raddr.WalletCount = count

	for c, addr := range addresses {
		addro := AddrItem{
			Addr:   addr.Address,
			Amount: 0,
			Alias:  "",
		}

		raddr.Index[addr.Address] = &addro
		if c+1 > count {
			continue
		}

		raddr.Data[c] = &addro

	}

	err = raddr.InitiateBalanceAddress()
	if err != nil {
		return raddr, err
	}

	return raddr, nil
}

func (dis *AddrList) InitiateBalanceAddress() error {
	unspents, err := dis.cli.GetUnspent()

	if err != nil {
		return err
	}

	dis.Lock()
	defer dis.Unlock()

	for _, unspent := range unspents {
		if !unspent.Spendable {
			continue
		}

		if !unspent.Solvable {
			continue
		}

		addr := dis.Index[unspent.Address]
		if addr == nil {
			addr = &AddrItem{
				Addr:   unspent.Address,
				Amount: 0,
				Alias:  "",
			}

			dis.Index[unspent.Address] = addr
		}

		addr.Amount += unspent.Amount
	}

	return nil
}

func (dis *AddrList) GetRandom() (*AddrItem, error) {
	index := rand.Intn(dis.WalletCount)
	return dis.Get(index)
}

func (dis *AddrList) Get(index int) (*AddrItem, error) {
	addr := dis.Data[index]
	if addr == nil {
		dis.Lock()
		defer dis.Unlock()
		newaddr, err := dis.cli.NewAddress()
		if err != nil {
			return addr, err
		}

		addr = &AddrItem{
			Addr: newaddr,
		}

		dis.Data[index] = addr
		dis.Index[addr.Addr] = addr

		name, _ := addr.AccountName(dis.cli)

		log.Println("creating address", newaddr, "with alias", name)

	}
	return addr, nil
}

func (dis *AddrList) GetBigUnspent(thresshold float64) []*AddrItem {
	hasil := []*AddrItem{}
	for _, addr := range dis.Index {
		// log.Println(addr.Amount, thresshold)
		if addr.Amount > thresshold {
			hasil = append(hasil, addr)
		}
	}

	sort.Slice(hasil, func(i, j int) bool {

		return hasil[i].Amount > hasil[j].Amount
	})

	return hasil
}

func (dis *AddrList) Balance() float64 {
	hasil := float64(0)
	for _, addr := range dis.Index {
		hasil += addr.Amount
	}

	return hasil
}
