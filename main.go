package main

import (
	"fmt"
)
func main() {
	base, concurrency, max_visit, err := getArg()
	if err != nil {
		fmt.Println(err)
		return 
	}

	/*var visited []string
	
	err := crawler(arg, arg, &visited)
	if err != nil {
		fmt.Println(err)
		return
	}*/
	fmt.Println("//-------------------------------------------")
	erro := start_crawler(base, concurrency, max_visit)
	if erro != nil {
		fmt.Println(erro)
		return 
	}
}