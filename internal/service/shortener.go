package service

import (
	"crypto/sha256"
)

func Shortener(url string) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	h := sha256.New()
	h.Write([]byte(url))
	hash256 := h.Sum(nil)

	var shortUrl string

	for i := 0; i < 10; i++ {
		shortUrl += string(chars[hash256[i]%63])
	}

	return shortUrl
}
