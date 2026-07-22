package util

import (
	"math/rand"
	"time"
)

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateUserID() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
