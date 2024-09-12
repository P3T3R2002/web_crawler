package main

import (
	"fmt"
)
func main() {
	arg := getArg()

	/*var visited []string
	
	err := crawler(arg, arg, &visited)
	if err != nil {
		fmt.Println(err)
	}*/
	fmt.Println("//-------------------------------------------")
	erro := start_crawler(arg)
	if erro != nil {
		fmt.Println(erro)
	}
}