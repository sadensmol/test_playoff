package utils

import (
	"crypto/rand"
	"math/big"
)

func GenRandomString(length int) (string, error) {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	randomString := ""

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			return "", err
		}
		randomString += string(characters[randomIndex.Int64()])
	}

	return randomString, nil
}
