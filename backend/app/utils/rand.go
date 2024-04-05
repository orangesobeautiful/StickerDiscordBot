package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"strings"
)

func RandString(letters []byte, n int) string {
	const maxLettersLen = 255
	lettersLen := len(letters)
	if lettersLen > maxLettersLen {
		panic(fmt.Sprintf("letters length must less than %d", maxLettersLen+1))
	}

	var mask byte
	maskLen := 0
	for lettersLen > 0 {
		lettersLen >>= 1
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
