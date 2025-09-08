package dymean

import (
	"hash"
	"hash/fnv"
)

// BloomFilter represents a probabilistic data structure for membership testing
type BloomFilter struct {
	bitArray     []bool
	size         uint
	hashFuncs    []hash.Hash64
	numHashFuncs int
}

// NewBloomFilter creates a new Bloom filter with the specified size and number of hash functions
func NewBloomFilter(size uint, numHashFuncs int) *BloomFilter {
	bf := &BloomFilter{
		bitArray:     make([]bool, size),
		size:         size,
		numHashFuncs: numHashFuncs,
		hashFuncs:    make([]hash.Hash64, numHashFuncs),
	}

	// Initialize hash functions
	for i := 0; i < numHashFuncs; i++ {
		bf.hashFuncs[i] = fnv.New64a()
	}

	return bf
}

// Add adds an item to the Bloom filter
func (bf *BloomFilter) Add(item string) {
	for i := 0; i < bf.numHashFuncs; i++ {
		bf.hashFuncs[i].Reset()
		bf.hashFuncs[i].Write([]byte(item))
		// Add salt to create different hash functions
		bf.hashFuncs[i].Write([]byte{byte(i)})
		hash := bf.hashFuncs[i].Sum64()
		index := hash % uint64(bf.size)
		bf.bitArray[index] = true
	}
}

// Contains checks if an item might be in the Bloom filter
// Returns true if the item is possibly in the set, false if definitely not
func (bf *BloomFilter) Contains(item string) bool {
	for i := 0; i < bf.numHashFuncs; i++ {
		bf.hashFuncs[i].Reset()
		bf.hashFuncs[i].Write([]byte(item))
		bf.hashFuncs[i].Write([]byte{byte(i)})
		hash := bf.hashFuncs[i].Sum64()
		index := hash % uint64(bf.size)
		if !bf.bitArray[index] {
			return false
		}
	}
	return true
}

// AddWords adds multiple words to the Bloom filter
func (bf *BloomFilter) AddWords(words []string) {
	for _, word := range words {
		bf.Add(word)
	}
}
