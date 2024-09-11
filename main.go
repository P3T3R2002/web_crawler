package main

import (
	"fmt"
)
func main() {
	arg := getArg()

	var visited []string
	
	err := crawler(arg, arg, &visited)
	if err != nil {
		fmt.Println(err)
	}
}