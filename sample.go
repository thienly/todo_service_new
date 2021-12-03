package main

import "fmt"

// mock dc.
type minioGetter interface {
	PresignedUrl()
}

// lib
type minioClient struct {
}

func (m *minioClient) PresignedUrl() {

}

// user code
type implClient struct {
	minio minioGetter
}

func (i *implClient) hello() {
	i.minio.PresignedUrl() // mock . tra ve
}

func main() {
	//

	// real code
	i := implClient{
		minio: &minioClient{},
	}
	fmt.Println(i)
}
