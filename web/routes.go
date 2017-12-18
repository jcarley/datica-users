package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Get(path string, router *mux.Router, handler http.HandlerFunc) {
	Handle(path, "GET", router, handler)
}

func Post(path string, router *mux.Router, handler http.HandlerFunc) {
	Handle(path, "POST", router, handler)
}

func Put(path string, router *mux.Router, handler http.HandlerFunc) {
	Handle(path, "PUT", router, handler)
}

func Patch(path string, router *mux.Router, handler http.HandlerFunc) {
	Handle(path, "PATCH", router, handler)
}

func Delete(path string, router *mux.Router, handler http.HandlerFunc) {
	Handle(path, "DELETE", router, handler)
}

func Handle(path string, method string, router *mux.Router, handler http.HandlerFunc) {
	router.HandleFunc(path, handler).Methods(method)
}
