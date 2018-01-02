package main

import (
	"log"
	"sync"

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
			//create genesis block
			block := NewBlock("lubitcoin genesis block", nil)
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
func (bc *BlockChain) AddBlock(str string) {

	block := NewBlock(str, bc.last)
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

// global variable
var (
	luBc   *BlockChain
	luOnce sync.Once
)

// BlockchainGenesis : create the chain with genesis block
func BlockchainGenesis() {

	luOnce.Do(func() {
		luBc = NewBlockChain(BlockchainBucket)
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
	luBc.AddBlock(info)

}
