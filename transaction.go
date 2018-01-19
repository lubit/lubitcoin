package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"
)

// GenesisRewards 100
const (
	GenesisRewards = 100 // 1 lubit = 1 * 10^6 发
	GenesisAuthor  = "罗发宣"
)

// Transaction struct
type Transaction struct {
	TXID      []byte
	TXInputs  []TXInput
	TXOutputs []TXOutput
	Timestamp time.Time
}

// TXInput struct
type TXInput struct {
	TXID    []byte
	Amount  int
	Address string
}

// TXOutput struct
type TXOutput struct {
	Amount  int
	Address string
}

// NewTransaction create a new TX
func NewTransaction(from, to []byte, amount int, utxo map[string]int) *Transaction {

	var (
		txins  []TXInput
		txouts []TXOutput
	)
	total := 0
	//INPUTS
	for k, v := range utxo {
		id, _ := hex.DecodeString(k)
		txin := TXInput{
			TXID:    id,
			Amount:  v,
			Address: string(from),
		}
		txins = append(txins, txin)
		total += v
	}
	//OUTPUTS
	txout := TXOutput{
		Amount:  amount,
		Address: string(to),
	}
	txouts = append(txouts, txout)
	if total-amount > 0 {
		txout = TXOutput{
			Amount:  total - amount,
			Address: string(from),
		}
		txouts = append(txouts, txout)
	}
	//CONSTRUCT
	tx := &Transaction{
		TXInputs:  txins,
		TXOutputs: txouts,
		Timestamp: time.Now(),
	}

	b, _ := json.Marshal(tx)
	hash := sha256.Sum256(b)
	tx.TXID = hash[:]
	return tx

}

// NewGenesisTransaction create genesis transaction
func NewGenesisTransaction(addr []byte) *Transaction {
	txout := TXOutput{
		Amount:  GenesisRewards,
		Address: string(addr),
	}
	tx := &Transaction{
		TXInputs:  nil,
		TXOutputs: []TXOutput{txout},
		Timestamp: time.Now(),
	}
	b, _ := json.Marshal(tx)
	hash := sha256.Sum256(b)
	tx.TXID = hash[:]
	return tx
}

// TXOutputs struct
type TXOutputs struct {
	TXOS []TXOutput
}

// Serialize to bytes
func (o TXOutputs) Serialize() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(o); err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// Deserialize to struct
func DeserializeTXO(enc []byte) *TXOutputs {
	var txo TXOutputs
	dec := gob.NewDecoder(bytes.NewReader(enc))
	if err := dec.Decode(&txo); err != nil {
		log.Panic(err)
	}
	return &txo
}
