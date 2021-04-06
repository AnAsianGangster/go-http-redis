package main

import (
	// "bufio"
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ClientQuery struct {
	command string
	key     string
	value   string
}

// receive input from client
func main() {

	// initialize buffer
	buf := bufio.NewReader(os.Stdin)
	fmt.Println("Client is successfully initialized")

	for {
		select {
		default:
			// get key and value from command line
			raw_input, err := buf.ReadBytes('\n')
			if err != nil {
				fmt.Println(err)
				continue
			}
			// slice the new line character and convert it to string
			input := string(raw_input[:len(raw_input)-1])
			input_array := strings.Split(input, " ")

			if len(input_array) < 2 {
				fmt.Println("Invalid Input!")
			}

			cq := ClientQuery{
				command: input_array[0],
				key:     input_array[1],
				value:   strings.Join(input_array[2:], " "),
			}

			/*
				Dynamo treats both the key and the object supplied by the caller  as an opaque array of bytes. It applies a MD5 hash on the key to
				generate a 128-bit identifier, which is used to determine the storage nodes that are responsible for serving the key
			*/

			// Calculate hash value
			hash := md5.Sum([]byte(cq.key))
			i := binary.LittleEndian.Uint64(hash[:])
			fmt.Printf("Looking for Nodes that handles: %v\n", i)

			// port_number := "5000"

			// Calculate which node handles

			fmt.Printf("Send request to server ----- %v Key: %v with Value: %v\n", cq.command, cq.key, cq.value)

			// send http request here and wait for reply

			if strings.Contains(cq.command, "store") {
				// query := string(`{"node":"server:1","key":"` + cq.key + `","value":"` + cq.value + `"}`)
				query := fmt.Sprintf("{\"node\":\"server:1\",\"key\":\"%v\",\"value\":\"%d\"\"}")
				// fmt.Println(query)
				byte_query := []byte(query)

				req, err := http.NewRequest("POST", "http://localhost:5000/key-value-pair", bytes.NewBuffer(byte_query))
				req.Header.Set("X-Custom-Header", "myvalue")
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()
				// defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				fmt.Println(string(body))
			}

			// if strings.Contains(cq.command, "fetch") {

			// /* 				fmt.Println(cq.key)
			// 				query := "http://localhost:5000/key-value-pair?node=server:1&key=" + cq.key

			// 				resp, err := http.Get(query)
			// 				if err != nil {
			// 					// handle error
			// 					fmt.Println(err)
			// 				}
			// 				fmt.Println(resp)
			// defer resp.Body.Close()
			// body, err := ioutil.ReadAll(resp.Body)
			// fmt.Println(string(body))
			// }

		}
	}

}
