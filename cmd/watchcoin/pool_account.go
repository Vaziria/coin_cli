package watchcoin

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/Vaziria/coin_cli/lib/xeggexlib"
)

var ErrClientExist error = errors.New("client exist")

type PoolAccount struct {
	sync.Mutex
	key      []string
	accounts map[string]*xeggexlib.XeggexClient
}

func NewPoolAccount() *PoolAccount {

	return &PoolAccount{
		accounts: map[string]*xeggexlib.XeggexClient{},
		key:      []string{},
	}
}

func (pool *PoolAccount) AddClient(key string, cred *xeggexlib.XeggexCredential) error {
	pool.Lock()
	defer pool.Unlock()

	if pool.accounts[key] != nil {
		return ErrClientExist
	}

	client := xeggexlib.NewXeggexClient(key, cred.ApiKey, cred.SecretKey)
	pool.accounts[key] = client
	pool.key = append(pool.key, key)

	return nil
}

func (pool *PoolAccount) RandomClient() *xeggexlib.XeggexClient {
	pool.Lock()
	defer pool.Unlock()

	in := rand.Intn(len(pool.key))
	client := pool.accounts[pool.key[in]]
	return client
}
