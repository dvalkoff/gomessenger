package utils

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
)

var curve = ecdh.X25519()
var random = rand.Reader

func GenerateKeys(len int) ([]*ecdh.PrivateKey, error) {
	if len <= 0 {
		return nil, fmt.Errorf("len must be bigger than 0")
	}
	keys := make([]*ecdh.PrivateKey, 0, len)
	for range len {
		key, err := GenerateKey()
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func GenerateKey() (*ecdh.PrivateKey, error) {
	return curve.GenerateKey(random)
}

func MapKey(privateKey []byte) (*ecdh.PrivateKey, error) {
	return curve.NewPrivateKey(privateKey)
}
