package main

import (
	"fmt"
	"sync"
)

func main(){
	c1:= gen(1,2,3)
	c2:= gen(4,5,6)
	c:= merge(c1, c2)
	for m:= range c {
		fmt.Println(m)
	}

}

func gen(nums ...int) <- chan int {
	result := make(chan int)
	go func() {
		for _, num := range nums {
			result <- num
		}
		close(result)
	}()
	return result
}

func sq(in <- chan int) <-chan int {
	r:=make(chan int)
	go func() {
		for i := range in {
			r <- i * i
		}
		close(r)
	}()
	return r
}

func merge(ins ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	r:= make(chan int)
	wg.Add(len(ins))
	for _, in := range ins {
		go func(c <-chan int) {
			for n:= range c {
				r <- n
			}
			wg.Done()
		}(in)
	}
	go func() {
		wg.Wait()
		close(r)
	}()
	return r
}