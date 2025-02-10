package bloom

import (
	"fmt"
	"hash"
	"math"

	"github.com/cespare/xxhash/v2"
)

type Bloom struct {
	bitmap       []byte
	numHashFuncs int
	numBits      uint64
	hashFunc     hash.Hash64
}

func NewBloomFilter(numElements int, falsePositiveProb float64) *Bloom {
	numBits, numHashFuncs := optimalParameters(numElements, falsePositiveProb)
	numOfBytes := numBits / 8
	fmt.Printf("numBits: %d,  numHashFuncs: %d\n", numBits, numOfBytes)
	return &Bloom{
		bitmap:       make([]byte, numOfBytes),
		numHashFuncs: numHashFuncs,
		numBits:      numBits,
		hashFunc:     xxhash.New(),
	}
}

func (b *Bloom) Contains(key string) (bool, error) {
	b.hashFunc.Reset()
	_, err := b.hashFunc.Write([]byte(key))
	if err != nil {
		return false, fmt.Errorf("failed to create hash: %v", err)
	}
	hash := b.hashFunc.Sum64()
	for i := 0; i < b.numHashFuncs; i++ {
		index, offset := b.getBitPosition(hash, i)
		if b.bitmap[index]&(1<<offset) == 0 {
			return false, nil
		}
	}
	return true, nil
}

func (b *Bloom) Insert(key string) error {
	b.hashFunc.Reset()
	_, err := b.hashFunc.Write([]byte(key))
	if err != nil {
		return fmt.Errorf("failed to create hash: %v", err)
	}
	hash := b.hashFunc.Sum64()
	for i := 0; i < b.numHashFuncs; i++ {
		index, offset := b.getBitPosition(hash, i)
		b.bitmap[index] = b.bitmap[index] | (1 << offset)
	}
	return nil
}

func (b *Bloom) getBitPosition(hash uint64, index int) (byteIndex, bitOffset uint64) {
	// Fash hashing
	h1 := hash
	h2 := hash >> 32
	hash = h1 + uint64(index)*h2
	position := hash % uint64(b.numBits/8)
	return position, position % 8
}

func optimalParameters(numElements int, falsePositiveProb float64) (numBits uint64, numHashFuncs int) {
	// m = -(n * ln(p)) / (ln(2))^2
	numBits = uint64(math.Ceil(-float64(numElements) * math.Log(falsePositiveProb) / math.Pow(math.Log(2), 2)))

	// k = (m/n) * ln(2)
	numHashFuncs = int(math.Ceil(float64(numBits) / float64(numElements) * math.Log(2)))

	return numBits, numHashFuncs
}
