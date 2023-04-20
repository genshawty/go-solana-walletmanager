package walletmanager

import (
	"context"
	"encoding/csv"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type Wallets []Wallet

type Wallet struct {
	Name string
	Key  solana.PrivateKey
}

func Clear() {
	fmt.Printf("\x1bc")
}

func ExportForUrbanMenu(wallets Wallets) {
	fmt.Println("Pick wallets to export")

	var choice string
	fmt.Scanf("\n%s", &choice)

	if strings.Contains(choice, "-") {
		index1, index2 := strings.Split(choice, "-")[0], strings.Split(choice, "-")[1]

		index1Int, err := strconv.Atoi(index1)

		if err != nil {
			panic(err)
		}

		index2Int, err := strconv.Atoi(index2)

		if err != nil {
			panic(err)
		}

		walletsExport := wallets[index1Int-1 : index2Int]

		exportForUrban(walletsExport)
	} else {
		index1Int, err := strconv.Atoi(choice)

		if err != nil {
			panic(err)
		}
		walletsExport := Wallets{wallets[index1Int-1]}
		exportForUrban(walletsExport)
	}
}

func exportForUrban(wallets Wallets) {
	exportWallets := make([][]string, len(wallets))

	for i, wallet := range wallets {
		exportWallets[i] = []string{wallet.Name, wallet.Key.String()}
	}

	filePath := "data/export.txt"

	var data []byte
	data, err := os.ReadFile(filePath)

	if err != nil {
		os.Create(filePath)
		data, _ = os.ReadFile(filePath)
	}
	// if string(data[len(data)-1]) != "\n" {
	// 	exportData = append([]string{"\n"}, exportData...)
	// }

	for _, v := range exportWallets {
		data = append(data, []byte(v[0]+":"+v[1]+"\n")...)
	}

	err = os.WriteFile(filePath, data, 0644)

	if err != nil {
		panic(err)
	}
}

// func DeleteWallet(wallets Wallets, name string, rpcClient *rpc.Client) Wallets {
// 	fmt.Println("Choose wallet for deleting")
// 	ShowWallets(wallets, rpcClient)
// 	var indexStr string

// 	fmt.Scanf("\n%s", &indexStr)

// 	index, err := strconv.Atoi(indexStr)

// 	if err != nil {
// 		panic(err)
// 	}

// 	newWallets := make(Wallets, 0, 20)
// 	for i := 0; i < len(wallets); i++ {
// 		if i != (index - 1) {
// 			newWallets = append(newWallets, wallets[i])
// 		}
// 	}
// 	WriteWallets(newWallets, name)
// 	return newWallets
// }

func ShowWallets(walletsWithNames Wallets, rpcClient *rpc.Client) error {
	wallets := make([]solana.PrivateKey, len(walletsWithNames))
	names := make([]string, len(walletsWithNames))

	for i := range wallets {
		wallets[i] = walletsWithNames[i].Key
		names[i] = walletsWithNames[i].Name
	}

	balances := make([]string, len(wallets))
	errors := make([]error, len(wallets))

	var wg sync.WaitGroup
	wg.Add(len(wallets))
	for i, key := range wallets {
		go writeBalance(key.PublicKey(), i, balances, errors, rpcClient, &wg)
	}
	wg.Wait()

	for i, key := range wallets {
		pub := key.PublicKey()
		if errors[i] != nil {
			fmt.Printf("[%v] [%v], public: %v, balance: -1, error: %v", names[i], i+1, pub, errors[i])
		} else {
			fmt.Printf("[%v] [%v], public: %v, balance: %v\n", names[i], i+1, pub, balances[i])
		}

	}
	return nil
}

func writeBalance(wallet solana.PublicKey, index int, balances []string, errors []error, rpcClient *rpc.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	balance, err := GetBalance(wallet, rpcClient)
	if err != nil {
		balances[index] = "-1"
		errors[index] = err
	} else {
		balances[index] = balance
	}

}

// func ImportWallet(wallets Wallets, name string) (solana.PrivateKey, error) {
// 	fmt.Println("Input wallet:")

// 	var priv string

// 	fmt.Scanf("\n%s", &priv)

// 	decoded, _ := solana.PrivateKeyFromBase58(priv)

// 	wallets = append(wallets, decoded)

// 	err := WriteWallets(wallets, name)

// 	if err != nil {
// 		return decoded, err
// 	} else {
// 		fmt.Printf("\nImported wallet: %v", priv)
// 		return decoded, nil
// 	}
// }

func CreateWallet(wallets Wallets, name string) Wallets {
	fmt.Println("Enter amout of wallets to create")
	var choice string

	fmt.Scanf("\n%s", &choice)
	index, _ := strconv.Atoi(choice)

	fmt.Println("Pick base name: ")
	var base string

	fmt.Scanf("\n%s", &base)

	for i := 0; i < index; i++ {
		priv := solana.NewWallet().PrivateKey
		walletName := base + fmt.Sprint(i+1)
		wallets = append(wallets, Wallet{walletName, priv})
	}
	WriteWallets(wallets, name)
	return wallets
}

func WriteWallets(wallets Wallets, name string) error {
	exportWallets := make([][]string, len(wallets))

	for i, wallet := range wallets {
		exportWallets[i] = []string{wallet.Name, wallet.Key.String()}
	}
	if err := os.Truncate(name, 0); err != nil {
		return err
	}

	csvFile, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return err
	}

	csvwriter := csv.NewWriter(csvFile)

	err = csvwriter.WriteAll(exportWallets)

	return err
}

func ReadWallets(name string) Wallets {
	csvFile, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	wallets := make(Wallets, 0, 30)
	for _, line := range csvLines {
		name := line[0]
		key := solana.MustPrivateKeyFromBase58(line[1])

		wallets = append(wallets, Wallet{name, key})
	}
	return wallets
}

func GetBalance(wallet solana.PublicKey, rpcClient *rpc.Client) (string, error) {
	out, err := rpcClient.GetBalance(
		context.TODO(),
		wallet,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return "-1", err
	}

	var lamportsOnAccount = new(big.Float).SetUint64(out.Value)
	// Convert lamports to sol:
	var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))

	// WARNING: this is not a precise conversion.
	return solBalance.Text('f', 10), nil
	// return string(out.Value), nil
}
