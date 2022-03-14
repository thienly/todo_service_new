//package main
//
//import (
//	"encoding/json"
//	"encoding/pem"
//	"errors"
//	"fmt"
//	"io"
//	"io/fs"
//	"io/ioutil"
//	"runtime/debug"
//	"sync"
//
//	"github.com/golang/protobuf/proto"
//)
//
//type Person struct {
//	FirstName string
//	LastName  string
//}
//
//func main() {
//	p := &Person{FirstName: "Thien", LastName: "Ly"}
//	bytes, err := json.Marshal(p)
//	if err != nil {
//		panic("Application panic")
//	}
//	ioutil.WriteFile("test1.text", bytes, fs.ModeDir)
//
//}
//
//// ParseRSAPublicKeyFromPEM parses a PEM encoded PKCS1 or PKCS8 public key
//func ParseRSAPublicKeyFromPEM(key []byte) error {
//
//	// Parse PEM block
//	var block *pem.Block
//	if block, _ = pem.Decode(key); block == nil {
//		return fmt.Errorf("%v", "error")
//	}
//
//	return errors.New("No error")
//}
//
//type request func(wg *sync.WaitGroup, done chan interface{}, results chan interface{}, in interface{})
//
//func ReplicateRequest(re request, in interface{}) interface{} {
//	var wg sync.WaitGroup
//	done := make(chan interface{})
//	results := make(chan interface{})
//	wg.Add(2)
//	for i := 0; i < 2; i++ {
//		go re(&wg, done, results, in)
//	}
//	// Waiting for first response
//	data := <-results
//	close(done)
//	wg.Wait()
//	return data
//}
//
//// before sending the message to component, we need cho check it is runing to make sure it is working.
//
//type MyError struct {
//	Inner   error
//	Message string
//	Stack   string
//	Misc    map[string]interface{}
//}
//
//func WrapErrr(err error, msgf string, msgArgs ...interface{}) *MyError {
//	return &MyError{
//		Inner:   err,
//		Message: fmt.Sprintf(msgf, msgArgs...),
//		Stack:   string(debug.Stack()),
//		Misc:    make(map[string]interface{}),
//	}
//}
//
//func generator(done <-chan interface{}, nums ...int) chan interface{} {
//	c := make(chan interface{})
//	go func() {
//		defer close(c)
//		for _, v := range nums {
//			select {
//			case <-done:
//			case c <- v:
//			}
//		}
//	}()
//	return c
//}
//
//func fanIn(done chan interface{}, channels ...chan interface{}) chan interface{} {
//	c := make(chan interface{})
//	var wg sync.WaitGroup
//	wg.Add(len(channels))
//	for i := 0; i < len(channels); i++ {
//		go func(in int) {
//			defer wg.Done()
//			for i := range channels[in] {
//				select {
//				case <-done:
//					return
//				case c <- i:
//				}
//			}
//		}(i)
//	}
//	go func() {
//		wg.Wait()
//		close(c)
//	}()
//	return c
//}
//
//func fanOut(done chan interface{}, in chan interface{}, num int) []chan interface{} {
//	sl := make([]chan interface{}, num)
//
//	return sl
//}
//
//type Options struct {
//	firstName string
//	lastName  string
//}
//type Option func(*Options)
//
//func withFirstName(s string) Option {
//	return func(o *Options) {
//		o.firstName = s
//	}
//}
//func withLastName(l string) Option {
//	return func(o *Options) {
//		o.lastName = l
//	}
//}
//func NewOptions(opts ...Option) *Options {
//	value := &Options{}
//	for _, v := range opts {
//		v(value)
//	}
//	return value
//}
//
//type Codec struct {
//	conn io.ReadWriteCloser
//}
//
//func NewCode(conn io.ReadWriteCloser) *Codec {
//	return &Codec{conn: conn}
//}
//func (c *Codec) ReadBody(b interface{}) error {
//	if b == nil {
//		return nil
//	}
//	buf, err := io.ReadAll(c.conn)
//	if err != nil {
//		return err
//	}
//	m, ok := b.(proto.Message)
//	if !ok {
//		return errors.New("error")
//	}
//	return proto.Unmarshal(buf, m)
//}
//
//// func (c *Codec) Write(m *codec.Message, b interface{}) error {
//// 	if b == nil {
//// 		// Nothing to write
//// 		return nil
//// 	}
//// 	p, ok := b.(proto.Message)
//// 	if !ok {
//// 		return errors.New("error")
//// 	}
//// 	buf, err := proto.Marshal(p)
//// 	if err != nil {
//// 		return err
//// 	}
//// 	_, err = c.conn.Write(buf)
//// 	return err
//// }
