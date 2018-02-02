package main

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

//Blockchain Variable
const (
	BlockChainLast = "tip"
	BlockchainFile = "lubit.db.block"
)

// BlockChain struct
type BlockChain struct {
	tip []byte
	lvl *leveldb.DB
}

// NewBlockChain return a new block chain
func NewBlockChain() *BlockChain {

	lvl, err := leveldb.OpenFile(BlockchainFile, nil)
	if err != nil {
		log.Panic(err)
	}
	tip, err := lvl.Get([]byte(BlockChainLast), nil)
	if err != nil {
		// empty chain
		if leveldb.ErrNotFound == err {
			block := NewGenesisBlock()
			tip = block.CurrHash
			lvl.Put(block.CurrHash, block.Serialize(), nil)
			lvl.Put([]byte(BlockChainLast), block.CurrHash, nil)
		} else {
			log.Panic("LevelDB get tip failed")
		}
	}
	chain := &BlockChain{
		tip: tip,
		lvl: lvl,
	}

	return chain
}

// NewGenesisBlock for blockchain
func NewGenesisBlock() *Block {
	tx := NewGenesisTransaction()
	block := NewBlock("lubitcoin genesis block", nil, []Transaction{*tx})
	return block
}

// AddBlock chain add new block
func (chain *BlockChain) AddBlock(block *Block) {
	block.PrevHash = chain.tip
	chain.lvl.Put(block.CurrHash, block.Serialize(), nil)
	chain.lvl.Put([]byte(BlockChainLast), block.CurrHash, nil)
	chain.tip = block.CurrHash
}

// ListBlocks list&dump block
func (chain *BlockChain) ListBlocks() {
	iter := chain.tip
	for {
		if iter == nil {
			return
		}
		enc, _ := chain.lvl.Get(iter, nil)
		block := DeserializeBlock(enc)
		block.Dump()
		if block.PrevHash == nil {
			return
		}
		iter = block.PrevHash
	}
}

// FindUTXOByAddress iterate address amount
func (chain *BlockChain) FindUTXOByAddress(addr []byte, amount int) (map[string]int, int, error) {
	balance := 0
	UTXO := make(map[string]int)  // Unspent Transaction Output
	STXI := make(map[string]bool) // spent transaction input

	iter := chain.tip
	for {

		enc, _ := chain.lvl.Get([]byte(iter), nil)
		block := DeserializeBlock(enc)
		txs := block.Transactions
		for _, tx := range txs {
			txid := hex.EncodeToString(tx.TXID)
			if (amount != -1) && (balance >= amount) {
				break
			}
			// TXOUT : check previous inputs
			if _, exist := STXI[txid]; exist {
				log.Println("STXI exist", STXI, txid)
			} else {
				for _, out := range tx.TXOutputs {
					// check address
					if !bytes.Equal([]byte(out.Address), addr) {
						continue
					}
					UTXO[txid] = out.Amount
					balance += out.Amount
				}
			}
			// TXINPUT
			for _, in := range tx.TXInputs {
				id := hex.EncodeToString(in.TXID)
				if false == bytes.Equal([]byte(in.Address), addr) {
					continue
				} else {
					id = hex.EncodeToString(in.TXID)
					STXI[id] = true
				}
			}
		}
		// genis block break
		if block.PrevHash == nil {
			log.Println("genesis block arrived")
			break
		}
		iter = block.PrevHash
	}
	return UTXO, balance, nil
}

// FindUTXO return all the UTXO
func (chain *BlockChain) FindUTXO() map[string][]TXOutput {
	UTXO := make(map[string][]TXOutput)
	STXI := make(map[string][]string)

	iter := chain.tip
	for {
		enc, err := chain.lvl.Get([]byte(iter), nil)
		if err != nil {
			log.Panic(err)
		}
		block := DeserializeBlock(enc)
		for _, tx := range block.Transactions {
			txid := hex.EncodeToString(tx.TXID)
			// UTXO collect
			for _, out := range tx.TXOutputs {

				exist := false
				if STXI[txid] != nil {
					for _, addr := range STXI[txid] {
						if addr == out.Address {
							exist = true
						}
					}
				}
				if !exist {
					UTXO[txid] = append(UTXO[txid], out)
					log.Printf("UTXO: %+v \n", UTXO)
				}
			}
			// STXI
			for _, in := range tx.TXInputs {
				id := hex.EncodeToString(in.TXID)
				STXI[id] = append(STXI[id], in.Address)
			}
		}
		if block.PrevHash == nil {
			break
		} else {
			iter = block.PrevHash
		}
	}

	return UTXO
}

func (chain *BlockChain) GenerateUTXO() {

}
