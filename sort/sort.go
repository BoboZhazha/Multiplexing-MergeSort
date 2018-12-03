package main

import (
	"fmt"
	"runtime"
	"sort"
)

func main() {
	a := []int{3, 4, 78, 9, 1, 5, 8, 9, 9, 1}
	sort.Ints(a)

	for _, v := range a {
		fmt.Println(v)
	}
	fmt.Println(runtime.GOMAXPROCS(0))
}
