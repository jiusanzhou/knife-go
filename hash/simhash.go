package hash

import (
	"github.com/jiusanzhou/knife-go/hash/simhash"
)

var (
	// Simhash returns a 64-bit simhash of the given feature set
	Simhash = simhash.Simhash

	// SimhashCompare calculates the Hamming distance between two 64-bit integers
	SimhashCompare = simhash.Compare
)
