package main

import (
	"bufio"
	"fmt"
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

			if len(input_array) != 3 {
				fmt.Println("Invalid Input!")
			}

			clientQuery := ClientQuery{
				command: input_array[0],
				key:     input_array[1],
				value:   input_array[2],
			}

			fmt.Printf("Send request to server ----- %v Key: %v with Value: %v\n", clientQuery.command, clientQuery.key, clientQuery.value)

			// send http request here and wait for reply

			// resp, err := http.PostForm("http://localhost:5000/", url.Values{"key": {clientQuery.key}, "value": {clientQuery.value}})

		}
	}

}
