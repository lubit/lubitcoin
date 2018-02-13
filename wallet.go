package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const WalletFile = "lubit.wallet"

type Wallet struct {
	KeyPairs map[string]*KeyPair
}

func NewWallet() *Wallet {
	w := &Wallet{
		KeyPairs: make(map[string]*KeyPair),
	}
	w.LoadFromFile()
	return w
}

func (w Wallet) GenKeyPair() string {
	key := NewKeyPair()
	addr := key.GetAddress()
	w.KeyPairs[addr] = key
	return addr
}

func (w Wallet) SaveToFile() {
	var buff bytes.Buffer
	gob.Register(elliptic.P256())
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(w); err != nil {
		log.Panic(err)
	}
	if err := ioutil.WriteFile(WalletFile, buff.Bytes(), 0644); err != nil {
		log.Panic(err)
	}
}

func (w Wallet) LoadFromFile() {
	// not exist
	if _, err := os.Stat(WalletFile); os.IsNotExist(err) {
		w.SaveToFile()
		return
	}
	// load
	content, err := ioutil.ReadFile(WalletFile)
	if err != nil {
		log.Panic(err)
	}
	gob.Register(elliptic.P256())
	dec := gob.NewDecoder(bytes.NewReader(content))
	err = dec.Decode(&w)
	if err != nil {
		log.Panic(err)
	}
	return
}

func (w Wallet) Dump() {
	log.Println("Wallet Dump :")
	for k, _ := range w.KeyPairs {
		log.Println("\t Address:", k)
	}
}

func (w Wallet) GetKeyPair(addr string) *KeyPair {
	if _, ok := w.KeyPairs[addr]; ok {
		return w.KeyPairs[addr]
	} else {
		return nil
	}
}

type KeyPair struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewKeyPair() *KeyPair {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	keypair := &KeyPair{
		PrivateKey: *private,
		PublicKey:  public,
	}
	return keypair
}

func (kp KeyPair) GetAddress() string {
	globalVersion := byte(0x00)
	pubhash := HashPublicKey(kp.PublicKey)
	version := globalVersion
	versionPayload := append([]byte{version}, pubhash...)
	checksum := checksum(versionPayload)
	finalPayload := append(versionPayload, checksum...)
	enc := base58.Encode(finalPayload)

	return enc
}

func HashPublicKey(pubkey []byte) []byte {
	sha := sha256.Sum256(pubkey)
	ripe := ripemd160.New()
	_, err := ripe.Write(sha[:])
	if err != nil {
		log.Fatal(err)
	}
	pub := ripe.Sum(nil)
	return pub
}

func checksum(payload []byte) []byte {
	s := sha256.Sum256(payload)
	ss := sha256.Sum256(s[:])
	return ss[:4]
}

func main() {
	w := NewWallet()
	w.Dump()
	w.GenKeyPair()
	w.Dump()
	w.SaveToFile()
	return
	/*
			priv, _ := GenerateKey(c, rand.Reader)

		  	hashed := []byte("testing")
		  	r, s, err := Sign(rand.Reader, priv, hashed)
		  	if err != nil {
		  		t.Errorf("%s: error signing: %s", tag, err)
		  		return
		  	}

		  	if !Verify(&priv.PublicKey, hashed, r, s) {
		  		t.Errorf("%s: Verify failed", tag)
		  	}

		  	hashed[0] ^= 0xff
		  	if Verify(&priv.PublicKey, hashed, r, s) {
		  		t.Errorf("%s: Verify always works!", tag)
			  }
	*/
}
