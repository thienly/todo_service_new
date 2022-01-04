package main1

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sync"
)

type Person struct {
	Name string
}

func (p Person) ChangeName() {
	fmt.Printf("%p\n", &p)
	p.Name = "Changed"
}
type sample struct {
	Ids []int32 `json:"ids"`
	Name string `json:"name"`
}


func main() {

	str :=  `{"ids":[a,2,3], "name" :"ABC"}`
	byt := []byte(str)

	vl := sample{}
	err:= json.Unmarshal(byt, &vl)
	if err !=nil {
		fmt.Println(err)
	}
	fmt.Println(vl.Ids)
	// m, err:= MD5All(os.Args[1])
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// var paths []string
	// for path := range m {
	// 	paths = append(paths, path)
	// }
	// sort.Strings(paths)
	// for _, path := range paths {
	// 	fmt.Printf("%x  %s\n", m[path], path)
	// }
	//p := Person{Name: "Thien"}
	//fmt.Printf("%p\n", &p)
	//p.ChangeName()
	//fmt.Println(p)

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
