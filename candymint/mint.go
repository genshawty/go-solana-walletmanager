package candymint

import (
	"context"
	"fmt"
	"log"
	"os"
	cm "sol/candy_machine"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/gagliardetto/solana-go/text"
)

// func getCandyMachineCreator(candyMachineId solana.PublicKey) solana.PublicKey {
// 	// PublicKey.findProgramAddress(
// 	// 	[Buffer.from('candy_machine'), candyMachine.toBuffer()],
// 	// 	CANDY_MACHINE_PROGRAM_ID,
// 	//   );
// 	candyMachineProgramId := solana.MustPublicKeyFromBase58("cndy3Z4yapfJBmL3ShUp5exZKqR3z33thTzeNMm2gRZ")

// 	candyMachineCreator, _, _ := solana.FindProgramAddress(
// 		[][]byte{
// 			[]byte("candy_machine"),
// 			candyMachineId.Bytes(),
// 		},
// 		candyMachineProgramId,
// 	)
// 	return candyMachineCreator
// }

// ищет MasterEdition, необходимый параметр для инструкции MintNft
func findMasterEdition(mint solana.PublicKey) solana.PublicKey {
	// Pubkey::find_program_address(
	//     &[
	//         PREFIX.as_bytes(),
	//         crate::id().as_ref(),
	//         mint.as_ref(),
	//         EDITION.as_bytes(),
	//     ],
	//     &crate::id(),
	// )
	const prefix string = "metadata"
	const edition string = "edition"

	programId := solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
	masterEdition, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte(prefix),
			programId.Bytes(),
			mint.Bytes(),
			[]byte(edition),
		},
		programId,
	)

	return masterEdition
}

// основная функция минта, в начало нужны рпц и ws клиент, айди кэндимашины и кошелек с которого минтишь я сейчас ввожу вручную, потому что для тестирования так проще
// тестирование происходит на моей коллекции на девнете, ее cmId 61CCEcc5LZCp9RJ2yPhmABHno3fPLKK6Lbxrpvqdj5Bu
func MetaplexMint(rpcClient *rpc.Client, wsClient *ws.Client) error {
	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)

	// константа, айди программы для взаимодействия с кэндимашиной Metaplex
	candyMachineProgramId, _ := solana.PublicKeyFromBase58("cndy3Z4yapfJBmL3ShUp5exZKqR3z33thTzeNMm2gRZ")

	infoLog.Printf("CandyMachineProgramId %s", candyMachineProgramId.String())

	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	// айди кэндимашины, кастомное для каждого минта и спокойно достается с сайта минта
	candyMachineId := solana.MustPublicKeyFromBase58("61CCEcc5LZCp9RJ2yPhmABHno3fPLKK6Lbxrpvqdj5Bu")

	infoLog.Printf("CandyMachineId %s", candyMachineId.String())

	// кошелек, с которого идет минт FIXME: необходимо передавать его в саму функцию, это будет прикручено как функция начнет впринципе работать
	payer := solana.MustPrivateKeyFromBase58("you wallet private key")

	infoLog.Printf("Payer public %s", payer.PublicKey().String())
	infoLog.Printf("Payer public %s", payer.String())

	// mint address, необходимо его сгенерить самому вручную
	mint := solana.NewWallet().PrivateKey

	infoLog.Printf("Mint public %s", mint.PublicKey().String())
	infoLog.Printf("Mint Private %s", mint.String())
	// применение: keys.go по поиску TokenProgramId 425 строка
	// нашел там чуть выше FindAssociatedTokenAddress, эквивалент тому что уже нашел

	// токен аккаунт для нашего минта, сгенерированного выше, возможно тут косячу, поскольку мы в блокчейн еще не передавали никаких данных и ничего не вызывали, но кажется все работает и так и надо
	recieverTokenAccount, _, err := solana.FindAssociatedTokenAddress(payer.PublicKey(), mint.PublicKey())

	infoLog.Printf("recieverTokenAdress public %s", recieverTokenAccount.String())

	if err != nil {
		fmt.Println("Не найден токен аккаунт")
		return err
	}

	lamports, _ := rpcClient.GetMinimumBalanceForRentExemption(
		context.TODO(),
		token.MINT_SIZE,
		rpc.CommitmentFinalized,
	)

	instructions := make([]solana.Instruction, 0, 10)

	// инструкция Create Account
	createAccountInstruction := system.NewCreateAccountInstruction(
		lamports,
		token.MINT_SIZE,
		solana.TokenProgramID,
		payer.PublicKey(),
		mint.PublicKey(),
	).Build()

	// fmt.Printf("create account inst \n")
	// writeAccountData(createAccountInstruction.Accounts())

	// инструкция Initialize Mint
	initializeMintInstruction := token.NewInitializeMintInstruction(
		0,
		payer.PublicKey(),
		payer.PublicKey(),
		mint.PublicKey(),
		solana.SysVarRentPubkey,
	).Build()

	fmt.Printf("initializeMintInstruction  %v\n", initializeMintInstruction.Accounts())

	// Инструкция создания минта для аккаунта из Create Account Instruction
	slpCreateInst := associatedtokenaccount.NewCreateInstruction(
		payer.PublicKey(),
		payer.PublicKey(),
		mint.PublicKey(),
	).Build()

	// fmt.Printf("slpCreateInst\n")
	// writeAccountData(slpCreateInst.Accounts())

	// Mint To instruction
	mintToInst := token.NewMintToInstruction(
		1,
		mint.PublicKey(),
		recieverTokenAccount,
		payer.PublicKey(),
		[]solana.PublicKey{},
	).Build()

	// fmt.Printf("Mint To Inst\n")
	// writeAccountData(mintToInst.Accounts())

	// создание массива (среза, терминология Go) инструкций
	instructions = append(instructions,
		createAccountInstruction,
		initializeMintInstruction,
		slpCreateInst,
		mintToInst,
	)

	// finding data for MintNft instruction

	// ищет аккаунт создателя кэндимашины, а так же bumpseed
	candyMachineCreator, creatorBump, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("candy_machine"),
			candyMachineId.Bytes(),
		}, candyMachineProgramId,
	)

	infoLog.Printf("candyMachineCreator public %s", candyMachineCreator.String())

	if err != nil {
		fmt.Println(err)
	}

	// поиск адреса метаданных для нашего минта, созданного ранее
	metadataAddress, _, _ := solana.FindTokenMetadataAddress(mint.PublicKey())

	// поиск мастер эдишн для нашего минта, созданного ранее
	masterEdition := findMasterEdition(mint.PublicKey())

	// checkMetadata, _, _ := solana.FindTokenMetadataAddress(solana.MustPublicKeyFromBase58("BPvoiNJLc5uVBBHpad7GMaX3hmXuP1SUTauVN1vgD1tX"))

	infoLog.Printf("metadataAddress public %s", metadataAddress.String())

	infoLog.Printf("masterEdition public %s", masterEdition.String())

	// infoLog.Printf("IS IT RIGHT %s", findMasterEdition(solana.MustPublicKeyFromBase58("BPvoiNJLc5uVBBHpad7GMaX3hmXuP1SUTauVN1vgD1tX")))

	// infoLog.Printf("IS IT RIGHT %s", checkMetadata.String())

	// ставит айди программы чтобы правильно создавалась инструкция MintNft
	cm.SetProgramID(candyMachineProgramId)

	// создание главной инструкции - MintNft, там где я создаю паблик кей внутри (solana.MustPublicKeyFromBase58("7eoXQLSQuKSaKuwK69BACmis6eyoTDqgmcXKtpCkNt41")), значит я это еще не нашел
	mintNftInst := cm.NewMintNftInstruction(
		creatorBump,
		candyMachineId,
		candyMachineCreator,
		payer.PublicKey(),
		solana.MustPublicKeyFromBase58("7eoXQLSQuKSaKuwK69BACmis6eyoTDqgmcXKtpCkNt41"),
		metadataAddress,
		mint.PublicKey(),
		payer.PublicKey(),
		payer.PublicKey(),
		masterEdition,
		solana.TokenMetadataProgramID,
		solana.TokenProgramID,
		solana.SystemProgramID,
		solana.SysVarRentPubkey,
		solana.SysVarClockPubkey,
		solana.SysVarRecentBlockHashesPubkey,
		solana.SysVarInstructionsPubkey,
	).Build()

	// это значение константно для кэндимашины
	collectionPda := solana.MustPublicKeyFromBase58("xbqsLYnVWWsziVyYSnjWdgnuaqHzJQ8nWAKEtWz3ARM")

	// это значение константно для кэндимашины
	collectionMint := solana.MustPublicKeyFromBase58("Dw4FNFm8UJVHUnt5cbFrg7LNSUZx6FyNbyAhHNRSJy2K")

	// это значение константно для кэндимашины
	collectionAuthorityRecord := solana.MustPublicKeyFromBase58("GwifH3VhXyefh2quEYsCPhWNpYHMVj7jL4jxXmkMeXMY")
	// TODO выдает какую то кастомную ошибку сейчас, запуск через go run main.go, затем в консоль цифру 2, затем 9 и читать консоль, в идеале в VS code это делать, тут вывод красивый
	setCollection := cm.NewSetCollectionDuringMintInstruction(
		candyMachineId,
		metadataAddress,
		payer.PublicKey(),
		collectionPda,
		solana.TokenMetadataProgramID,
		solana.SysVarInstructionsPubkey,
		collectionMint,
		metadataAddress,
		masterEdition,
		solana.MustPublicKeyFromBase58("7eoXQLSQuKSaKuwK69BACmis6eyoTDqgmcXKtpCkNt41"),
		collectionAuthorityRecord,
	).Build()

	fmt.Printf("mintNftInst")

	infoLog.Printf("program id %s", mintNftInst.ProgramID().String())

	instructions = append(instructions,
		mintNftInst,
		setCollection,
	)

	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(payer.PublicKey()),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("-----------------\n\n")
	fmt.Println(tx.String())
	fmt.Printf("\n\n-----------------")
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if payer.PublicKey().Equals(key) {
				return &payer
			}
			if mint.PublicKey().Equals(key) {
				return &mint
			}
			return nil
		},
	)
	if err != nil {
		panic(fmt.Errorf("unable to sign transaction: %w", err))
	}
	spew.Dump(tx)
	// Pretty print the transaction:
	tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Transfer SOL"))

	// Send transaction, and wait for confirmation:
	sig, err := confirm.SendAndConfirmTransaction(
		context.TODO(),
		rpcClient,
		wsClient,
		tx,
	)
	if err != nil {
		panic(err)
	}
	spew.Dump(sig)

	return nil
}

// побочная функция для отладки, выписывает список всех аккаунтов для инструкции, вроде как больше в ней не нуждаюсь, надо будет снести (уже снес)
// func writeAccountData(accs []*solana.AccountMeta) {
// 	for i := 0; i < len(accs); i++ {
// 		acc := accs[i]
// 		fmt.Printf("Account: %v, isSigner: %v, isWriter: %v\n", acc.PublicKey.String(), acc.IsSigner, acc.IsWritable)
// 	}
// }
