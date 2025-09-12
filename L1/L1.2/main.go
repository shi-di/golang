package main

import (
	"fmt"
	"sync"
)

func main() {
	nums := []int{2, 4, 6, 8, 10}

	var wg sync.WaitGroup
	wg.Add(len(nums))

	for _, n := range nums {
		n := n
		go func() {
			res := n * n
			fmt.Println(res)
			wg.Done()
		}()
	}

	wg.Wait()
}
