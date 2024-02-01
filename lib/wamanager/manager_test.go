package wamanager_test

import (
	"testing"

	"github.com/Vaziria/coin_cli/lib/commonlib"
	"github.com/Vaziria/coin_cli/lib/wamanager"
	"github.com/stretchr/testify/assert"
)

type DataYml struct {
	Price int
	CC    int
}

func TestBasicFunction(t *testing.T) {
	mngr := wamanager.Manager{}

	base := commonlib.MockBaseLocation()

	fname := base.Path("qr.jpg")
	t.Run("test generate qr", func(t *testing.T) {

		err := mngr.GenerateQr("https://google.com", fname)
		assert.Nil(t, err)
	})

	t.Run("open browser qr", func(t *testing.T) {
		err := mngr.OpenBrowser(fname)

		assert.Nil(t, err)
	})

}

func TestWaInitialising(t *testing.T) {

	base := commonlib.MockBaseLocation()
	store := wamanager.NewDeviceStore(base.Path("wa_device.json"))
	err := store.Load()

	assert.Nil(t, err)

	_, err = wamanager.NewManager(
		&wamanager.ManagerConfiguration{
			Debug:  true,
			MeowDB: base.Path("wa_database.db"),
		},
		store,
	)

	assert.Nil(t, err)

	// t.Run("test manager add device", func(t *testing.T) {
	// 	mngr.AddDevice(context.Background(), base.Path("qr_data.jpg"))
	// 	time.Sleep(time.Minute)
	// })

	t.Run("getting device", func(t *testing.T) {

		// dev, err := mngr.ConnectDevice(&wamanager.WaDevice{
		// 	Phone: "6285804152031",
		// }, nil)

		// assert.Nil(t, err)

		// datas, err := dev.Client.GetJoinedGroups()
		// assert.Nil(t, err)

		// for _, data := range datas {
		// 	log.Println(data.GroupName.Name)
		// 	if data.GroupName.Name == "testing wa" {

		// 		log.Println("sending message")

		// 		dd := DataYml{
		// 			Price: 100,
		// 			CC:    213,
		// 		}

		// 		dastr, _ := yaml.Marshal(dd)

		// 		message := "test wa dari bot [*prepare vish lauch*]\n\n" + string(dastr)
		// 		_, err = dev.Client.SendMessage(context.Background(), data.JID, &proto.Message{
		// 			Conversation: &message,
		// 		})
		// 		assert.Nil(t, err)

		// 	}
		// }

		// time.Sleep(time.Hour)
	})
}
