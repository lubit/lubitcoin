package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Block struct
type Block struct {
	Timestamp    int64
	Extras       string
	PrevHash     []byte
	CurrHash     []byte
	Transactions []Transaction // add for TX
}

// NewBlock return a new block
func NewBlock(extras string, prevHash []byte, txs []Transaction) *Block {
	hash := sha256.Sum256([]byte(extras))
	block := &Block{
		Timestamp:    time.Now().Unix(),
		Extras:       extras,
		PrevHash:     prevHash,
		CurrHash:     hash[:],
		Transactions: txs,
	}
	return block
}

// DeserializeBlock decode a block from []byte
func DeserializeBlock(enc []byte) *Block {
	var block Block
	dec := gob.NewDecoder(bytes.NewReader(enc))
	if err := dec.Decode(&block); err != nil {
		log.Panic(err)
	}
	return &block
}

// Serialize the block
func (b *Block) Serialize() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(b); err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// Dump pretty print the block
func (b *Block) Dump() {
	dump, _ := json.MarshalIndent(b, "", "  ")
	fmt.Println(string(dump))
}
