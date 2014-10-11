package main

import (
	"log"
	"net/http"

	"./proxy"
	"github.com/gorilla/mux"
)

// TODO make cli interface with useful help and file handle
// TODO ensure/test response is application/json
// TODO generate godoc
// TODO write a nice README.md
// TODO logging

func main() {

	p, err := proxy.NewProxy("./users.json")
	if err != nil {
		log.Fatalf("Proxy init failed: %s", err)
	}

	r := mux.NewRouter()

	// define the api method expectations, useful for future api handles
	api := r.
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded").
		Subrouter()

	api.Handle("/api/2/domains/{domain}/proxyauth/", p.Authenticate())
	http.ListenAndServe(":8080", r) // TODO make cli argument
}
