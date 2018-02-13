package main

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

// variable
const (
	UTXOFile = "lubit.db.utxo"
)

// UTXOSet  owns all the UTXOs
type UTXOSet struct {
	lvl   *leveldb.DB
	chain *BlockChain
}

func NewUTXOSet(bc *BlockChain) *UTXOSet {
	lvl, _ := leveldb.OpenFile(UTXOFile, nil)
	set := &UTXOSet{
		lvl:   lvl,
		chain: bc,
	}
	return set
}

func (u UTXOSet) Reindex() {

	// delete
	iter := u.lvl.NewIterator(nil, nil)
	for iter.Next() {
		u.lvl.Delete(iter.Key(), nil)
	}
	// rebuild
	UTXOS := u.chain.FindUTXO()
	for k, v := range UTXOS {
		key, _ := hex.DecodeString(k)
		val := TXOutputs{v}
		u.lvl.Put(key, val.Serialize(), nil)
	}
}

func (u UTXOSet) Update(b *Block) {

	for _, tx := range b.Transactions {
		// inputs
		for _, ins := range tx.TXInputs {

			val, err := u.lvl.Get(ins.TXID, nil)
			if err != nil { //没找到
				log.Println(err)
				continue
			}
			u.lvl.Delete(ins.TXID, nil)
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

// FindUTXOByAddress get UTXO by address
func (u UTXOSet) FindUTXO(addr []byte, amount int) (int, map[string][]TXOutput) {
	var acc int
	utxo := make(map[string][]TXOutput)
	iter := u.lvl.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		txid := hex.EncodeToString(iter.Key())
		tx := DeserializeTXO(iter.Value())
		var txos []TXOutput
		for _, txo := range tx.TXOS {
			if bytes.EqualFold(addr, txo.ScriptPubKey) {
				acc += txo.Amount
				txos = append(txos, txo)
				if acc > amount {
					return acc, utxo
				}
			}
		}
		if len(txos) > 1 {
			utxo[txid] = txos
		}
	}
	return acc, utxo
}
