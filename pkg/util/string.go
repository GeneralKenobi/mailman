package util

import (
	"math/rand"
	"strings"
)

// RandomAlphanumericString generates a string of given length composed of random alphanumeric characters.
func RandomAlphanumericString(length int) string {
	if length < 0 {
		panic("length cannot be 0")
	}

	var builder strings.Builder
	for i := 0; i < length; i++ {
		const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		builder.WriteByte(alphanumeric[rand.Intn(len(alphanumeric))])
	}
	return builder.String()
}
