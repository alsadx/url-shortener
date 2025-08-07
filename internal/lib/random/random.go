package random

import (
	"crypto/rand"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func NewRandomString(length int) (string, error) {
	res := make([]byte, length)

	for i := range res {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		res[i] = charset[n.Int64()]
	}

	return string(res), nil
}
