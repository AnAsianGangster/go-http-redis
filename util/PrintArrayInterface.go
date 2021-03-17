package util

import "fmt"

func PrintArrayInterface(items []interface{}) {
	for i := 0; i < len(items); i++ {
		fmt.Println(items[i])
	}
}
