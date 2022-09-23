package handler

import (
	"fmt"
	"net/http"

	"github.com/kaiya/play/httprouter"
)

var router *httprouter.Router

func init() {
	router = httprouter.New()
	registerHttpRouter(router)
}

func registerHttpRouter(r *httprouter.Router) {
	r.NotFound = http.HandlerFunc(NotFoundHandler)
	r.GET("/hello/:name", HelloHandler)
}

func RouterHandler(w http.ResponseWriter, r *http.Request) {
	// preprocess
	router.ServeHTTP(w, r)
	// postprocess
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("opps, resource not found"))
}

func HelloHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("hello, %s", p.ByName("name"))))
}
