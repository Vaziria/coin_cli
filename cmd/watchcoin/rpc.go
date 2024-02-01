package watchcoin

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/xeggexlib"
	"github.com/dustin/go-humanize"
	"github.com/pdcgo/common_conf/pdc_common"
	"github.com/urfave/cli/v2"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"gopkg.in/yaml.v3"
)

type FilterOrderConfig struct {
	PairTreeshold      float64
	PriceChangePercent float64
}

type Accounts struct {
	*xeggexlib.XeggexCredential
	*FilterOrderConfig

	Name string
}

type WatchCoinConfig struct {
	BookDebug   bool
	fname       string
	WaDebug     bool
	GroupName   string
	Accounts    map[string]*Accounts
	SleepSecond int
	MarketPair  []string
}

func (cfg *WatchCoinConfig) Save() error {
	file, err := os.OpenFile(cfg.fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(cfg)
}

type WaStdOut struct {
	sync.Mutex

	Buffer  []byte
	Size    int
	GroupID types.JID
	Client  *whatsmeow.Client
}

func NewWaStdOut(
	GroupID types.JID,
	Client *whatsmeow.Client,
) (*WaStdOut, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	writer := WaStdOut{
		GroupID: GroupID,
		Client:  Client,
		Buffer:  []byte{},
		Size:    0,
	}

	go func() {
	Parent:
		for {
			select {
			case <-ctx.Done():
				break Parent
			default:

				if len(writer.Buffer) == 0 {
					time.Sleep(time.Second * 5)
					continue Parent
				}

				func() {
					writer.Lock()
					defer writer.Unlock()
					defer func() {
						writer.Size = 0
						writer.Buffer = []byte{}
					}()

					message := string(writer.Buffer)
					_, err := writer.Client.SendMessage(context.Background(), writer.GroupID, &proto.Message{
						Conversation: &message,
					})

					if err != nil {
						pdc_common.ReportError(err)
					}
				}()

				time.Sleep(time.Second * 5)
			}
		}
	}()

	return &writer, cancel

}

func (stdout *WaStdOut) Write(p []byte) (n int, err error) {
	stdout.Lock()
	defer stdout.Unlock()

	cp := len(p)

	stdout.Buffer = append(stdout.Buffer, p...)
	stdout.Size = cp
	return cp, nil
}

func CallCommand(config *WatchCoinConfig, writer io.Writer, action []string) {

	app := &cli.App{
		Name:      "rpc",
		Usage:     "command line toolkit watchcoin",
		Writer:    writer,
		ErrWriter: writer,
		Commands: []*cli.Command{
			CreateFakeVolumeScript(writer, false),
			{
				Name: "filtershow",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "alias",
						Aliases:     []string{"a"},
						DefaultText: "default",
						Value:       "default",
					},
				},

				Action: func(ctx *cli.Context) error {
					alias := ctx.String("alias")

					cfgaccount := config.Accounts[alias]

					if cfgaccount == nil {
						fmt.Fprintf(
							writer,
							"[*%s*] config kosong",
							alias,
						)
						return nil
					}

					fmt.Fprintf(
						writer,
						"[*%s*] jika *USDT* lebih *%s* dan *price impact* kurang *%.2f%%*",
						alias,
						humanize.Comma(int64(cfgaccount.PairTreeshold)),
						cfgaccount.PriceChangePercent*100,
					)

					return nil
				},
			},
			{
				Name: "filter",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "alias",
						Aliases:     []string{"a"},
						DefaultText: "default",
						Value:       "default",
					},
					&cli.Float64Flag{
						Name:    "quantity",
						Aliases: []string{"qty"},
					},
					&cli.Float64Flag{
						Name:    "percent",
						Aliases: []string{"prcnt"},
					},
				},
				Action: func(ctx *cli.Context) error {
					alias := ctx.String("alias")

					cfgaccount := config.Accounts[alias]

					if cfgaccount == nil {
						fmt.Fprintf(
							writer,
							"[*%s*] config kosong",
							alias,
						)
						return nil
					}

					qty := ctx.Float64("quantity")
					if qty != 0 {
						cfgaccount.PairTreeshold = qty
					}

					percent := ctx.Float64("percent")
					if percent != 0 {
						cfgaccount.PriceChangePercent = percent
					}

					err := config.Save()
					if err != nil {
						fmt.Fprintln(writer, err)
					}

					fmt.Fprintf(
						writer,
						"[*%s*] jika *USDT* lebih *%s* dan *price impact* kurang *%.2f%%* updated config",
						alias,
						humanize.Comma(int64(cfgaccount.PairTreeshold)),
						cfgaccount.PriceChangePercent*100,
					)

					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Fprintln(
				writer,
				"\"*rpc help*\" for help",
			)
			return nil
		},
	}

	err := app.Run(action)

	if err != nil {
		writer.Write([]byte(err.Error()))
	}
}
