package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

/*Cmd Varaiable*/
const (
	CmdBlockchain          = "blockchain"
	CmdBlockchainCreate    = "-create"
	CmdBlockchainList      = "-list"
	CmdBlockchainListBlock = "block"
	CmdBlockchainAdd       = "-add"
	CmdTransaction         = "transaction"
	CmdTransactionQuery    = "-query"
	CmdTransactionSend     = "-send"
	CmdTransactionSendFrom = "-from"
	CmdTransactionSendTo   = "-to"
	CmdUTXOSet             = "utxoset"
	CmdUTXOSetReindex      = "-reindex"
	CmdUTXOSetQuery        = "-query"
)

func main() {

	parseFlags()
	return
}

func printUsage() {

	usage := `	
	USAGE : lubitcoin <module> [<action>] [action parameter]

	 <blockchain> [xxx]
		 blockchain -create : create blockchain with genesis block
		 blockchain -list  : list blockchain all blocks
		 blockchain -add xxxx : add a block

	 <transaction> [xx]
		 transaction -query xxx : query address balance
		 transaction -send 'amount_00' -from 'address_xx' -to 'address_yy' : transfer

	 <utxoset>
		 utxoset -reindex : rebuild utxoset
		 utxoset -query xxx : query address by utxoset
		
	`
	fmt.Println(usage)
}

func parseFlags() {

	// blockchainCmd := flag.NewFlagSet("blockchain", flag.ExitOnError)
	// transactionCmd := flag.NewFlagSet("transaction", flag.ExitOnError)
	// utxosetCmd := flag.NewFlagSet("utxoset", flag.ExitOnError)
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case CmdBlockchain:
		cmdBlockChain(os.Args[2:])
	case CmdTransaction:
		cmdTransaction(os.Args[2:])
	case CmdUTXOSet:
		cmdUTXOSet(os.Args[2:])
	}

	return

}

//////* Transaction

func cmdTransaction(args []string) {

	if len(args) < 2 {
		printUsage()
		os.Exit(0)
	}

	log.Println(args)

	switch args[0] {
	case CmdTransactionQuery:
		AddressQuery(args[1])
	case CmdTransactionSend:
		amount, _ := strconv.Atoi(args[1])
		AddressTransfer([]byte(args[3]), []byte(args[5]), amount)
	}
	log.Println(args)
}

/////  UTXO
func cmdUTXOSet(args []string) {
	if len(args) < 2 {
		printUsage()
		os.Exit(0)
	}
	switch args[0] {
	case CmdUTXOSetReindex:
		UTXOReindex()
	case CmdUTXOSetQuery:
		UTXOQuery([]byte(args[1]))
	}
	return
}

///// Blockchain
func cmdBlockChain(args []string) {

	if len(args) == 0 {
		printUsage()
		os.Exit(0)
	}
	switch args[0] {
	case CmdBlockchainCreate:
		BlockchainGenesis()
	case CmdBlockchainList:
		if len(args) < 2 {
			BlockchainListBlocks()
		} else if CmdBlockchainListBlock == args[1] {
			BlockchainListBlocks()
		} else {
			printUsage()
			os.Exit(0)
		}
	case CmdBlockchainAdd:
		if len(args) < 2 {
			printUsage()
			os.Exit(0)
		} else {
			BlockchainAddBlock(args[1])
		}
	default:
		printUsage()
	}

}

func cmdAddress(args []string) {
	if len(args) < 2 {
		printUsage()
		os.Exit(0)
	}

	switch args[0] {
	case CmdTransactionQuery:
		AddressQuery(args[1])
	case CmdTransactionSend:
		amount, _ := strconv.Atoi(args[1])
		AddressTransfer([]byte(args[3]), []byte(args[5]), amount)
	}
	log.Println(args)
}
