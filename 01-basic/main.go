package main

import (
	"fmt"

	"basic/slice"
)

func main() {
	result := slice.Delete([]int64{10, 20, 30, 40, 50}, 2)
	fmt.Printf("%v \n\n", result)
}
