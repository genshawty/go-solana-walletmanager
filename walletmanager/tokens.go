package walletmanager

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func getTokenCount(pubKey solana.PublicKey, rpcClient *rpc.Client, i int, balances []int, errors []error, wg *sync.WaitGroup) error {
	defer wg.Done()
	count := 0

	var wgMini sync.WaitGroup
	var out *rpc.GetTokenAccountsResult
	var err error

	retries := 0
	for (err != nil) || (retries < 3) {
		// fmt.Printf("get token count: %v\n", (err != nil) || (retries <= 5))
		out, err = rpcClient.GetTokenAccountsByOwner(
			context.TODO(),
			pubKey,
			&rpc.GetTokenAccountsConfig{
				ProgramId: solana.TokenProgramID.ToPointer(),
			},
			&rpc.GetTokenAccountsOpts{
				Encoding: solana.EncodingBase64Zstd,
			},
		)
		retries++
	}
	if err != nil {
		errors[i] = err
		return err
	}

	if len(out.Value) == 0 {
		balances[i] = 0

	} else {
		wgMini.Add(len(out.Value))
		for _, rawAccount := range out.Value {
			go getTokenBalance(rawAccount.Pubkey.ToPointer(), rpcClient, &count, errors, i, &wgMini)
		}
		wgMini.Wait()
		if errors[i] != nil {
			balances[i] = -1
		} else {
			balances[i] = count
		}
	}
	return nil
}

// i - номер кошелька
func getTokenBalance(pubKey *solana.PublicKey, rpcClient *rpc.Client, count *int, errors []error, i int, wg *sync.WaitGroup) error {
	defer wg.Done()
	var err error
	var amount *rpc.GetTokenAccountBalanceResult
	retries := 0
	for (err != nil) || (retries < 3) {
		// fmt.Printf("get token balance: %v\n", (err != nil) || (retries <= 5))
		amount, err = rpcClient.GetTokenAccountBalance(
			context.TODO(),
			*pubKey,
			rpc.CommitmentFinalized,
		)

		retries++
	}

	if err != nil {
		errors[i] = err
		return err
	}

	if err == nil {
		amountInt, err := strconv.Atoi(amount.Value.UiAmountString)
		if err == nil {
			*count = *count + amountInt
			return nil
		} else {
			return nil
		}
	} else {
		errors[i] = err
		return err
	}
}

func ShowTokenBalance(walletsWithNames Wallets, rpcClient *rpc.Client) {
	wallets := make([]solana.PrivateKey, len(walletsWithNames))
	names := make([]string, len(walletsWithNames))

	for i := range wallets {
		wallets[i] = walletsWithNames[i].Key
		names[i] = walletsWithNames[i].Name
	}

	start := time.Now()
	balances := make([]int, len(wallets))

	var wg sync.WaitGroup
	wg.Add(len(wallets))
	errors := make([]error, len(wallets))

	for i, wallet := range wallets {
		go getTokenCount(wallet.PublicKey(), rpcClient, i, balances, errors, &wg)
	}
	wg.Wait()

	for i, wallet := range wallets {
		if errors[i] != nil {
			fmt.Printf("[%v] [%v], public: %v, ERROR: %v\n", names[i], i+1, wallet.PublicKey().String(), errors[i])
		} else {
			fmt.Printf("[%v] [%v] public: %v, balance: %v\n", names[i], i+1, wallet.PublicKey().String(), balances[i])
		}
	}
	fmt.Printf("taken time: %v", time.Since(start))
}

func CollectTokensMenu(wallets Wallets, rpcClient *rpc.Client, wsClient *ws.Client) {
	fmt.Println("Choose wallet to recieve tokens")
	// ShowTokenBalance(wallets, rpcClient)

	var choice string
	fmt.Scanf("\n%s", &choice)
	index, _ := strconv.Atoi(choice)

	walletTo := wallets[index-1]

	walletsFrom := make(Wallets, 1)

	fmt.Println("Choose wallets to send tokens from (\"all\" for all)")
	fmt.Scanf("\n%s", &choice)
	if choice == "all" {
		walletsFrom = append(wallets[:index-1], wallets[index:]...)
	} else if strings.Contains(choice, "-") {
		index1, index2 := strings.Split(choice, "-")[0], strings.Split(choice, "-")[1]

		index1int, _ := strconv.Atoi(index1)
		index2int, _ := strconv.Atoi(index2)

		walletsFrom = wallets[index1int-1 : index2int]
	} else {
		indexFrom, _ := strconv.Atoi(choice)
		walletsFrom[0] = wallets[indexFrom-1]
	}

	var wg sync.WaitGroup
	wg.Add(len(walletsFrom))
	for _, walletFrom := range walletsFrom {
		go TransferTokens(walletTo, walletFrom, rpcClient, &wg)
	}
	wg.Wait()
}

// Returns list of token accounts by public key of solana wallet
func GetAllTokens(wallet solana.PublicKey, rpcClient *rpc.Client) ([]solana.PublicKey, error) {
	var out *rpc.GetTokenAccountsResult
	var err error

	out, err = rpcClient.GetTokenAccountsByOwner(
		context.TODO(),
		wallet,
		&rpc.GetTokenAccountsConfig{
			ProgramId: solana.TokenProgramID.ToPointer(),
		},
		&rpc.GetTokenAccountsOpts{
			Encoding: solana.EncodingBase64Zstd,
		},
	)

	if err != nil {
		return []solana.PublicKey{}, err
	}

	tokenAccounts := make([]solana.PublicKey, 0, len(out.Value))

	var amount *rpc.GetTokenAccountBalanceResult

	for _, account := range out.Value {

		amount, err = rpcClient.GetTokenAccountBalance(
			context.TODO(),
			account.Pubkey,
			rpc.CommitmentFinalized,
		)
		if err != nil {
			return []solana.PublicKey{}, err
		}
		if amount.Value.UiAmountString != "0" {
			_, err1 := strconv.Atoi(amount.Value.UiAmountString)
			if err1 != nil {
				return []solana.PublicKey{}, err
			}
			tokenAccounts = append(tokenAccounts, account.Pubkey)
		}
	}
	return tokenAccounts, nil
}
