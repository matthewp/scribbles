package main

import (
	"encoding/hex"
	"log"

	"github.com/btcsuite/btcd/btcec"
)

func GetPubKey(privateKey string) string {
	if keyb, err := hex.DecodeString(config.PrivateKey); err != nil {
		log.Printf("Error decoding key from hex: %s\n", err.Error())
		return ""
	} else {
		_, pubkey := btcec.PrivKeyFromBytes(btcec.S256(), keyb)
		return hex.EncodeToString(pubkey.X.Bytes())
	}
}
