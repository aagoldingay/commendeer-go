package utilities

import (
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Setup initiates a random seed, or, providing -1, is randomly selected
func Setup(seed int) {
	if seed > -1 {
		rand.Seed(int64(seed))
		return
	}
	rand.Seed(time.Now().UnixNano())
}

// GenerateCodes creates a 'quantity' (int) of random strings using the constant: letters
// l : length of code
func GenerateCodes(quantity, l int) []string {
	var codes []string
	m := make(map[string]bool)

	for i := 0; i < quantity; i++ {
		unique := false
		c := ""
		for !unique {
			c = generateCode(l)

			_, prs := m[c] // check if c present in map
			if !prs {
				unique = true
			}
		}
		m[c] = true // add unique code to map
	}
	for k := range m {
		codes = append(codes, k)
	}
	return codes
}

func generateCode(l int) string {
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
