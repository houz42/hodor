package server

import (
	"net/http"
)

// middleware - a type that describes a middleware, at the core of this
// implementation a middleware is merely a function that takes a handler
// function, and returns a handler function.
type middleware func(next http.HandlerFunc) http.HandlerFunc

// chain middlewares
func chain(f http.HandlerFunc, mws ...middleware) http.HandlerFunc {
	if f == nil {
		f = func(http.ResponseWriter, *http.Request) {}
	}
	if len(mws) == 0 {
		return f
	}
	return mws[0](chain(f, mws[1:]...))
}
