package main

import "fmt"

type res struct {

}
type req struct {

}
type handlerFunc func(res2 res, req2 req)
type middleware func(next handlerFunc) handlerFunc

func chainMiddleware(mw ...middleware) middleware{
	return func(final handlerFunc) handlerFunc {
		return func(w res, r req) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

func main(){
	// x:= tracing(logging(x,y)) execute: x();
    var tracing middleware = func(next handlerFunc) handlerFunc{
		return func(res2 res, req2 req) {
			fmt.Println("tracing")
			if next != nil {
				next(res2, req2)
			}
		}
	}

	var logging middleware = func(next handlerFunc) handlerFunc{
		return func(res2 res, req2 req) {
			fmt.Println("logging")
			if next != nil {
				next(res2, req2)
			}
		}
	}


	logging(tracing(nil))(res{},req{})
	chainMiddleware(logging,tracing)(nil)(res{},req{})

}


