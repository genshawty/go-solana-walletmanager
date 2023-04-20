package walletmanager

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

var diff int = 10e5

func SplitMoney(wallets Wallets, rpcClient *rpc.Client, wsClient *ws.Client) error {
	fmt.Println("Choose wallet to drain money from")
	ShowWallets(wallets, rpcClient)
	var choice string
	fmt.Scanf("\n%s", &choice)
	index, err := strconv.Atoi(choice)
	walletFrom := wallets[index-1]

	if err != nil {
		return err
	}

	fmt.Println("Choose wallet to get money")
	fmt.Scanf("\n%s", &choice)

	fmt.Println("Enter amount of money should left on account")

	var amount string
	fmt.Scanf("\n%s", &amount)

	amountFloat, err := strconv.ParseFloat(amount, 64)

	if err != nil {
		return err
	}

	amountInt := int(amountFloat * 10e8)

	var wg sync.WaitGroup

	if strings.Contains(choice, "-") {
		index1, index2 := strings.Split(choice, "-")[0], strings.Split(choice, "-")[1]

		index1Int, err := strconv.Atoi(index1)

		if err != nil {
			return err
		}

		index2Int, err := strconv.Atoi(index2)

		if err != nil {
			return err
		}

		walletsTo := wallets[index1Int-1 : index2Int]

		wg.Add(len(walletsTo))
		errors := make([]error, len(walletsTo))

		for i, walletTo := range walletsTo {
			go CalculateAndTransfer(i, errors, index, index1Int+i, walletFrom, walletTo, amountInt, rpcClient, wsClient, &wg)
		}
		wg.Wait()

		for i, v := range errors {
			if v != nil {
				fmt.Printf("Error while sending money [%v] [%v] -> [%v] [%v], error: %v\n", index, walletFrom.Name, index1Int+i, walletsTo[i].Name, v)
			}
		}

		if err != nil {
			return err
		}

	} else {
		index1Int, err := strconv.Atoi(choice)

		if err != nil {
			return err
		}
		wg.Add(1)
		walletTo := wallets[index1Int-1]
		errors := make([]error, 1)
		go CalculateAndTransfer(0, errors, index1Int, index, walletFrom, walletTo, amountInt, rpcClient, wsClient, &wg)
		wg.Wait()
		for i, v := range errors {
			if v != nil {
				fmt.Printf("Error while sending money [%v] [%v] -> [%v] [%v], error: %v\n", index, walletFrom.Name, index1Int+i, walletTo.Name, v)
			}
		}
	}
	return nil
}

func CalculateAndTransfer(errorIndex int, errors []error, indexFrom int, indexTo int, accountFrom Wallet, accountTo Wallet, amount int, rpcClient *rpc.Client, wsClient *ws.Client, wg *sync.WaitGroup) error {
	defer wg.Done()
	balance, err := rpcClient.GetBalance(context.TODO(), accountTo.Key.PublicKey(), rpc.CommitmentFinalized)
	if err != nil {
		errors[errorIndex] = err
		return err
	}
	balanceValue := int(balance.Value)

	if (amount - balanceValue) > diff {
		TransferMoney(errorIndex, errors, indexTo, indexFrom, accountFrom, accountTo, amount-balanceValue, rpcClient, wsClient, wg)
	}
	return nil
}

func CollectMoney(wallets Wallets, rpcClient *rpc.Client, wsClient *ws.Client) error {
	fmt.Println("Choose wallet to get money")
	ShowWallets(wallets, rpcClient)
	var choice string

	fmt.Scanf("\n%s", &choice)
	index, err := strconv.Atoi(choice)
	walletTo := wallets[index-1]

	if err != nil {
		return err
	}

	fmt.Println("Choose wallets to collect money from")
	fmt.Scanf("\n%s", &choice)
	var wg sync.WaitGroup

	if strings.Contains(choice, "-") {
		index1, index2 := strings.Split(choice, "-")[0], strings.Split(choice, "-")[1]

		index1Int, err := strconv.Atoi(index1)

		if err != nil {
			return err
		}

		index2Int, err := strconv.Atoi(index2)

		if err != nil {
			return err
		}

		walletsFrom := wallets[index1Int-1 : index2Int]
		wg.Add(len(walletsFrom))
		errors := make([]error, len(walletsFrom))
		for i := 0; i < len(walletsFrom); i++ {
			go TransferMaxMoney(i, errors, index1Int+i, index, walletsFrom[i], walletTo, rpcClient, wsClient, &wg)
		}

		wg.Wait()
		for i, v := range errors {
			if v != nil {
				fmt.Printf("Error while sending money [%v] [%v] -> [%v] [%v], error: %v\n", index, walletsFrom[i].Name, index1Int+i, walletTo.Name, v)
			}
		}

	} else {
		index1Int, err := strconv.Atoi(choice)

		if err != nil {
			return err
		}
		walletFrom := wallets[index1Int-1]
		wg.Add(1)
		errors := make([]error, 1)
		go TransferMaxMoney(0, errors, index1Int, index, walletFrom, walletTo, rpcClient, wsClient, &wg)
		wg.Wait()
		for i, v := range errors {
			if v != nil {
				fmt.Printf("Error while sending money [%v] [%v] -> [%v] [%v], error: %v\n", index, walletFrom.Name, index1Int+i, walletTo.Name, v)
			}
		}
	}
	return nil
}

func TransferMoney(errorIndex int, errors []error, indexTo int, indexFrom int, accountFromWithName Wallet, accountToWithName Wallet, amount int, rpcClient *rpc.Client, wsClient *ws.Client, wg *sync.WaitGroup) error {
	accountFrom := accountFromWithName.Key
	accountTo := accountToWithName.Key.PublicKey()

	balance, err := rpcClient.GetBalance(context.TODO(), accountFrom.PublicKey(), rpc.CommitmentFinalized)

	if err != nil {
		errors[errorIndex] = err
		return err
	}
	if amount > int(balance.Value) {
		// fmt.Printf("balance: %v, amount: %v", balance.Value, amount)
		// fmt.Printf("Not enough balance [%v] -> [%v]\n", indexFrom, indexTo)
		errors[errorIndex] = fmt.Errorf("not enough money")
		return fmt.Errorf("not enough money")
	}

	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		errors[errorIndex] = err
		return err
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				uint64(amount),
				accountFrom.PublicKey(),
				accountTo,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(accountFrom.PublicKey()),
	)
	if err != nil {
		errors[errorIndex] = err
		return err
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if accountFrom.PublicKey().Equals(key) {
				return &accountFrom
			}
			return nil
		},
	)
	if err != nil {
		errors[errorIndex] = err
		return err
	}
	// spew.Dump(tx)
	// Pretty print the transaction:
	// tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Transfer SOL"))

	// Send transaction, and wait for confirmation:
	// sig, err := confirm.SendAndConfirmTransaction(
	// 	context.TODO(),
	// 	rpcClient,
	// 	wsClient,
	// 	tx,
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// spew.Dump(sig)

	// Or just send the transaction WITHOUT waiting for confirmation:
	sig, err := rpcClient.SendTransactionWithOpts(
		context.TODO(),
		tx,
		false,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		errors[errorIndex] = err
		return err
	}
	status, err := ValidateTransaction(rpcClient, sig)

	if err != nil {
		fmt.Printf("Transfering money from: [%v] [%v] to [%v] [%v], ERROR: %v\n", indexFrom, accountFromWithName.Name, indexTo, accountToWithName.Name, err)
		time.Sleep(30 * time.Second)
		status, err = ValidateTransaction(rpcClient, sig)

		if err != nil {
			TransferMoney(errorIndex, errors, indexTo, indexFrom, accountFromWithName, accountToWithName, amount, rpcClient, wsClient, wg)
		}
	}
	fmt.Printf("Transfering money from: [%v] [%v] to [%v] [%v], sig: %v, STATUS: %v\n", indexFrom, accountFromWithName.Name, indexTo, accountToWithName.Name, sig.String(), status)
	return nil
	// checkTransferStatus(indexFrom, indexTo, sig, wsClient)
	// spew.Dump(sig)
}

func TransferMaxMoney(errorIndex int, errors []error, indexFrom int, indexTo int, accountFromWithName Wallet, accountToWithName Wallet, rpcClient *rpc.Client, wsClient *ws.Client, wg *sync.WaitGroup) error {
	defer wg.Done()

	accountFrom := accountFromWithName.Key
	accountTo := accountToWithName.Key.PublicKey()

	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		errors[errorIndex] = err
		return err
	}

	amount, err := rpcClient.GetBalance(context.TODO(), accountFrom.PublicKey(), rpc.CommitmentFinalized)
	if err != nil {
		errors[errorIndex] = err
		return err
	}

	if amount.Value < uint64(diff) {
		errors[errorIndex] = fmt.Errorf("already drained wallet")
		return nil
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amount.Value,
				accountFrom.PublicKey(),
				accountTo,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(accountFrom.PublicKey()),
	)
	if err != nil {
		errors[errorIndex] = err
		return err
	}

	out, err := rpcClient.GetFeeForMessage(context.TODO(), tx.Message.ToBase64(), rpc.CommitmentFinalized)

	if err != nil {
		errors[errorIndex] = err
		return err
	}

	finalAmount := amount.Value - *out.Value

	tx, err = solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				finalAmount,
				accountFrom.PublicKey(),
				accountTo,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(accountFrom.PublicKey()),
	)
	if err != nil {
		errors[errorIndex] = err
		return err
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if accountFrom.PublicKey().Equals(key) {
				return &accountFrom
			}
			return nil
		},
	)
	if err != nil {
		errors[errorIndex] = err
		return err
	}
	// spew.Dump(tx)
	// Pretty print the transaction:
	// tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Transfer SOL"))

	// Send transaction, and wait for confirmation:
	sig, err := rpcClient.SendTransactionWithOpts(
		context.TODO(),
		tx,
		false,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		errors[errorIndex] = err
		return err
	}
	status, err := ValidateTransaction(rpcClient, sig)

	if err != nil {
		fmt.Printf("Transfering money from: [%v] [%v] to [%v] [%v], ERROR: %v\n", indexFrom, accountFromWithName.Name, indexTo, accountToWithName.Name, err)
		time.Sleep(30 * time.Second)
		status, err = ValidateTransaction(rpcClient, sig)

		if err != nil {
			TransferMaxMoney(errorIndex, errors, indexFrom, indexTo, accountFromWithName, accountToWithName, rpcClient, wsClient, wg)
		}
	}
	fmt.Printf("Transfering money from: [%v] [%v] to [%v] [%v], sig: %v, STATUS: %v\n", indexFrom, accountFromWithName.Name, indexTo, accountToWithName.Name, sig.String(), status)
	return nil
}

// func checkTransferStatus(indexFrom int, indexTo int, sig solana.Signature, wsClient *ws.Client) {
// 	sub, err := wsClient.SignatureSubscribe(
// 		sig,
// 		rpc.CommitmentFinalized,
// 	)
// 	if err != nil {
// 		fmt.Printf("[%v] -> [%v], sig: %v, err: %v\n", indexFrom, indexTo, sig, err)
// 	}
// 	defer sub.Unsubscribe()
// 	defer fmt.Println("unsubscribe")
// 	for {
// 		got, _ := sub.Recv()
// 		spew.Dump(got)
// 		time.Sleep(8 * time.Second)
// 	}
// }
