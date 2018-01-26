package main

import (
	"fmt"
	"log"
	"sync"
)

// global variable
var (
	global_chain      *BlockChain
	global_chain_once sync.Once
)

// BlockchainGenesis : create the chain with genesis block
func BlockchainGenesis() {

	global_chain_once.Do(func() {
		global_chain = NewBlockChain(BlockchainBucket)
	})

}

// BlockchainListBlocks : list all the block on the chain
func BlockchainListBlocks() {
	if global_chain == nil {
		BlockchainGenesis()
	}
	global_chain.ListBlocks()

}

// BlockchainAddBlock : add a block at the end of the chain
func BlockchainAddBlock(info string) {
	if global_chain == nil {
		BlockchainGenesis()
	}
	log.Println("BlockchainAddBlock")
	b := NewBlock(info, nil, nil)
	global_chain.AddBlock(b)

}

// AddressQuery : return
func AddressQuery(addr string) int {
	if global_chain == nil {
		BlockchainGenesis()
	}
	txm, balance, err := global_chain.FindUTXOByAddress([]byte(addr), -1)
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
	if global_chain == nil {
		BlockchainGenesis()
	}
	txm, balance, err := global_chain.FindUTXOByAddress([]byte(from), amount)
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
	global_chain.AddBlock(block)

}

// UTXOReindex : reindex utxo
func UTXOReindex() {
	if global_chain == nil {
		BlockchainGenesis()
	}
	u := NewUTXOSet()
	//u := UTXOSet{global_chain}
	u.Reindex(*global_chain)
}

// UTXOQuery : utxo query
func UTXOQuery(addr []byte) int {
	if global_chain == nil {
		BlockchainGenesis()
	}
	u := NewUTXOSet()
	amount := u.GetBalance(addr)
	log.Println(string(addr), amount)
	return amount
}
