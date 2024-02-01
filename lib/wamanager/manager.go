package wamanager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type ManagerConfiguration struct {
	Debug  bool
	MeowDB string
}

type Manager struct {
	Config    *ManagerConfiguration
	container *sqlstore.Container
	db        *DeviceStore
}

func NewManager(config *ManagerConfiguration, store *DeviceStore) (*Manager, error) {

	manager := Manager{
		Config: config,
		db:     store,
	}

	dbLog := waLog.Stdout("Database", manager.GetLevelLog(), true)
	meowUri := fmt.Sprintf("file:%s?_foreign_keys=on", config.MeowDB)
	container, err := sqlstore.New("sqlite3", meowUri, dbLog)

	if err != nil {
		return nil, err
	}

	manager.container = container

	return &manager, nil
}

func (mngr *Manager) GetLevelLog() string {
	if mngr.Config.Debug {
		return "DEBUG"
	}
	return "INFO"
}
func (mngr *Manager) AddDevice(ctx context.Context, qrfname string) (*WaDevice, error) {
	var err error
	var device *WaDevice

	deviceStore := mngr.container.NewDevice()

	clientLog := waLog.Stdout("Client", mngr.GetLevelLog(), true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	successChan := make(chan error, 1)

	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.AppStateSyncComplete:
			dev := WaDevice{
				Phone:  client.Store.ID.User,
				Client: client,
			}
			err := mngr.db.Save(&dev)
			log.Println("completed sync", client.Store.ID.User)

			device = &dev

			successChan <- err
		case *events.Message:
			fmt.Println("Received a message!", v.Message.GetConversation())
		}
	})

	go func() {

		if client.Store.ID == nil {
			// No ID stored, new login
			qrChan, _ := client.GetQRChannel(ctx)
			err = client.Connect()
			if err != nil {
				successChan <- err
				return
			}
			for evt := range qrChan {
				if evt.Event == "code" {
					// Render the QR code here
					// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
					fmt.Println("QR code:", evt.Code)
					err = mngr.GenerateQr(evt.Code, qrfname)

					if err != nil {
						successChan <- err
					}

					mngr.OpenBrowser(qrfname)
					if err != nil {
						successChan <- err
					}
				} else {
					fmt.Println("Login event:", evt.Event, client.Store.ID.ADString())

				}
			}
		} else {
			// Already logged in, just connect
			err = client.Connect()
			if err != nil {
				successChan <- err
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		return device, ctx.Err()
	case errauth := <-successChan:
		client.Disconnect()
		if errauth != nil {
			return device, errauth
		}
		return device, nil
	}

}

func (mngr *Manager) GenerateQr(url string, filename string) error {

	qrCode, _ := qrcode.New(url, qrcode.Medium)
	err := qrCode.WriteFile(256, filename)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("QR code generated and saved as %s", filename)

	return nil
}

func (mngr *Manager) OpenBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return err
	}

	return nil

}

func (mngr *Manager) ConnectDevice(device *WaDevice, msgChan chan *events.Message) (*WaDevice, error) {

	datas, err := mngr.container.GetAllDevices()
	if err != nil {
		return nil, err
	}

	var deviceStore *store.Device

	for _, data := range datas {
		if data.ID.User == device.Phone {
			deviceStore = data
			break
		}
	}

	if deviceStore == nil {
		return nil, errors.New("phone not found")
	}

	clientLog := waLog.Stdout("Client", mngr.GetLevelLog(), true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.AppStateSyncComplete:
			log.Println("completed sync", client.Store.ID.User)
		case *events.Message:
			msg := v.Message.GetConversation()
			fmt.Println("Received a message!", msg)
		}
	})

	err = client.Connect()

	device.Client = client
	return device, err
}
