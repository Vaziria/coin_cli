package commonlib

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type BaseLocation struct {
	base string
}

func (base *BaseLocation) Path(elem ...string) string {
	newelem := append([]string{base.base}, elem...)
	return filepath.Join(newelem...)
}

func (base *BaseLocation) SaveYaml(cfg any, elem ...string) error {
	fname := base.Path(elem...)
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(cfg)
}

func NewBaseLocation(ctx *cli.Context) *BaseLocation {

	base := ctx.String("base")
	return &BaseLocation{
		base: base,
	}
}

func MockBaseLocation(elem ...string) *BaseLocation {
	_, filename, _, _ := runtime.Caller(0)
	basedir := filepath.Dir(filename)
	listPath := []string{basedir, "../../test"}

	listPath = append(listPath, elem...)

	pathdata := filepath.Join(listPath...)
	os.MkdirAll(pathdata, os.ModeDir)

	return &BaseLocation{
		base: pathdata,
	}
}
