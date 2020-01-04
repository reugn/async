package internal

import "math/rand"

// Cas returns compare-and-set stamp value
func Cas() int64 {
	return rand.Int63()
}
