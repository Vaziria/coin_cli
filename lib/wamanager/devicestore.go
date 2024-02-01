package wamanager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"gopkg.in/yaml.v3"
)

type WaDevice struct {
	Phone  string
	Client *whatsmeow.Client `json:"-" yaml:"-"`
}

func (wa *WaDevice) SendMessage(ctx context.Context, phone string, data string) error {
	jidstr := fmt.Sprintf("%s@s.whatsapp.net", phone)
	jid, err := types.ParseJID(jidstr)
	if err != nil {
		return err
	}

	_, err = wa.Client.SendMessage(ctx, jid, &proto.Message{
		Conversation: &data,
	})

	if err != nil {
		return err
	}

	return nil
}

func (dev *WaDevice) GetID() (types.JID, error) {
	jidstr := fmt.Sprintf("%s@s.whatsapp.net", dev.Phone)
	return types.ParseJID(jidstr)
}

type DeviceStore struct {
	sync.Mutex
	Path string
	Data map[string]*WaDevice
}

func NewDeviceStore(path string) *DeviceStore {
	store := DeviceStore{
		Path: path,
		Data: map[string]*WaDevice{},
	}

	return &store
}

func (store *DeviceStore) Load() error {
	if _, err := os.Stat(store.Path); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	data, err := os.ReadFile(store.Path)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &store.Data)

}

func (store *DeviceStore) Get(phone string) (*WaDevice, error) {
	store.Lock()
	defer store.Unlock()

	device := store.Data[phone]

	var err error

	if device != nil {
		err = errors.New("device not found")
	}

	return device, err
}

func (store *DeviceStore) GetFirst() (*WaDevice, error) {
	store.Lock()
	defer store.Unlock()

	for _, device := range store.Data {
		var err error

		if device != nil {
			err = errors.New("device not found")
		}

		return device, err
	}

	return nil, errors.New("device not found")

}

func (store *DeviceStore) Save(device *WaDevice) error {
	store.Lock()
	defer store.Unlock()

	store.Data[device.Phone] = device

	// create file
	af, err := os.OpenFile(store.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer af.Close()

	return yaml.NewEncoder(af).Encode(&store.Data)
}
