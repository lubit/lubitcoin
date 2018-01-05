package main

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"
)

//Blockchain Variable
const (
	BlockchainName   = "lubit"
	BlockchainFile   = "lubit.db"
	BlockchainBucket = "lubit.bucket"
	BlockChainLast   = "l"
)

// BlockChain struct
type BlockChain struct {
	name string
	last []byte
	db   *bolt.DB
}

// NewBlockChain return a new block chain
func NewBlockChain(str string) *BlockChain {

	// open storage
	db, err := bolt.Open(BlockchainFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	var hash []byte
	// Create or Get Genesis Block
	db.Update(func(tx *bolt.Tx) error {
		dbk, err := tx.CreateBucketIfNotExists([]byte(str))
		if err != nil {
			log.Panic(err)
			return err
		}
		hash = dbk.Get([]byte(BlockChainLast))
		if nil == hash {
			//creage genesis transaction
			tx := NewGenesisTransaction([]byte(GenesisAuthor))
			//create genesis block
			block := NewBlock("lubitcoin genesis block", nil, []Transaction{*tx})

			hash = block.CurrHash
			dbk.Put(hash, block.Serialize())
			log.Println("Create Blockchain with genesis block")
			dbk.Put([]byte(BlockChainLast), hash)
		}
		return nil
	})
	// Construct blockchain
	bc := &BlockChain{
		name: str,
		last: hash,
		db:   db,
	}
	return bc
}

// AddBlock bc add a block
func (bc *BlockChain) AddBlock(str string, txs []Transaction) {

	block := NewBlock(str, bc.last, txs)
	block.Dump()

	bc.db.Update(func(tx *bolt.Tx) error {
		// db store block
		dbk := tx.Bucket([]byte(bc.name))
		dbk.Put(block.CurrHash, block.Serialize())
		//update last
		dbk.Put([]byte(BlockChainLast), block.CurrHash)
		bc.last = block.CurrHash
		return nil
	})

}

// ListBlocks bc list blocks
func (bc *BlockChain) ListBlocks() {

	bc.db.View(func(tx *bolt.Tx) error {
		dbk := tx.Bucket([]byte(bc.name))
		iter := bc.last
		for {
			enc := dbk.Get(iter)
			block := DeserializeBlock(enc)
			block.Dump()
			if block.PrevHash == nil {
				return nil
			}
			iter = block.PrevHash
		}
	})
}

// FindUTXO iterate address amount
func (bc *BlockChain) FindUTXO(addr []byte, amount int) (map[string]int, int, error) {

	balance := 0
	UTXO := make(map[string]int)  // Unspent Transaction Output
	STXI := make(map[string]bool) // spent transaction input
	bc.db.View(func(tx *bolt.Tx) error {
		dbk := tx.Bucket([]byte(bc.name))
		iter := bc.last
		for {
			enc := dbk.Get(iter)
			block := DeserializeBlock(enc)

			txs := block.Transactions
			for _, tx := range txs {
				txid := hex.EncodeToString(tx.TXID)
				if (amount != -1) && (balance >= amount) {
					break
				}
				// check previous inputs
				if _, ok := STXI[txid]; ok {
					continue
				}
				// TXOUT
				for _, out := range tx.TXOutputs {
					// check address
					if !bytes.Equal(out.Address, addr) {
						continue
					}
					UTXO[txid] = out.Amount
					balance += out.Amount
				}
				// TXINPUT
				for _, in := range tx.TXInputs {
					if !bytes.Equal(in.Address, addr) {
						continue
					}
					id := hex.EncodeToString(in.TXID)
					STXI[id] = true
				}
			}
			// genis block break
			if block.PrevHash == nil {
				break
			}
			iter = block.PrevHash
		}
		return nil
	})

	return UTXO, balance, nil
}
