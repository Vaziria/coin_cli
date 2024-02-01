package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

var ManifestName string

type GithubCoin struct {
	Username string
	Name     string
	Token    string
}

type CoinConfig struct {
	Port           int        `yaml:"port"`
	Ticker         string     `yaml:"ticker"`
	GithubCoin     GithubCoin `yaml:"github_coin"`
	WorkCoinVolume []string   `yaml:"work_coin_volume"`

	// CoinDirectory string `yaml:"coin_directory"`

	// ansible_user: root
	// coin_port: 33333
	// coin_directory: .tonnagecore
	// fname_conf: tonnage.conf
	// service_name: tonnaged
}

func GetConfiguration() (*CoinConfig, error) {
	config := CoinConfig{
		Port:   33333,
		Ticker: "TNN",
		GithubCoin: GithubCoin{
			Username: "Vaziria",
			Name:     "tonnage",
		},
		WorkCoinVolume: []string{},
	}

	data, err := os.ReadFile(ManifestName)

	if err != nil {
		return &config, err
	}

	err = yaml.Unmarshal(data, &config)
	return &config, err

}

func CheckOrInitConfiguration() error {

	if _, err := os.Stat(ManifestName); errors.Is(err, os.ErrNotExist) {
		color.RedString("manifest %s not exist", ManifestName)

		err := InitConfiguration()

		if err != nil {
			return err
		}
	}

	return nil
}

func InitConfiguration() error {

	log.Println("initializing config", ManifestName)

	config := CoinConfig{
		Port:   33333,
		Ticker: "TNN",
		GithubCoin: GithubCoin{
			Username: "Vaziria",
			Name:     "tonnage",
		},
	}

	fmt.Print("Port Yang Digunakan Coin: ")
	fmt.Scan(&config.Port)
	fmt.Print("Username Github: ")
	fmt.Scan(&config.GithubCoin.Username)

	fmt.Print("Github Coin Name: ")
	fmt.Scan(&config.GithubCoin.Name)

	fmt.Print("Github Coin Token: ")
	fmt.Scan(&config.GithubCoin.Token)

	fmt.Print("Ticker: ")
	fmt.Scan(&config.Ticker)

	// create file
	af, err := os.OpenFile(ManifestName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer af.Close()

	err = yaml.NewEncoder(af).Encode(config)
	if err != nil {
		return err
	}

	return nil
}
