package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

/*Cmd Varaiable*/
const (
	CmdBlockchain          = "blockchain"
	CmdBlockchainCreate    = "create"
	CmdBlockchainList      = "list"
	CmdBlockchainListBlock = "block"
	CmdBlockchainAdd       = "add"
	CmdAddress             = "address"
	CmdAddressQuery        = "query"
	CmdAddressTransfer     = "transfer"
)

func main() {

	parseFlags()
	return
}

func printUsage() {

	usage := `	
	USAGE : lubitcoin <module> [<action>] [action parameter]

	lubitcoin <blockchain> [xxx]
		lubitcoin blockchain create : create blockchain with genesis block
		lubitcoin blockchain list block : list blockchain all blocks
		lubitcoin blockchain add xxxx : add a block

	lubitcoin <address> [xx]
		lubitcoin address query xxx : query address balance
		lubitcoin address transfer 'amount_00' from 'address_xx' to 'address_yy' : transfer

	`
	fmt.Println(usage)
}

func parseFlags() {

	flagSet := flag.NewFlagSet("lubitcoin", flag.ExitOnError)
	flagSet.Parse(os.Args)
	if flagSet.NArg() < 2 {
		printUsage()
		os.Exit(0)
	}
	switch flagSet.Arg(1) {
	case CmdBlockchain:
		cmdBlockChain(flagSet.Args()[2:])
	case CmdAddress:
		cmdAddress(flagSet.Args()[2:])
	default:
		cmdBlockChain(flagSet.Args()[2:])
	}
}

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
	case CmdAddressQuery:
		AddressQuery(args[1])
	case CmdAddressTransfer:
		amount, _ := strconv.Atoi(args[1])
		AddressTransfer([]byte(args[3]), []byte(args[5]), amount)
	}
	log.Println(args)
}
