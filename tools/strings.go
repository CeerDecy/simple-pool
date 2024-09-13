package tools

import (
	"math/rand/v2"
	"time"
)

var CHARACTERS = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandString(length int) string {
	var bytes = make([]byte, length)
	rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
	for i := range bytes {
		bytes[i] = CHARACTERS[rand.IntN(len(CHARACTERS))%len(CHARACTERS)]
	}
	return string(bytes)
}
