package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

// GenesisRewards 100
const (
	GenesisRewards  = 100 // 1 lubit = 1 * 10^6 发
	GenesisAuthor   = "thues"
	CoinbaseSubsidy = 10
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
	TXID      []byte
	Amount    int
	Address   string
	Signature []byte
	PubKey    []byte
}

// TXOutput struct
type TXOutput struct {
	Amount       int
	Address      string
	ScriptPubKey []byte
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
func NewGenesisTransaction() *Transaction {
	txout := TXOutput{
		Amount:  GenesisRewards,
		Address: GenesisAuthor,
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

func (tx Transaction) Hash() []byte {
	b, _ := json.Marshal(tx)
	hash := sha256.Sum256(b)
	return hash[:]
}

// NewTxCoinbase create coinbase by mine block
func NewCoinbaseTx(to_addr string) *Transaction {

	script := []byte(to_addr)

	txin := TXInput{[]byte{}, -1, "", nil, nil}
	txout := TXOutput{
		Amount:       CoinbaseSubsidy,
		Address:      "",
		ScriptPubKey: script,
	}
	tx := Transaction{
		TXInputs:  []TXInput{txin},
		TXOutputs: []TXOutput{txout},
		Timestamp: time.Now(),
	}
	tx.TXID = tx.Hash()

	return &tx
}

// NewScriptPubTx create Transaction with ScriptPub
func NewSignedTx(kp *KeyPair, to string, amount, remain int, utxo map[string][]TXOutput) *Transaction {
	var txins []TXInput
	var txouts []TXOutput
	// input
	for k, txos := range utxo {
		txid, _ := hex.DecodeString(k)
		for _, out := range txos {
			input := TXInput{
				TXID:   txid,
				Amount: out.Amount,
				PubKey: kp.PublicKey,
			}
			txins = append(txins, input)
		}
	}
	// output
	addr := base58.Decode(to)
	tohash := addr[1 : len(addr)-4]
	output := TXOutput{amount, "", tohash}
	txouts = append(txouts, output)
	if remain > 0 {
		rem := TXOutput{remain, "", HashPublicKey(kp.PublicKey)}
		txouts = append(txouts, rem)
	}
	// todo

	tx := &Transaction{
		TXInputs:  txins,
		TXOutputs: txouts,
		Timestamp: time.Now(),
	}
	tx.TXID = tx.Hash()
	tx.TxSign(kp)
	return tx
}

func (tx Transaction) TxVerify() bool {
	dup := tx.TxDuplicate()
	curve := elliptic.P256()
	for idx, vin := range tx.TXInputs {
		// setup src
		dup.TXInputs[idx].PubKey = vin.PubKey
		signSrc := fmt.Sprintf("%x", dup)
		// setup key
		r, s, x, y := big.Int{}, big.Int{}, big.Int{}, big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		pubLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(pubLen / 2)])
		y.SetBytes(vin.PubKey[(pubLen / 2):])
		pubkey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		// verify
		if !ecdsa.Verify(&pubkey, []byte(signSrc), &r, &s) {
			log.Printf("txin[%d] signature verify failed \n", idx)
			return false
		}
	}
	return true
}

func (tx Transaction) TxSign(kp *KeyPair) {

	dup := tx.TxDuplicate()
	for _, vin := range dup.TXInputs {
		vin.PubKey = kp.PublicKey
		signSrc := fmt.Sprintf("%x", dup)
		r, s, err := ecdsa.Sign(rand.Reader, &kp.PrivateKey, []byte(signSrc))
		if err != nil {
			log.Panic(err)
		}
		sig := append(r.Bytes(), s.Bytes()...)
		vin.Signature = sig
	}

	return
}

func (tx Transaction) TxDuplicate() Transaction {

	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.TXInputs {
		inputs = append(inputs, TXInput{vin.TXID, vin.Amount, "", nil, nil})
	}
	for _, vout := range tx.TXOutputs {
		outputs = append(outputs, TXOutput{vout.Amount, "", vout.ScriptPubKey})
	}
	dup := Transaction{tx.TXID, inputs, outputs, time.Now()}

	return dup
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

/* TODO
1. coinbase
2. sign
3. verify
*/
