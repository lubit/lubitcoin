package main

import (
	"bytes"
	"log"

	"github.com/boltdb/bolt"
)

// UTXOSet  owns all the UTXOs
type UTXOSet struct {
	BlockChain *BlockChain // For BlockChain.db
}

// Reindex build UTXOSet From Blockchain
func (u UTXOSet) Reindex() {
	db := u.BlockChain.db
	//Rebuild all UTXO
	UTXOS := u.BlockChain.FindUTXO()
	//ReCreate the bucket
	err := db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(BlockchainUTXO))
		buck, _ := tx.CreateBucket([]byte(BlockchainUTXO))

		for k, v := range UTXOS {
			val := TXOutputs{v}
			buck.Put([]byte(k), val.Serialize())
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

}

// Update the UTXOSet
func (u UTXOSet) Update(b *Block) {
	db := u.BlockChain.db
	err := db.Update(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(BlockchainUTXO))

		for _, tx := range b.Transactions {
			// Inputs
			for _, ins := range tx.TXInputs {
				if val := buck.Get(ins.TXID); val != nil {
					// delete
					buck.Delete(ins.TXID)
					// update
					txos := DeserializeTXO(val)
					var newo []TXOutput
					for _, v := range txos.TXOS {
						if ins.Address != v.Address {
							newo = append(newo, v)
						}
					}
					utxo := TXOutputs{newo}
					buck.Put(ins.TXID, utxo.Serialize())

				}
			}
			// Outputs
			utxo := TXOutputs{tx.TXOutputs}
			buck.Put(tx.TXID, utxo.Serialize())
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// FindUTXOByAddress get UTXO by address
func (u UTXOSet) FindUTXOByAddress(addr []byte) []TXOutput {
	var utxo []TXOutput
	db := u.BlockChain.db
	err := db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(BlockchainUTXO))
		curs := buck.Cursor()
		for k, v := curs.First(); k != nil; k, v = curs.Next() {
			val := DeserializeTXO(v)
			for _, txo := range val.TXOS {
				if bytes.EqualFold([]byte(txo.Address), addr) {
					utxo = append(utxo, txo)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return utxo
}

// GetBalance get balance by address
func (u UTXOSet) GetBalance(addr []byte) int {
	amount := 0
	txos := u.FindUTXOByAddress(addr)
	for _, v := range txos {
		amount += v.Amount
	}
	return amount
}
