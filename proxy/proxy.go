package proxy

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Proxy struct {
	Domains []Domain
}

type Domain struct {
	Address string `json:"domain"`
	Users   []User `json:"users"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Success bool   `json:"access_granted"`
	Reason  string `json:"reason,omitempty"`
}

// generate and return base64 encoded sha256 digest of provided password
func b64sha256(password string) string {
	s256 := sha256.New()
	s256.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(s256.Sum(nil))
}

// Parse the json users file and return the proxy data type
// Unfortunately the process of decoding from json means plaintext
// passwords are handled temporarily and then overwritten with b64sha256
// representations
func NewProxy(filePath string) (*Proxy, error) {
	usersJson, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	p := &Proxy{}
	err = json.NewDecoder(usersJson).Decode(&p.Domains)
	if err != nil {
		return nil, err
	}

	for i, d := range p.Domains {
		for j, u := range d.Users {
			p.Domains[i].Users[j].Password = "{SHA256}" + b64sha256(u.Password)
		}
	}
	return p, nil

}

// search within the proxy for a domain
func (p *Proxy) get(reqDomain string) (*Domain, error) {
	for i, d := range p.Domains {
		if d.Address == reqDomain {
			return &p.Domains[i], nil
		}
	}
	return nil, fmt.Errorf("No such domain")
}

// search within a domain for a user
func (d *Domain) get(reqUser string) (*User, error) {
	for i, u := range d.Users {
		if u.Username == reqUser {
			return &d.Users[i], nil
		}
	}
	return nil, fmt.Errorf("No such user")
}

// Per policy: in case of authentication failure or validation errors.
// The 'reason' is always is always same: "denied by policy".
// Additionally, success is always simply access_granted: true.
// Having this separated as a function will be useful if different
// responses are needed in the future
// Assumption: 200 OK is golang default
func writeSuccess(w http.ResponseWriter, success bool) {
	var r *Response
	out := json.NewEncoder(w)
	if success {
		r = &Response{
			Success: success,
		}
	} else {
		r = &Response{
			Success: success,
			Reason:  "denied by policy",
		}
	}
	out.Encode(r)
}

// used for testing, less performant than using io.Writer interface
func successBody(success bool) string {
	var r *Response
	if success {
		r = &Response{
			Success: success,
		}
	} else {
		r = &Response{
			Success: success,
			Reason:  "denied by policy",
		}
	}
	// live dangerously, ignoring error for this specific use case
	b, _ := json.Marshal(r)
	return string(b) + "\n" // http body fix, expects newline at end
}

// Proxy Authentication handler expects looks up domain (mux url variable),
// username and password (query parameters) and returns the appropriate
// Response and Status Code
func (p *Proxy) Authenticate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// domain lookup
		vars := mux.Vars(r)
		urlDomain, ok := vars["domain"]
		if !ok {
			w.WriteHeader(500) // Server error
			return
		}
		d, err := p.get(urlDomain)
		if err != nil {
			w.WriteHeader(404) // No such domain
			return
		}

		// parse parameters
		err = r.ParseForm()
		if err != nil {
			w.WriteHeader(500) // Server error
			return
		}
		username := r.Form.Get("username")
		if username == "" {
			writeSuccess(w, false) // no username provided
			return
		}
		password := r.Form.Get("password")
		if password == "" {
			writeSuccess(w, false) // no password
			return
		}

		// user lookup
		u, err := d.get(username)
		if err != nil {
			writeSuccess(w, false) // no such user
			return
		}

		// password validation
		if u.Password != password {
			writeSuccess(w, false) // password mismatch
			return
		} else {
			writeSuccess(w, true) // successful authenticaton
			return
		}

	})
}
