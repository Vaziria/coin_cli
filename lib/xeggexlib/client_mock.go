package xeggexlib

import (
	"errors"
	"os"
	"testing"

	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/commonlib"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func GetXeggexTestClient(t *testing.T, base *commonlib.BaseLocation) *XeggexClient {
	data := XeggexCredential{
		ApiKey:    "",
		SecretKey: "",
	}

	fname := base.Path("xeggex_credential.yaml")
	t.Log(fname)

	if _, err := os.Stat(fname); errors.Is(err, os.ErrNotExist) {
		file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			assert.Nil(t, err)
			return nil
		}
		defer file.Close()

		err = yaml.NewEncoder(file).Encode(&data)
		assert.Nil(t, err)
		return nil
	}

	raw, err := os.ReadFile(fname)
	assert.Nil(t, err)

	err = yaml.Unmarshal(raw, &data)
	assert.Nil(t, err)

	return NewXeggexClient("testing", data.ApiKey, data.SecretKey)

}
