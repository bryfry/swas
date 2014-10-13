// proxyauth package provides data models and associated functions required to
// facilitate a proxy authentication server as specified in SPEC.md
package proxyauth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

// Proxy data model of a slice of Domains searched when authentication request is made.
// Proxy is the highest level data model in proxyauth, the golang type analogous to users.json
type Proxy struct {
	Domains []Domain
}

// Domains are identified by their Address (unique) and contain a slice of all the
// registered Users for that domain
type Domain struct {
	Address string `json:"domain"`
	Users   []User `json:"users"`
}

// Users are identified by their Username (unique) and each has a Base 64 encoded, SHA 256 digest
// of the users password.  See SPEC.md or b64sha256 for specific implementation details
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Response is used for communicating the result from an attempted authentication
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

// Parse the json users file (filePath) and return the proxy data type.
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
			log.WithFields(log.Fields{
				"domain": d.Address,
				"user":   u.Username,
			}).Info("User Initalized")
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
	w.Header().Set("Content-Type", "application/json")
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

// Proxy Authentication handler expects: domain (mux url variable) username and password
// (query parameters). HTTP response returns the appropriate json response and Status
// Code as specified in SPEC.md
func (p *Proxy) Authenticate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logfields := log.Fields{
			"Method": r.Method,
		}

		// domain lookup
		vars := mux.Vars(r)
		urlDomain, ok := vars["domain"]
		if !ok {
			w.WriteHeader(500) // Server error
			logfields["Status"] = 500
			log.WithFields(logfields).Warn("Domain parsing failed")
			return
		}
		d, err := p.get(urlDomain)
		logfields["Domain"] = urlDomain
		if err != nil {
			w.WriteHeader(404) // No such domain
			logfields["Status"] = 404
			log.WithFields(logfields).Info("No such domain")
			return
		}

		// parse parameters
		err = r.ParseForm()
		if err != nil {
			w.WriteHeader(500)
			logfields["Status"] = 500
			log.WithFields(logfields).Warn("Parse form failure")
			return
		}
		username := r.Form.Get("username")
		if username == "" {
			writeSuccess(w, false)
			log.WithFields(logfields).Info("No username provided")
			return
		}
		logfields["Username"] = username
		password := r.Form.Get("password")
		if password == "" {
			writeSuccess(w, false)
			log.WithFields(logfields).Info("No password provided")
			return
		}

		// user lookup
		u, err := d.get(username)
		if err != nil {
			writeSuccess(w, false)
			log.WithFields(logfields).Info("No such user")
			return
		}

		// password validation
		if u.Password != password {
			writeSuccess(w, false)
			log.WithFields(logfields).Info("Password mismatch")
			return
		} else {
			writeSuccess(w, true)
			log.WithFields(logfields).Info("Successful authentication")
			return
		}

	})
}
