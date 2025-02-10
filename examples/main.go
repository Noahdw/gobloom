package main

import (
	"fmt"

	bloom "github.com/noahdw/gobloom"
)

func main() {

	var t test
	t.Start()
}

type test struct {
}

func (t *test) Start() {
	t.testFalsePositive(10000000, 100000, 10000000)
	t.testFalsePositive(1000, 50, 157)
	t.testFalsePositive(1000, 100, 1000)
	t.testFalsePositive(10000000, 100000, 1000000)
}

func (t *test) testFalsePositive(numSlots, numKeys, numTrys int) {
	filter := bloom.NewBloomFilter(numSlots, 0.1)

	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("%d", i)
		filter.Insert(key)
	}

	count := 0
	for i := 0; i < numTrys; i++ {
		key := fmt.Sprintf("%d", i)
		if ok, _ := filter.Contains(key); ok {
			if i > numKeys {
				count++
			}
		}
	}
	fmt.Printf("false positive count: %d over numTrys: %d\n", count, numTrys-numKeys)
	fmt.Printf("false positive percent: %f\n", float64(count)/float64(numTrys-numKeys))
}
