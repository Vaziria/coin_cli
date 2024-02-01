package main

import (
	"embed"
	"log"
	"os"
	"strings"

	"github.com/Vaziria/coin_cli/lib/commonlib"
	"github.com/urfave/cli/v2"
)

//go:embed templates
var Templates embed.FS

func CreateBuildScript() *cli.Command {

	return &cli.Command{
		Name:  "generate",
		Usage: "generate script tool",

		Action: func(ctx *cli.Context) error {
			dirctx := commonlib.NewBaseLocation(ctx)

			initDirectory(dirctx)

			config, err := GetConfiguration()

			if err != nil {
				return err
			}

			scripts := []string{"buildwin.sh", "buildunix.sh"}
			for _, name := range scripts {
				err := createScript(config, dirctx, name)
				if err != nil {
					return err
				}
			}

			err = GenerateCompose(ctx)
			if err != nil {
				log.Println("Generating compose gagal")
				return err
			}

			log.Println("Generating script success")

			return nil
		},
	}
}

func createScript(config *CoinConfig, dirctx *commonlib.BaseLocation, name string) error {

	data, err := Templates.ReadFile("templates/build/" + name)
	if err != nil {
		return err
	}

	content := strings.ReplaceAll(string(data), "\r", "")
	content = strings.ReplaceAll(content, "COIN_NAME", config.GithubCoin.Name)

	fname := dirctx.Path("tool", name)
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)

	return err
}

func initDirectory(dirctx *commonlib.BaseLocation) {
	os.MkdirAll(dirctx.Path("tool", "unixdist"), os.ModeDir)
	os.MkdirAll(dirctx.Path("tool", "windist"), os.ModeDir)
}
