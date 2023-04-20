package main

import (
	"context"
	"fmt"
	"os"
	"sol/monitors"
	wm "sol/walletmanager"
	"strconv"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

// const rpcAdress string = "https://rpc.ankr.com/solana"
// const lamport float64 = 10e8

func main() {
	// rpcClient, wsClient, name := pickNetwork()
	rpcClient := getRpcClient()
	var err error
	wsClient, err := ws.Connect(context.Background(), rpc.MainNetBeta_WS)

	name := "data/wallets.csv"

	if err != nil {
		panic(err)
	}

	wallets := wm.ReadWallets(name)

	menu(wallets, rpcClient, wsClient, name)
}

// func testWriteWallets() {
// 	wallets := make([][]byte, 0, 20)
// 	var walletS string

// 	fmt.Scanf("%s", &walletS)

// 	wallet, _ := base58.Decode(walletS)

// 	wallets = append(wallets, wallet)

// 	wm.WriteWallets(wallets)
// }

func getRpcClient() *rpc.Client {
	file, err := os.ReadFile("data/rpc.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(file))
	return rpc.New(string(file))
}

func menu(wallets wm.Wallets, rpcClient *rpc.Client, wsClient *ws.Client, name string) {
	fmt.Printf("\n")
	fmt.Println("Use input area to pick the action: ")
	fmt.Println("[0] CMV2 MINT (TODO)\n[1] Create Wallet\n[2] Split money\n[3] Collect money\n[4] Show wallets\n[5] Get token Balance\n[6] Collect tokens\n[7] Export for Urban\n[8] Monitors list")

	var choice string
	fmt.Scanf("%s", &choice)
	if choice == "" {
		fmt.Scanf("%s", &choice)
	}

	choiceInt, err := strconv.Atoi(choice)

	if err != nil {
		fmt.Printf("Use number")
		menu(wallets, rpcClient, wsClient, name)
	}

	switch choiceInt {
	// case 0:
	// 	cmv2.MetaplexMint(rpcClient, wsClient)
	// case 1:
	// 	wallet, err := wm.ImportWallet(wallets, name)
	// 	wallets = append(wallets, wallet)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	case 1:
		wallets = wm.CreateWallet(wallets, name)

		fmt.Println("Actual wallets")

		wm.ShowWallets(wallets, rpcClient)

	case 2:
		err = wm.SplitMoney(wallets, rpcClient, wsClient)

		if err != nil {
			fmt.Println(err)
		}
	case 3:
		wm.CollectMoney(wallets, rpcClient, wsClient)

	case 4:
		err = wm.ShowWallets(wallets, rpcClient)

		if err != nil {
			fmt.Println(err)
		}
	case 5:
		wm.ShowTokenBalance(wallets, rpcClient)
	case 6:
		wm.CollectTokensMenu(wallets, rpcClient, wsClient)
	case 7:
		wm.ExportForUrbanMenu(wallets)
	case 8:
		monitors.MonitorMenu(wallets, rpcClient)
	}
	defer menu(wallets, rpcClient, wsClient, name)
	defer fmt.Println("---------------------------------")
}

func pickNetwork() (*rpc.Client, *ws.Client, string) {
	fmt.Println("Pick network: 1 - MainNetwork, 2 - DevNet")
	var choice string
	fmt.Scanf("%s", &choice)

	choiceInt, _ := strconv.Atoi(choice)

	var rpcClient *rpc.Client
	var wsClient *ws.Client
	var name string

	if choiceInt == 1 {
		rpcClient = rpc.New("https://purple-silent-snowflake.solana-mainnet.quiknode.pro/f377280d0b9d0475ede1ee8c88ac5f45c5165c57/")
		var err error
		wsClient, err = ws.Connect(context.Background(), rpc.MainNetBeta_WS)

		name = "data/wallets.csv"

		if err != nil {
			panic(err)
		}
	} else {
		rpcClient = rpc.New(rpc.DevNet_RPC)
		var err error
		wsClient, err = ws.Connect(context.Background(), rpc.DevNet_WS)

		name = "data/devwallets.csv"

		if err != nil {
			panic(err)
		}
	}
	return rpcClient, wsClient, name
}
