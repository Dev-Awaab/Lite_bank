package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer 
func RandomInt(min, max int64) int64{
	return min + rand.Int63n(max - min + 1)
}

//RandomString generate a random string of length n
func RandomString(n int) string {
    var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < k; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generate a random owner
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generate a random amount of money
func RandomMoney() int64 {
	return	RandomInt(0, 1000)
}

// RandomCurrency generate a random curreny
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

