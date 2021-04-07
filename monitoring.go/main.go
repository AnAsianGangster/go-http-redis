package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
)

func getHashValue(key string) uint64 {
	// Calculate hash value
	hash := md5.Sum([]byte(key))
	fmt.Println(hash)
	value := binary.LittleEndian.Uint64(hash[:8])

	return value
}

func main() {

	// Listen at some port

	number_of_nodes := 3
	node_value_mapping := make([]uint64, number_of_nodes)
	MAX_VAL := 9223372036854775807

	for i := 0; i < number_of_nodes; i++ {
		// 64* bit / number of nodes * (node id)
		node_value_mapping[i] = uint64(MAX_VAL / number_of_nodes * i)
	}

	// When the client ask the node
	hash := getHashValue("keys14d")
	var idx int
	for i := number_of_nodes - 1; i >= 0; i-- {
		// 64* bit / number of nodes * (node id)
		if hash > node_value_mapping[i] {
			idx = i
			break
		}
	}

	fmt.Println("this is the main.go for go-consistent-hashing", idx)

	// println(getHashValue(key))

	// For demo purpose

}
