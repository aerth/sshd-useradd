package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	log.SetPrefix("")
	log.SetFlags(log.Ltime)
}

var addr = "127.0.0.1:8666"
var handler = func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		log.Println(r.Form)
	}
	log.Println(r.Method, r.UserAgent(), r.URL.Path)
	w.Write([]byte(`{}`))
}

func main() {
	http.ListenAndServe(addr, http.HandlerFunc(handler))
}
