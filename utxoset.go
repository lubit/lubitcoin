package main

import (
	"bytes"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	UTXOFile = "lubit.utxo"
)

// UTXOSet  owns all the UTXOs
type UTXOSet struct {
	lvl *leveldb.DB
	//BlockChain *BlockChain // For BlockChain.db
}

func NewUTXOSet() *UTXOSet {
	lvl, _ := leveldb.OpenFile(UTXOFile, nil)
	set := &UTXOSet{lvl}
	return set
}

func (u UTXOSet) Reindex(bc BlockChain) {
	UTXOS := bc.FindUTXO()
	for k, v := range UTXOS {
		val := TXOutputs{v}
		u.lvl.Put([]byte(k), val.Serialize(), nil)
	}
}

func (u UTXOSet) Update(b *Block) {
	for _, tx := range b.Transactions {
		// inputs
		for _, ins := range tx.TXInputs {
			val, err := u.lvl.Get([]byte(ins.TXID), nil)
			if err != nil { //没找到
				log.Println(err)
				continue
			}
			u.lvl.Delete([]byte(ins.TXID), nil)
			var oxv []TXOutput
			txos := DeserializeTXO(val)
			for _, v := range txos.TXOS {
				if ins.Address != v.Address {
					oxv = append(oxv, v)
				}
			}
			u.lvl.Put(ins.TXID, TXOutputs{oxv}.Serialize(), nil)
		}
		// outpus
		u.lvl.Put(tx.TXID, TXOutputs{tx.TXOutputs}.Serialize(), nil)
	}
}

// FindUTXOByAddress get UTXO by address
func (u UTXOSet) FindUTXOByAddress(addr []byte) []TXOutput {
	var utxo []TXOutput
	iter := u.lvl.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		tx := DeserializeTXO(iter.Value())
		for _, txo := range tx.TXOS {
			if bytes.EqualFold(addr, []byte(txo.Address)) {
				utxo = append(utxo, txo)
			}
		}
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
