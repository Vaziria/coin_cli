package main

import (
	"errors"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Vaziria/bitcoin_development_env/coin_cli/cmd/watchcoin"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func GettingWslInterface(name string) (net.IP, error) {
	ifaces, err := net.Interfaces()

	if err != nil {
		return nil, err
	}

	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		if i.Name != name {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip.IsPrivate() {
				return ip, nil
			}
		}
	}

	return nil, errors.New("getting wsl local ip")
}

func CheckXlaunch() error {
	out, err := exec.Command("tasklist").Output()

	if err != nil {
		return err
	}

	if !strings.Contains(string(out), "vcxsrv.exe") {
		return errors.New("xlaunch not installed or not running")
	}

	return nil
}

func CheckDockerInstalled() error {

	_, err := exec.Command("docker", "version").Output()

	if err != nil {
		return err
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ManifestName = "pdc_coin_manifest.yaml"

	app := &cli.App{
		Name:  "command line toolkit coin",
		Usage: "command line toolkit coin checking dependencies coin",
		Commands: []*cli.Command{
			watchcoin.CreateFakeVolumeScript(os.Stdout, true),
			CreateBuildScript(),
			CreateLaundryScript(),
			watchcoin.CreateWatchCoinScript(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "base",
				Aliases: []string{"b"},
				Value:   "./",
			},
		},
		Action: func(*cli.Context) error {
			checks := map[string]func() error{
				"checking xlaunch":          CheckXlaunch,
				"checking docker installed": CheckDockerInstalled,
				"checking configuration":    CheckOrInitConfiguration,
			}

			for key, fn := range checks {
				err := fn()

				if err != nil {
					color.Red(key)
					return err
				}

				color.Green(key)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		color.Red(err.Error())
	}

	// ip, err := GettingWslInterface("vEthernet (WSL)")

	// if err != nil {
	// 	log.Println("tidak bisa dapet ip local", err)
	// }

	// log.Println(ip)
}
