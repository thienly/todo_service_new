package main1

import (
	"crypto/md5"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sync"
)

func main() {

}

func MD5All(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		m[path] = md5.Sum(data)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}

func sq(in <-chan int) <-chan int {
	r := make(chan int)
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
	r := make(chan int)
	wg.Add(len(ins))
	for _, in := range ins {
		go func(c <-chan int) {
			for n := range c {
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
