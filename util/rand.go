package util

import "crypto/rand"

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func RandomString(length uint) string {
	b := make([]byte, length)
	rand.Read(b)

	for i := uint(0); i < length; i++ {
		b[i] = alphabet[b[i]%byte(len(alphabet))]
	}

	return string(b)
}
