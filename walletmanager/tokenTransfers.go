package walletmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

// Recieves token account address returns token mint address
func GetTokenMintAddress(pubKey solana.PublicKey, rpcClient *rpc.Client) (string, error) {
	var resp *rpc.GetAccountInfoResult
	var err error
	resp, err = rpcClient.GetAccountInfoWithOpts(
		context.TODO(),
		pubKey,
		// You can specify more options here:
		&rpc.GetAccountInfoOpts{
			Encoding:   solana.EncodingJSONParsed,
			Commitment: rpc.CommitmentFinalized,
			// You can get just a part of the account data by specify a DataSlice:
			// DataSlice: &rpc.DataSlice{
			//  Offset: pointer.ToUint64(0),
			//  Length: pointer.ToUint64(1024),
			// },
		},
	)

	if err != nil {
		return "", err
	}
	type mintResponce struct {
		Parsed struct {
			Info struct {
				IsNative    bool   `json:"isNative"`
				Mint        string `json:"mint"`
				Owner       string `json:"owner"`
				State       string `json:"state"`
				TokenAmount struct {
					Amount         string  `json:"amount"`
					Decimals       int     `json:"decimals"`
					UIAmount       float64 `json:"uiAmount"`
					UIAmountString string  `json:"uiAmountString"`
				} `json:"tokenAmount"`
			} `json:"info"`
			Type string `json:"type"`
		} `json:"parsed"`
		Program string `json:"program"`
		Space   int    `json:"space"`
	}

	var mapResp mintResponce

	jsonResp, err := resp.Value.Data.GetRawJSON().MarshalJSON()

	if err != nil {
		return "", err
	}

	json.Unmarshal(jsonResp, &mapResp)

	mint := mapResp.Parsed.Info.Mint

	return mint, nil
}

// writing to mintAccounts array
func getAndWriteMintAccount(tokenAccount solana.PublicKey, rpcClient *rpc.Client, i int, mintAccounts []string, errors []error, wg *sync.WaitGroup) error {
	defer wg.Done()
	var err error
	mintAccounts[i], err = GetTokenMintAddress(tokenAccount, rpcClient)

	if err != nil {
		errors[i] = err
		return err
	}
	return nil
}

// returns map of {tokenAcc: mintAddress}
func getTokenAndMintAccounts(wallet solana.PublicKey, rpcClient *rpc.Client) (map[string]string, error) {
	var err error
	tokensMap := make(map[string]string)

	tokenAccounts, err := GetAllTokens(wallet, rpcClient)

	if err != nil {

		return tokensMap, err
	}

	mintAccounts := make([]string, len(tokenAccounts))

	var wg sync.WaitGroup

	wg.Add(len(mintAccounts))

	errors := make([]error, len(tokenAccounts))
	for i, acc := range tokenAccounts {
		go getAndWriteMintAccount(acc, rpcClient, i, mintAccounts, errors, &wg)
	}
	wg.Wait()

	for _, val := range errors {
		if val != nil {
			return tokensMap, err
		}
	}

	if err != nil {
		return tokensMap, err
	}

	for i, value := range tokenAccounts {
		tokensMap[value.String()] = mintAccounts[i]
	}

	return tokensMap, nil
}

func TransferTokens(walletToStruct Wallet, walletFromStruct Wallet, rpcClient *rpc.Client, wg *sync.WaitGroup) error {
	defer wg.Done()

	walletTo := walletToStruct.Key
	toName := walletToStruct.Name

	walletFrom := walletFromStruct.Key
	fromName := walletFromStruct.Name

	var err error
	tokensMap, err := getTokenAndMintAccounts(walletFrom.PublicKey(), rpcClient)

	if err != nil {
		fmt.Printf("Token sent from: %v, to: %v, ERROR: %v\n", fromName, toName, err)
		return err
	}

	// for k, v := range tokensMap {
	// 	fmt.Printf("Token account: %v, mint address: %v\n", k, v)
	// }
	var wgMini sync.WaitGroup
	wgMini.Add(len(tokensMap))

	for k, v := range tokensMap {
		source := solana.MustPublicKeyFromBase58(k)
		mint := solana.MustPublicKeyFromBase58(v)

		go transferByMintAddressWithDefer(fromName, toName, walletFrom, source, mint, walletTo, rpcClient, &wgMini)
	}
	wgMini.Wait()
	return nil
}

func transferByMintAddressWithDefer(fromName string, toName string, walletFrom solana.PrivateKey, source solana.PublicKey, mint solana.PublicKey, walletTo solana.PrivateKey, rpcClient *rpc.Client, wg *sync.WaitGroup) error {
	defer wg.Done()
	_ = transferByMintAddress(fromName, toName, walletFrom, source, mint, walletTo, rpcClient)
	return nil
}

func transferByMintAddress(fromName string, toName string, walletFrom solana.PrivateKey, source solana.PublicKey, mint solana.PublicKey, walletTo solana.PrivateKey, rpcClient *rpc.Client) error {
	instructions := make([]solana.Instruction, 0, 2)

	destination, _, err := solana.FindAssociatedTokenAddress(
		walletTo.PublicKey(),
		mint,
	)
	if err != nil {
		fmt.Printf("Token sent from: [%v], to: [%v], ERROR while finding destination: %v\n", fromName, toName, err)
		return err
	}
	// fmt.Printf("wallet to public: %v", walletTo.PublicKey().String())
	// fmt.Printf("destination: %v", destination.String())

	_, err = GetTokenMintAddress(destination, rpcClient) // here was recieverAddress

	// if err != nil {
	// 	fmt.Printf("Token sent from: %v, to: %v, ERROR: %v\n", walletFrom.PublicKey().String(), walletTo.PublicKey().String(), err)
	// 	return err
	// }

	// fmt.Printf("reciever address: %v", recieverAddress)

	if err != nil {
		instructions = append(instructions,
			associatedtokenaccount.NewCreateInstruction(
				walletFrom.PublicKey(),
				walletTo.PublicKey(),
				mint,
			).Build(),
		)
	}

	instructions = append(instructions,
		token.NewTransferInstruction(
			1,
			source,
			destination,
			walletFrom.PublicKey(),
			[]solana.PublicKey{},
		).Build(),
	)

	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)

	if err != nil {
		fmt.Printf("Token sent from: [%v], to: [%v], ERROR: %v\n", fromName, toName, err)
		return err
	}

	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(walletFrom.PublicKey()),
	)
	if err != nil {
		fmt.Printf("Token sent from: [%v], to: [%v], ERROR: %v\n", fromName, toName, err)
		return err
	}
	// spew.Dump(tx)

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if walletFrom.PublicKey().Equals(key) {
				return &walletFrom
			}
			return nil
		},
	)
	if err != nil {
		fmt.Printf("Token sent from: [%v], to: [%v], ERROR: %v\n", fromName, toName, err)
		return err
	}
	sig, err := rpcClient.SendTransactionWithOpts(
		context.TODO(),
		tx,
		false,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		fmt.Printf("Token sent from: [%v], to: [%v], ERROR: %v\n", fromName, toName, err)
		return err
	}

	status, err := ValidateTransaction(rpcClient, sig)

	if err != nil {
		fmt.Printf("Token sent from: [%v], to: [%v], ERROR: %v\n", fromName, toName, err)
		time.Sleep(30 * time.Second)
		status, err = ValidateTransaction(rpcClient, sig)

		if err != nil {
			transferByMintAddress(fromName, toName, walletFrom, source, mint, walletTo, rpcClient)
		}
	}

	fmt.Printf("Token sent from: [%v], to: [%v], sig: %v, status: %v\n", fromName, toName, sig.String(), status)
	return nil
}
