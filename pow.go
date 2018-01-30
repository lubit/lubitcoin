package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
)

var (
	maxNonce = math.MaxInt64
	maxZero  = 2
)

type ProofOfWork struct {
	block *Block
}

func (pow *ProofOfWork) Pad() []byte {
	src := bytes.Join(
		[][]byte{
			pow.block.PrevHash,
			pow.block.HashTransaction(),
			IntToByte(pow.block.Timestamp),
			IntToByte(pow.block.Nonce),
		},
		[]byte{})
	return src
}

func (pow *ProofOfWork) Mine() (int64, []byte) {
	var (
		hash  [32]byte
		nonce int64
	)
	for nonce < int64(maxNonce) {
		pow.block.Nonce = nonce
		hash = sha256.Sum256(pow.Pad())
		hashWin := true
		for i := 0; i < maxZero; i++ {
			if hash[i] != 0 {
				hashWin = false
				break
			}
		}
		fmt.Printf("\r nonce[%d], hash[%x]", nonce, hash)

		if hashWin {
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}

func IntToByte(num int64) []byte {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes()
}

/*
func main() {
	block := NewGenesisBlock()
	pow := ProofOfWork{block}
	pow.Mine()
	return
}
*/
