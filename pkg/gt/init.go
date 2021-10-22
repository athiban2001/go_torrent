package gt

import (
	"math/rand"
	"strings"
	"time"
)

func GetPeerID() string {
	prefix := "-GT0001-"
	possibleDigits := []rune("0123456789")

	rand.Seed(time.Now().UnixNano())
	b := strings.Builder{}

	for i := 1; i <= 12; i++ {
		b.WriteRune(possibleDigits[rand.Intn(len(possibleDigits))])
	}

	return prefix + b.String()
}
