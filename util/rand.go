package util

import "crypto/rand"

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)

	for i := 0; i < length; i++ {
		b[i] = alphabet[b[i]%byte(len(alphabet))]
	}

	return string(b)
}
