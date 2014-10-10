package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type proxy struct {
	Domains []domain
}

type domain struct {
	Address string `json:"domain"`
	Users   []user `json:"users"`
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func checkErr(err error, m string) {
	if err != nil {
		log.Fatalln(m, ": ", err)
	}
}

func b64sha256(password string) string {
	s256 := sha256.New()
	s256.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(s256.Sum(nil))
}

func (p *proxy) get(domain, username string) (*user, error) {
	for i, d := range p.Domains {
		if d.Address == domain {
			for j, u := range d.Users {
				if u.Username == username {
					return &p.Domains[i].Users[j], nil
				}
			}
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (p *proxy) authenticate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var username, password, domain string
		vars := mux.Vars(r)
		if val, ok := vars["domain"]; ok {
			domain = val
		}
		err := r.ParseForm()
		checkErr(err, "ParseForm: ")
		if val := r.Header.Get("Content-Type"); val == "" {
			// no content-type
		}

		if val := r.Form.Get("username"); val != "" {
			username = val
		} else {
			// no username

		}

		if val := r.Form.Get("password"); val != "" {
			password = val
		} else {
			// no password

		}
		fmt.Println(r.Header.Get("Content-Type"), domain, username, password)
		if _, err := p.get(domain, username); err != nil {
			fmt.Fprintf(w, "{'access_granted':true}")
		}

	})
}

func main() {
	filePath := "./users.json"
	usersJson, err := os.Open(filePath)
	checkErr(err, "Open "+filePath+": ")

	proxy := &proxy{}
	err = json.NewDecoder(usersJson).Decode(&proxy.Domains)
	checkErr(err, "Decode: ")

	for i, d := range proxy.Domains {
		for j, u := range d.Users {
			proxy.Domains[i].Users[j].Password = "{SHA256}" + b64sha256(u.Password)
		}
	}

	r := mux.NewRouter()
	get := r.Methods("POST").Subrouter()
	get.Handle("/api/2/domains/{domain}/proxyauth/", proxy.authenticate())
	http.ListenAndServe(":8080", r)
}
