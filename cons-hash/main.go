package main

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type message struct {
	src          int
	dst          int
	message_type string
	key          string
	value        string
	data         value
}

type value struct {
	value        string
	vector_clock int
	coordinator  int
	confirm      bool
}

var EMPTY_VALUE value = value{"", -1, -1, false}

var NUMBER_OF_NODES int = 10
var NUMBER_OF_REPLICAS int = 7
var channel_array []chan message = []chan message{}
var wg sync.WaitGroup

// Calculate hash value
func getHashValue(key string) uint64 {
	hash := md5.Sum([]byte(key))
	// fmt.Println(hash)
	value := binary.LittleEndian.Uint64(hash[0:8])

	return value
}

// Find Coordinator node responsible for the value
func findNode(hash uint64, node_value_mapping []uint64) int {
	var idx int
	for i := NUMBER_OF_NODES - 1; i >= 0; i-- {
		// 64* bit / number of nodes * (node id)
		if hash > node_value_mapping[i] {
			idx = i
			break
		}
	}
	return idx
}

func distributeValue(parent int, number_of_nodes int) []int {
	number_of_children := number_of_nodes / 2
	children_array := make([]int, number_of_children) // children nodes are the next n / 2 --- able to handle n / 2 - 1 node failures
	for i := 0; i < (number_of_children); i++ {
		children_array[i] = (parent + i) % 12
	}

	return children_array
}

// Formatted print of the hash map
func printMap(m map[string]value, node_id int) {
	var maxLenKey int
	for k, _ := range m {
		if len(k) > maxLenKey {
			maxLenKey = len(k)
		}
	}

	var str string = ""
	for k, v := range m {
		str = fmt.Sprintf(str + k + ": " + strings.Repeat(" ", maxLenKey-len(k)) + v.value + " --- vectorclock: " + strconv.Itoa(v.vector_clock) + "\n")
		// str = str + "" + k + ": " + strings.Repeat(" ", maxLenKey-len(k)) + v.value, v.vector_clock
	}
	fmt.Printf("NODE %d \t======>\n%v\n", node_id, str)
}

func node_write_consensus(ch chan message, message_rcv message, node_id int, latest_value value, m map[string]value) map[string]value {

	m[message_rcv.key] = latest_value
	for i := 0; i < NUMBER_OF_REPLICAS-1; i++ {
		idx := (latest_value.coordinator + i + 1) % NUMBER_OF_NODES
		fmt.Printf("NODE %d \t======> Sending check write request to Node %d\n", node_id, idx)
		channel_array[idx] <- message{node_id, idx, "check_write", message_rcv.key, message_rcv.value, latest_value}
	}

	// Timeout if it has not receive (N / 2) / 2 replies => minimum number of reply needed for successful GET request
	timeout := time.NewTimer(1 * time.Second)
	reply_count := 0
	exit := false
	for {
		select {
		case check_message := <-ch:
			// getting check_read_reply message
			if check_message.message_type == "check_write_reply" {
				// report to client if it has received (N / 2) / 2 reply
				reply_count = reply_count + 1
				if reply_count == (NUMBER_OF_REPLICAS / 2) {
					fmt.Printf("NODE %d \t======> Done getting consensus, received %d replies (including self) out of %d and proceed with returning write request response\n", node_id, reply_count+1, NUMBER_OF_REPLICAS)
					channel_array[message_rcv.src] <- message{node_id, message_rcv.src, "success", message_rcv.key, latest_value.value, EMPTY_VALUE}
					exit = true
				}
			}
		case <-timeout.C:
			// If timeout tell client that GET request has failed
			fmt.Printf("NODE %d \t======> Node timeout, failed getting consensus\n", node_id)
			channel_array[message_rcv.src] <- message{node_id, message_rcv.src, "fail", message_rcv.key, "", EMPTY_VALUE}
			exit = true
		}
		if exit {
			break
		}
	}
	return m
}

// Node go routine
func node(node_id int, ch chan message) {
	defer wg.Done()
	m := make(map[string]value)
	dead := false

	for {
		select {
		case message_rcv := <-ch:
			// Handle Client GET key request
			if message_rcv.message_type == "get" && dead != true {
				// If key exists, retrieve the value of the key else infrom the client

				// Ask the nodes that have the replicas
				for i := 0; i < NUMBER_OF_REPLICAS-1; i++ {
					idx := (message_rcv.data.coordinator + i + 1) % NUMBER_OF_NODES
					fmt.Printf("NODE %d \t======> Sending check read request to Node %d\n", node_id, idx)
					channel_array[idx] <- message{node_id, idx, "check_read", message_rcv.key, "", EMPTY_VALUE}
				}
				// Get the latest replica from local map
				var latest_value value
				_, exists := m[message_rcv.key]
				if exists {
					latest_value = m[message_rcv.key]
				} else {
					latest_value = EMPTY_VALUE
				}
				// Timeout if it has not receive (N / 2) / 2 replies => minimum number of reply needed for successful GET request
				timeout := time.NewTimer(1 * time.Second)
				reply_count := 0
				exit := false
				for {
					select {
					case check_message := <-ch:
						// getting check_read_reply message
						if check_message.message_type == "check_read_reply" {
							// update latest_value if it receives latest value
							if check_message.data.vector_clock > latest_value.vector_clock {
								latest_value = check_message.data
							}
							// report to client if it has received (N / 2) / 2 reply
							reply_count = reply_count + 1
							if reply_count == (NUMBER_OF_REPLICAS / 2) {
								fmt.Printf("NODE %d \t======> Done getting consensus, received %d replies (including self) out of %d and proceed with returning read request response\n", node_id, reply_count+1, NUMBER_OF_REPLICAS)
								channel_array[message_rcv.src] <- message{node_id, message_rcv.src, "success", message_rcv.key, latest_value.value, EMPTY_VALUE}
								// update current entry
								m[message_rcv.key] = latest_value
								// inform all other replicas
								exit = true
							}
						}
					case <-timeout.C:
						// If timeout tell client that GET request has failed
						fmt.Printf("NODE %d \t======> Node timeout, failed getting consensus\n", node_id)
						channel_array[message_rcv.src] <- message{node_id, message_rcv.src, "fail", message_rcv.key, "", EMPTY_VALUE}
						exit = true
					}
					if exit {
						break
					}
				}
			}
			// handle internal consensus for GET request
			if message_rcv.message_type == "check_read" && dead != true {
				var latest_value value
				_, exists := m[message_rcv.key]
				if exists {
					latest_value = m[message_rcv.key]
				} else {
					latest_value = EMPTY_VALUE
				}
				fmt.Printf("NODE %d \t======> Sending check read reply to Node %d\n", node_id, message_rcv.src)
				channel_array[message_rcv.src] <- message{node_id, 0, "check_read_reply", message_rcv.key, latest_value.value, latest_value}
			}
			if message_rcv.message_type == "post" && dead != true {
				// If key exists, add a new key with some value else infrom the client
				_, exists := m[message_rcv.key]

				if exists {
					if m[message_rcv.key].value != "" {
						fmt.Printf("NODE %d \t======> Key exist, cannot be rewritten, try put instead", node_id)
						continue
					}
				}
				m = node_write_consensus(ch, message_rcv, node_id, value{message_rcv.value, 1, node_id, false}, m)
			}
			// Handle Client PUT request
			if message_rcv.message_type == "put" && dead != true {
				// If key exists, edit the value of the key else infrom the client
				_, exists := m[message_rcv.key]
				if exists {
					curr_timestamp := m[message_rcv.key].vector_clock
					m = node_write_consensus(ch, message_rcv, node_id, value{message_rcv.value, curr_timestamp + 1, m[message_rcv.key].coordinator, false}, m)
				} else {
					fmt.Printf("NODE %d \t======> Cannot execute delete command, Key does not exist", node_id)
				}
			}
			// Handle Client DELETE request
			if message_rcv.message_type == "delete" && dead != true {
				// If key exists, delete the key value pair else infrom the client
				_, exists := m[message_rcv.key]

				if exists {
					curr_timestamp := m[message_rcv.key].vector_clock
					m = node_write_consensus(ch, message_rcv, node_id, value{"", curr_timestamp + 1, m[message_rcv.key].coordinator, false}, m)
				} else {
					fmt.Printf("NODE %d \t======> Cannot execute delete command, Key does not exist", node_id)
				}
			}
			// handle internal consensus for GET request
			if message_rcv.message_type == "check_write" && dead != true {
				m[message_rcv.key] = message_rcv.data
				fmt.Printf("NODE %d \t======> Sending check read reply to Node %d\n", node_id, message_rcv.src)
				channel_array[message_rcv.src] <- message{node_id, message_rcv.src, "check_write_reply", message_rcv.key, message_rcv.value, message_rcv.data}
			}
			if (message_rcv.message_type == "check" || message_rcv.message_type == "checkall") && dead != true {
				// If key exists, delete the key value pair else infrom the client
				// fmt.Printf("Node %d hashmap\n", node_id)
				if len(m) > 0 {
					printMap(m, node_id)
				} else {
					fmt.Printf("NODE %d \t======>\nNo keys\n", node_id)
				}
			}
			if message_rcv.message_type == "revive" {
				fmt.Printf("NODE %d \t======> REVIVED\n", node_id)
				dead = false
			}
			if dead {
				fmt.Printf("NODE %d \t======> DEAD\n", node_id)
			}
			if message_rcv.message_type == "kill" {
				fmt.Printf("NODE %d \t======> KILLED\n", node_id)
				dead = true
			}
		}
	}
}

func main() {

	if NUMBER_OF_REPLICAS > NUMBER_OF_NODES {
		fmt.Println("Number of Replicas more than number of Nodes")
		return
	}
	// Listen at some port
	node_value_mapping := make([]uint64, NUMBER_OF_NODES)
	MAX_VAL := ^uint64(0)

	// 64* bit / number of nodes * (node id)
	for i := 0; i < NUMBER_OF_NODES; i++ {
		node_value_mapping[i] = uint64(MAX_VAL / uint64(NUMBER_OF_NODES) * uint64(i))
	}

	var client_channel chan message = make(chan message, NUMBER_OF_NODES*10) // node communicate to client via this channel

	for i := 0; i < NUMBER_OF_NODES; i++ {
		// Create communicating channel for each node and store in array so other routine can refer
		var node_channel chan message = make(chan message, NUMBER_OF_NODES*10)
		channel_array = append(channel_array, node_channel)
		go node(i, node_channel)
		// wg.Add(1)
	}
	channel_array = append(channel_array, client_channel)

	// initialize buffer
	var idx int
	var hash uint64

	buf := bufio.NewReader(os.Stdin)
	fmt.Printf("========================================\nThis is Client Command Line, please use the following commands with the correct number of arguments: \nget	--- get key\npost	--- post key value\nput	--- put key value\ndelete	--- delete key\ncheck	--- check\nkill	--- kill node_id\nrevive	--- revive node_id\n")
	for {
		select {
		default:
			time.Sleep(time.Duration(100 * time.Millisecond))
			fmt.Printf("> ")
			// get key and value from command line
			raw_input, err := buf.ReadBytes('\n')
			if err != nil {
				fmt.Println(err)
				continue
			}

			// slice the new line character and convert it to string
			fmt.Println("LOG:")
			input := string(raw_input[:len(raw_input)-1])
			input_array := strings.Split(input, " ")

			// fix slicing when input doesn't have white space
			input_array[len(input_array)-1] = strings.TrimSpace(input_array[len(input_array)-1])
			// if len(input_array) < 2 {
			// 	input_array[0] = input[:len(input)-1]
			// }

			if input_array[0] == "get" {
				// get key
				if len(input_array) < 2 {
					fmt.Println("CLIENT \t======> Invalid format, need at least 1 argument")
					continue
				}
				key := input_array[1]

				// fmt.Printf("Client getting the key %v\n", key)
				hash = getHashValue(key)
				idx = findNode(hash, node_value_mapping)

				for i := 0; i < (NUMBER_OF_NODES/2)+1; i++ {
					curr_idx := (idx + i) % NUMBER_OF_NODES
					replied := false
					fmt.Printf("CLIENT\t======> Finding the hash value of key: %v => hash value: %d => coordinator for this key: Node %v => asking node: %v\n", key, hash, idx, curr_idx)
					channel_array[curr_idx] <- message{NUMBER_OF_NODES, curr_idx, "get", key, "", value{"", -1, idx, false}}

					timeout := time.NewTimer(2 * time.Second)
					select {
					case message_rcv := <-client_channel:
						if message_rcv.message_type == "success" {
							if message_rcv.value != "" {
								fmt.Printf("CLIENT \t======> Successfully retrieved the value of the key. key: %v => value: %v\n", key, message_rcv.value)
							} else {
								fmt.Printf("CLIENT \t======> Successfully retrieved the value of the key. key: %v has no value\n", key)
							}
						} else {
							fmt.Printf("CLIENT \t======> Failed to retrieve the value of the key.\n")
						}
						replied = true
					case <-timeout.C:
						fmt.Println("CLIENT \t======> Timeout")
					}
					if replied {
						break
					}
				}

				continue
			}
			if input_array[0] == "post" {
				// post key value
				if len(input_array) < 3 {
					fmt.Println("CLIENT \t======> Invalid format, need at least 2 argument")
					continue
				}

				key := input_array[1]
				value := input_array[2]

				hash = getHashValue(key)
				idx = findNode(hash, node_value_mapping)

				for i := 0; i < (NUMBER_OF_NODES/2)+1; i++ {
					curr_idx := (idx + i) % NUMBER_OF_NODES
					replied := false
					fmt.Printf("CLIENT\t======> Finding the hash value of key: %v => hash value: %d => coordinator for this key: Node %v => asking node: %v\n", key, hash, idx, curr_idx)
					channel_array[curr_idx] <- message{NUMBER_OF_NODES, curr_idx, "post", key, value, EMPTY_VALUE}
					timeout := time.NewTimer(2 * time.Second)
					select {
					case message_rcv := <-client_channel:
						if message_rcv.message_type == "success" {
							fmt.Printf("CLIENT \t======> Successfully post the value of the key. key: %v => value: %v\n", key, value)
						} else {
							fmt.Printf("CLIENT \t======> Failed to post the value of the key.\n")
						}
						replied = true
					case <-timeout.C:
						fmt.Println("CLIENT \t======> Timeout")
					}
					if replied {
						break
					}
				}
				continue
			}
			if input_array[0] == "put" {
				// put key value
				if len(input_array) < 3 {
					fmt.Println("CLIENT \t======> Invalid format, need at least 2 argument")
					continue
				}

				key := input_array[1]
				value := input_array[2]

				hash = getHashValue(key)
				idx = findNode(hash, node_value_mapping)

				for i := 0; i < (NUMBER_OF_NODES/2)+1; i++ {
					curr_idx := (idx + i) % NUMBER_OF_NODES
					replied := false
					fmt.Printf("CLIENT\t======> Finding the hash value of key: %v => hash value: %d => coordinator for this key: Node %v => asking node: %v\n", key, hash, idx, curr_idx)
					channel_array[curr_idx] <- message{NUMBER_OF_NODES, curr_idx, "put", key, value, EMPTY_VALUE}
					timeout := time.NewTimer(2 * time.Second)
					select {
					case message_rcv := <-client_channel:
						if message_rcv.message_type == "success" {
							fmt.Printf("CLIENT \t======> Successfully put the value of the key. key: %v => value: %v\n", key, value)
						} else {
							fmt.Printf("CLIENT \t======> Failed to put the value of the key.\n")
						}
						replied = true
					case <-timeout.C:
						fmt.Println("CLIENT \t======> Timeout")
					}
					if replied {
						break
					}
				}
				continue
			}
			if input_array[0] == "delete" {
				// delete key value
				if len(input_array) < 2 {
					fmt.Println("CLIENT \t======> Invalid format, need at least 1 argument")
					continue
				}

				key := input_array[1]
				hash = getHashValue(key)
				idx = findNode(hash, node_value_mapping)

				for i := 0; i < (NUMBER_OF_NODES/2)+1; i++ {
					curr_idx := (idx + i) % NUMBER_OF_NODES
					replied := false
					fmt.Printf("CLIENT\t======> Finding the hash value of key: %v => hash value: %d => coordinator for this key: Node %v => asking node: %v\n", key, hash, idx, curr_idx)
					channel_array[curr_idx] <- message{NUMBER_OF_NODES, curr_idx, "delete", key, "", EMPTY_VALUE}
					timeout := time.NewTimer(2 * time.Second)
					select {
					case message_rcv := <-client_channel:
						if message_rcv.message_type == "success" {
							fmt.Printf("CLIENT \t======> Successfully delete the value of the key. key: %v\n", key)
						} else {
							fmt.Printf("CLIENT \t======> Failed to delete the value of the key.\n")
						}
						replied = true
					case <-timeout.C:
						fmt.Println("CLIENT \t======> Timeout")
					}
					if replied {
						break
					}
				}
				continue
			}
			if input_array[0] == "kill" {
				// kill node_id
				if len(input_array) < 2 {
					fmt.Println("CLIENT \t======> Invalid format, need at least 1 argument")
					continue
				}
				node_id, err := strconv.Atoi(strings.TrimSpace(input_array[1]))
				if err != nil {
					fmt.Println("CLIENT \t======> Invalid node")
				}

				if node_id > NUMBER_OF_NODES && node_id < 0 {
					fmt.Println("CLIENT \t======> Node does not exist try int between 0 and ", NUMBER_OF_NODES)
				}

				channel_array[node_id] <- message{NUMBER_OF_NODES, node_id, "kill", "", "", EMPTY_VALUE}
				fmt.Printf("CLIENT \t======> Client killed Node %d\n", node_id)
				continue
			}
			if input_array[0] == "revive" {
				// revive node_id
				if len(input_array) < 2 {
					fmt.Println("CLIENT \t======> Invalid format, need at least 1 argument")
					continue
				}
				node_id, err := strconv.Atoi(strings.TrimSpace(input_array[1]))
				if err != nil {
					fmt.Println("CLIENT \t======> Invalid node")
					continue
				}

				if node_id > NUMBER_OF_NODES && node_id < 0 {
					fmt.Println("CLIENT \t======> Node does not exist try int between 0 and ", NUMBER_OF_NODES)
					continue
				}

				channel_array[node_id] <- message{NUMBER_OF_NODES, node_id, "revive", "", "", EMPTY_VALUE}
				fmt.Printf("CLIENT \t======> Client revived Node %d\n", node_id)
				continue
			}
			if input_array[0] == "checkall" {
				for i := 0; i < NUMBER_OF_NODES; i++ {
					channel_array[i] <- message{NUMBER_OF_NODES, i, "check", "", "", EMPTY_VALUE}
				}
				continue
			}
			if input_array[0] == "check" {
				if len(input_array) < 2 {
					fmt.Println("CLIENT \t======> Invalid format, need at least 1 argument")
					continue
				}
				node_id, err := strconv.Atoi(strings.TrimSpace(input_array[1]))
				if err != nil {
					fmt.Println("CLIENT \t======> Invalid node")
					continue
				}

				if node_id > NUMBER_OF_NODES && node_id < 0 {
					fmt.Println("CLIENT \t======> Node does not exist try int between 0 and ", NUMBER_OF_NODES)
					continue
				}
				channel_array[node_id] <- message{NUMBER_OF_NODES, node_id, "check", "", "", EMPTY_VALUE}
				continue
			}
			fmt.Println("CLIENT \t======> Wrong command")
		}
	}
}
