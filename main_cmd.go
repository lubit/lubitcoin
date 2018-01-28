package main

import (
	"fmt"
	"log"
	"sync"
)

// global variable
var (
	globalChain     *BlockChain
	globalUTXO 		*UTXOSet
	globalChainOnce sync.Once
)

// BlockchainGenesis : create the chain with genesis block
func BlockchainGenesis() {
	globalChainOnce.Do(func() {
		globalChain = NewBlockChain()
		globalUTXO = NewUTXOSet(globalChain)
	})
}

// BlockchainListBlocks : list all the block on the chain
func BlockchainListBlocks() {
	if globalChain == nil {
		BlockchainGenesis()
	}
	globalChain.ListBlocks()

}

// BlockchainAddBlock : add a block at the end of the chain
func BlockchainAddBlock(info string) {
	if globalChain == nil {
		BlockchainGenesis()
	}
	log.Println("BlockchainAddBlock")
	b := NewBlock(info, nil, nil)
	globalChain.AddBlock(b)

}

// AddressQuery : return
func AddressQuery(addr string) int {
	if globalChain == nil {
		BlockchainGenesis()
	}
	txm, balance, err := globalChain.FindUTXOByAddress([]byte(addr), -1)
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
	if globalChain == nil {
		BlockchainGenesis()
	}
	txm, balance, err := globalChain.FindUTXOByAddress([]byte(from), amount)
	if err != nil {
		log.Println(err)
	}
	if balance < amount {
		log.Fatal("Not enough balance: ", balance)
		return
	}

	tx := NewTransaction(from, to, amount, txm)
	str := fmt.Sprintf("Transfer From [%s] TO [%s]: %d", string(from), string(to), amount)
	block := NewBlock(str, nil, []Transaction{*tx})
	globalChain.AddBlock(block)

}

// UTXOReindex : reindex utxo
func UTXOReindex() {
	if globalUTXO == nil {
		BlockchainGenesis()
	}
	globalUTXO.Reindex()
}

// UTXOQuery : utxo query
func UTXOQuery(addr []byte) int {
	if globalUTXO == nil {
		BlockchainGenesis()
	}
	amount :=globalUTXO.GetBalance(addr)
	log.Println(string(addr), amount)
	return amount
}
