package util

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const alphabet string = "abcdefghijklmnopqrstuvwxyz"

func RandomInt(min, max int) int {
	return min + int(rand.Int63n(int64(max-min)+1))
}

func RandomString(length int) string {
	var rs strings.Builder

	for i := 0; i < length; i++ {
		r := alphabet[rand.Intn(len(alphabet))]
		rs.WriteByte(r)
	}

	return rs.String()
}

// return random attributes for testing purpose
func RandomOwner() string {
	return RandomString(12)
}

func RandomBalance() int {
	return RandomInt(0, 100000)
}

func RandomCurrency() string {
	currency := []string{"USD","IDR","KRW","EUR","GBP"}

	r := rand.Intn(len(currency))
	return currency[r]
}