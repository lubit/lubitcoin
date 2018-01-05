package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

// GenesisRewards 100
const (
	GenesisRewards = 100
	GenesisAuthor  = "luofaxuan"
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
	Address []byte
}

// TXOutput struct
type TXOutput struct {
	Amount  int
	Address []byte
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
			Address: from,
		}
		txins = append(txins, txin)
		total += v
	}
	//OUTPUTS
	txout := TXOutput{
		Amount:  amount,
		Address: to,
	}
	txouts = append(txouts, txout)
	if total-amount > 0 {
		txout = TXOutput{
			Amount:  total - amount,
			Address: from,
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
		Address: addr,
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
