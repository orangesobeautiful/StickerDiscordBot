package utils_test

import (
	"crypto/rand"
	"log"
	"math/big"
	"strings"
	"testing"

	"backend/utils"
)

var benchLetters = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func RandInt(letters []byte, n int) string {
	lettersLen := len(letters)

	var res strings.Builder
	res.Grow(n)
	for i := 0; i < n; i++ {
		bn, err := rand.Int(rand.Reader, big.NewInt(int64(lettersLen)))
		if err != nil {
			log.Fatal(err)
		}
		res.WriteByte(letters[bn.Int64()])
	}

	return res.String()
}

func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandInt(benchLetters, 10)
	}
}

func BenchmarkRandRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.RandString(benchLetters, 10)
	}
}
