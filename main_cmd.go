package main

import (
	"fmt"
	"log"
	"sync"
)

// global variable
var (
	luBc   *BlockChain
	luOnce sync.Once
)

// BlockchainGenesis : create the chain with genesis block
func BlockchainGenesis() {
	luOnce.Do(func() {
		luBc = NewBlockChain(BlockchainBucket)
		UTXOSet{luBc}.Reindex()
	})
}

// BlockchainListBlocks : list all the block on the chain
func BlockchainListBlocks() {
	if luBc == nil {
		BlockchainGenesis()
	}
	luBc.ListBlocks()

}

// BlockchainAddBlock : add a block at the end of the chain
func BlockchainAddBlock(info string) {
	if luBc == nil {
		BlockchainGenesis()
	}
	log.Println("BlockchainAddBlock")
	luBc.AddBlock(info, nil)
}

// AddressQuery : return
func AddressQuery(addr string) int {
	if luBc == nil {
		BlockchainGenesis()
	}
	txm, balance, err := luBc.FindUTXOByAddress([]byte(addr), -1)
	if err != nil {
		log.Println(err)
	}

	for k, v := range txm {
		log.Printf("TXID [%s] : TXAmount[%d] \n", k, v)
	}
	log.Println("Amount:", balance)
	return balance
}

// AddressTransfer : send from xx to yy
func AddressTransfer(from, to []byte, amount int) {
	log.Println(from, to, amount)
	if luBc == nil {
		BlockchainGenesis()
	}
	txm, balance, err := luBc.FindUTXOByAddress([]byte(from), amount)
	if err != nil {
		log.Println(err)
	}
	if balance < amount {
		log.Fatal("Not enough balance: ", balance)
		return
	}

	tx := NewTransaction(from, to, amount, txm)
	str := fmt.Sprintf("Transfer From [%s] TO [%s]: %d", string(from), string(to), amount)
	luBc.AddBlock(str, []Transaction{*tx})

}

// UTXOReindex : reindex utxo
func UTXOReindex() {
	if luBc == nil {
		BlockchainGenesis()
	}
	u := UTXOSet{luBc}
	u.Reindex()
}

// UTXOQuery : utxo query
func UTXOQuery(addr []byte) int {
	if luBc == nil {
		BlockchainGenesis()
	}
	u := UTXOSet{luBc}
	amount := u.GetBalance(addr)
	log.Println(string(addr), amount)
	return amount
}
