package utils

import (
	"crypto/rand"
	"log"
	"strings"
)

func RandString(letters []byte, n int) string {
	lettersLen := len(letters)
	if lettersLen > 255 {
		panic("letters length must less than 256")
	}

	var mask byte
	maskLen := 0
	for lettersLen > 0 {
		lettersLen = lettersLen >> 1
		maskLen++
	}
	mask = 1<<maskLen - 1
	lettersLen = len(letters)

	var res strings.Builder
	res.Grow(n)

	randBytes := make([]byte, n)

	for res.Len() < n {
		diff := n - res.Len()
		randBytes = randBytes[:diff]
		_, err := rand.Read(randBytes)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < diff; i++ {
			t := int(randBytes[i] & mask)
			if t < lettersLen {
				res.WriteByte(letters[t])
			}
		}
	}

	return res.String()
}
