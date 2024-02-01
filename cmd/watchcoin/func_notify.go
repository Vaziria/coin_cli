package watchcoin

import (
	"context"
	"log"
	"strings"

	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/commonlib"
	"github.com/Vaziria/bitcoin_development_env/coin_cli/lib/wamanager"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func CreateFuncNotify(base *commonlib.BaseLocation, cfg *WatchCoinConfig) (func(message string) error, error) {
	store := wamanager.NewDeviceStore(base.Path("wa_device.json"))
	err := store.Load()
	if err != nil {
		return nil, err
	}

	mngr, err := wamanager.NewManager(
		&wamanager.ManagerConfiguration{
			Debug:  cfg.WaDebug,
			MeowDB: base.Path("wa_database.db"),
		},
		store,
	)

	if err != nil {
		return nil, err
	}

	device, err := store.GetFirst()
	if err != nil {
		return nil, err
	}

	if device == nil {
		device, err = mngr.AddDevice(context.Background(), base.Path("qr_data.jpg"))
		if err != nil {
			return nil, err
		}
	}

	var groupID types.JID

	client, err := mngr.ConnectDevice(device, nil)

	if err != nil {
		return nil, err
	}

	datas, err := client.Client.GetJoinedGroups()
	if err != nil {
		return nil, err
	}

	for _, data := range datas {
		if data.GroupName.Name == cfg.GroupName {
			groupID = data.JID

		}
	}

	// writer := WaStdOut{
	// 	GroupID: groupID,
	// 	Client:  client.Client,
	// 	Size:    0,
	// }

	writer, _ := NewWaStdOut(groupID, client.Client)

	client.Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			msg := v.Message.GetConversation()

			if msg == "ping" {
				message := "still online"
				client.Client.SendMessage(context.Background(), groupID, &proto.Message{
					Conversation: &message,
				})
			}

			actions := strings.Split(msg, " ")
			if actions[0] == "rpc" {

				log.Println("executing", actions)

				CallCommand(cfg, writer, actions)
			}
		}
	})

	return func(message string) error {
		_, err = client.Client.SendMessage(context.Background(), groupID, &proto.Message{
			Conversation: &message,
		})
		return err
	}, nil
}
