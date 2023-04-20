package walletmanager

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// func monitorTransaction(rpcClient *rpc.Client, sig solana.Signature) error {
// 	for {
// 		status, err := fetchTxStatus(rpcClient, sig)

// 		if err == nil {
// 			fmt.Printf("status: %v\n", status)

// 			if status == "finalized" {
// 				return nil
// 			}
// 		} else {
// 			fmt.Printf("error: %v\n", err)
// 		}
// 		time.Sleep(time.Second)
// 	}
// }

func ValidateTransaction(rpcClient *rpc.Client, sig solana.Signature) (string, error) {
	time.Sleep(time.Second)
	var status string
	var err error
	retries := 0

	for (err != nil) || (retries < 5) {
		status, err = fetchTxStatus(rpcClient, sig)
		time.Sleep(time.Second)
		retries++
	}
	return status, err
}

func fetchTxStatus(rpcClient *rpc.Client, sig solana.Signature) (string, error) {
	out, err := rpcClient.GetSignatureStatuses(context.TODO(), true, sig)
	if err != nil {
		return "", err
	}

	value := out.Value

	if len(value) == 0 {
		return "", fmt.Errorf("cant find the transaction")
	}
	if value[0] == nil {
		return "", fmt.Errorf("cant find the transaction")
	}
	if value[0].Err == nil {
		return string(value[0].ConfirmationStatus), nil
	} else {
		return "", err
	}
}
