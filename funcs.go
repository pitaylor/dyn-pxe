package main

import (
	"crypto/sha256"
	"encoding/hex"
)

func ShortHash(s string) string {
	hash := sha256.Sum256([]byte(s))
	str := hex.EncodeToString(hash[:])
	return str[:7]
}
