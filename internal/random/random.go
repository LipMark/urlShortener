package random

import (
	"math/rand"
	"time"
)

// NewRandomAlias generates random string with given size. It is used when desired alias is empty.
func NewRandomAlias(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	allowedChars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	randomAlias := make([]rune, size)
	for i := range randomAlias {
		randomAlias[i] = allowedChars[rnd.Intn(len(allowedChars))]
	}

	return string(randomAlias)
}
