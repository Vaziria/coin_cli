package walletcli_test

// func TestScenarioPakZen(t *testing.T) {
// 	cancel, err := walletcli.RunServiceDaemon("D:/testvish/vishaid.exe", "D:/testvish2")

// 	// time.Sleep(time.Hour)
// 	assert.Nil(t, err)

// 	defer cancel()

// 	cli := walletcli.WalletCli{
// 		Host:     "http://localhost:14277",
// 		Username: "virtuoso",
// 		Password: "virtuoso",
// 	}

// 	cli.WaitFullSync(time.Minute * 10)

// 	dis := walletcli.DistributeToWallets{
// 		Config: &walletcli.DistributeConfig{
// 			UnspentThreeshold: 200,
// 			SendPercent:       0.6,
// 			SendSleep:         0.1,
// 		},
// 		WalletCount: 500,
// 		Cli:         &cli,
// 	}

// 	// err = dis.Distribute(context.Background())
// 	err = dis.DistributeAllWallet(context.Background())
// 	assert.Nil(t, err)
// }
