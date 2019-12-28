package internal

import "math/rand"

// Cas stamp number
func Cas() int64 {
	return rand.Int63()
}
