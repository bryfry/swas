package main

import (
	"log"
	"net/http"

	"./proxy"
	"github.com/gorilla/mux"
)

func main() {

	proxy, err := proxy.NewProxy("./users.json")
	if err != nil {
		log.Fatalln("Proxy init failed: %s", err)
	}

	r := mux.NewRouter()

	// define the api method expectations, useful for future api handles
	api := r.
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded").
		Subrouter()

	api.Handle("/api/2/domains/{domain}/proxyauth/", proxy.Authenticate())
	http.ListenAndServe(":8080", r)
}
