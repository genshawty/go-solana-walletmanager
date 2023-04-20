package walletmanager

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func PrintConstants() {
	candyMachineId := solana.MustPublicKeyFromBase58("2NfTUdcY7Yz9ja33LSj8X9jC8DGmmwUAwUxbkmcoN87Y")
	candyMachineProgramId, _ := solana.PublicKeyFromBase58("cndy3Z4yapfJBmL3ShUp5exZKqR3z33thTzeNMm2gRZ")
	candyMachineCreator, creatorBump, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("candy_machine"),
			candyMachineId.Bytes()[:],
		}, candyMachineProgramId,
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(candyMachineCreator, creatorBump)
}

func GetToken() {
	payer, _ := solana.PrivateKeyFromBase58("55b1tzCoj3Npbx5jPVAj8i8M33R3DnrgQRurac39vUE54YVZHRyA9hR8XfwDr2Dt37XAsMaMhppGBU5gkN6XVb9r")

	mint := solana.MustPublicKeyFromBase58("3BLdpFW9ocZaoTruHSi46DmY9TTjErNd4MfqUcE5HoCq")

	//candyMachineId := solana.MustPublicKeyFromBase58("2NfTUdcY7Yz9ja33LSj8X9jC8DGmmwUAwUxbkmcoN87Y")

	// tokenAddress := solana.FindProgramAddress([][]byte{
	// 	payer.PublicKey()[:],
	// 	candyMachineId[:],
	// 	mint[:],
	// },)
	recieverTokenAdress, _, err := solana.FindAssociatedTokenAddress(payer.PublicKey(), mint)

	fmt.Println(err)
	fmt.Println(recieverTokenAdress.String())
}

func Test() {
	endpoint := rpc.DevNet_RPC
	client := rpc.New(endpoint)

	txSig := solana.MustSignatureFromBase58("xJPXeC8bBUoQi9oSaFoaAm8DVhY6u4R86jfySqGDz7MHxe4ugY8N2pUrWa8L4j9itA3Pei9oYscGLryTR1LyeT7")
	{
		out, err := client.GetTransaction(
			context.TODO(),
			txSig,
			&rpc.GetTransactionOpts{
				Encoding: solana.EncodingBase64,
			},
		)
		if err != nil {
			panic(err)
		}
		spew.Dump(out)
		spew.Dump(out.Transaction.GetBinary())

		decodedTx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(out.Transaction.GetBinary()))
		if err != nil {
			panic(err)
		}
		spew.Dump(decodedTx)
	}
	{
		out, err := client.GetTransaction(
			context.TODO(),
			txSig,
			nil,
		)
		if err != nil {
			panic(err)
		}
		spew.Dump(out)
		spew.Dump(out.Transaction.GetParsedTransaction())
	}
}
