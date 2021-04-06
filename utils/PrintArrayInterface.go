/*
 * Author: Yang Aobo
 * Telegram: @AnAsianGangster
 * Created At: March 17, 2021
 * Updated At: March 17, 2021
 * Last Modified By: Yang Aobo
 */

/**
 * This package contains utility functions
 *
 *
 * This file exposes a **PrintArratInterface** function that prints all items
 * in a **[]interface{}** in the console
 */
package utils

import "fmt"

func PrintArrayInterface(items []interface{}) {
	for i := 0; i < len(items); i++ {
		fmt.Println(items[i])
	}
}
