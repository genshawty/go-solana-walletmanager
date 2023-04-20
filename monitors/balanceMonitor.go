package monitors

import (
	"fmt"
	wm "sol/walletmanager"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
)

func MonitorMenu(wallets wm.Wallets, rpcClient *rpc.Client) {
	fmt.Printf("\n")
	fmt.Println("Use input area to pick the action: ")
	fmt.Printf("[1] Sol balance monitor\n[2] Token balance monitor\n[3] Cmv2 monitor (TODO)\n")

	var choice string
	fmt.Scanf("%s", &choice)
	if choice == "" {
		fmt.Scanf("%s", &choice)
	}

	choiceInt, err := strconv.Atoi(choice)

	if err != nil {
		fmt.Printf("Use number")
		MonitorMenu(wallets, rpcClient)
	}

	switch choiceInt {
	case 1:
		solBalanceMonitor(wallets, rpcClient)
	case 2:
		tokenBalanceMonitor(wallets, rpcClient)
	}

	defer MonitorMenu(wallets, rpcClient)
	defer fmt.Println("---------------------------------")
}

func solBalanceMonitor(wallets wm.Wallets, rpcClient *rpc.Client) {
	for {
		fmt.Printf("\n-------------------------------------\n\n")
		wm.ShowWallets(wallets, rpcClient)
		time.Sleep(5 * time.Second)
	}
}

func tokenBalanceMonitor(wallets wm.Wallets, rpcClient *rpc.Client) {
	for {
		fmt.Printf("\n-------------------------------------\n\n")
		wm.ShowTokenBalance(wallets, rpcClient)
		time.Sleep(5 * time.Second)
	}
}
