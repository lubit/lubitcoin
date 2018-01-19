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
	BlockchainBucket = "lubit.bucket" //block
	BlockchainUTXO   = "lubit.utxo"   //chainstate
	BlockChainLast   = "l"
)

// BlockChain struct
type BlockChain struct {
	name string
	tip  []byte
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
		tip:  hash,
		db:   db,
	}
	return bc
}

// AddBlock bc add a block
func (bc *BlockChain) AddBlock(str string, txs []Transaction) {

	block := NewBlock(str, bc.tip, txs)
	block.Dump()

	bc.db.Update(func(tx *bolt.Tx) error {
		// db store block
		dbk := tx.Bucket([]byte(bc.name))
		dbk.Put(block.CurrHash, block.Serialize())
		//update last
		dbk.Put([]byte(BlockChainLast), block.CurrHash)
		return nil
	})

	bc.tip = block.CurrHash

}

// ListBlocks bc list blocks
func (bc *BlockChain) ListBlocks() {

	bc.db.View(func(tx *bolt.Tx) error {
		dbk := tx.Bucket([]byte(bc.name))
		iter := bc.tip
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

// FindUTXOByAddress iterate address amount
func (bc *BlockChain) FindUTXOByAddress(addr []byte, amount int) (map[string]int, int, error) {
	balance := 0
	UTXO := make(map[string]int)  // Unspent Transaction Output
	STXI := make(map[string]bool) // spent transaction input
	bc.db.View(func(tx *bolt.Tx) error {
		dbk := tx.Bucket([]byte(bc.name))
		//iter := bc.last
		iter := dbk.Get([]byte(BlockChainLast))
		for {
			enc := dbk.Get(iter)
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
		return nil
	})

	return UTXO, balance, nil
}

// FindUTXO return all the UTXO
func (bc *BlockChain) FindUTXO() map[string][]TXOutput {
	UTXO := make(map[string][]TXOutput)
	STXI := make(map[string][]string)
	bc.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(bc.name))
		iter := buck.Get([]byte(BlockChainLast))
		for {
			block := DeserializeBlock(buck.Get(iter))
			for _, tx := range block.Transactions {
				txid := hex.EncodeToString(tx.TXID)
				// UTXO collect
				for _, out := range tx.TXOutputs {
					// check if spent in STXI
					if STXI[txid] != nil {
						for _, addr := range STXI[txid] {
							if addr == out.Address {
								continue
							}
						}
					}
					UTXO[txid] = append(UTXO[txid], out)
				}
				// STXI collect
				for _, in := range tx.TXInputs {
					inid := hex.EncodeToString(in.TXID)
					STXI[inid] = append(STXI[inid], in.Address)
				}
			}
			//
			if block.PrevHash == nil {
				break
			} else {
				iter = block.PrevHash
			}
		}
		return nil
	})
	return nil
}
