package proxyauth

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
)

// We need access to the top level gorilla/mux router which, being a mux, also
// implements the ServeHTTP(w,r) interface which we will call in the tests
func routerInit() *mux.Router {
	p, err := NewProxy("../users.json")
	if err != nil {
		log.Fatalf("Proxy init failed: %s", err)
	}

	r := mux.NewRouter()
	// define the api method expectations, useful for future api handles
	api := r.
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded").
		Subrouter()

	api.Handle("/api/2/domains/{domain}/proxyauth", p.Authenticate())
	return r
}

func TestAuthenticate(t *testing.T) {

	var tests = []struct {
		domain   string
		username string
		password string
		code     int
		success  bool
	}{
		// Case1 Success, topcoder.com domain, StatusCode 200
		{"topcoder.com", "takumi", "{SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=", 200, true},

		// Case2 Success, appirio.com domain, StatusCode 200
		{"appirio.com", "jun", "{SHA256}/Hnfw7FSM40NiUQ8cY2OFKV8ZnXWAvF3U7/lMKDwmso=", 200, true},

		// Case3 Failure, password unmatch, StatusCode 200
		{"topcoder.com", "takumi", "{SHA256}/Hnfw7FSM40NiUQ8cY2OFKV8ZnXWAvF3U7/lMKDwmso=", 200, false},

		// Case4 Failure, username not found, StatusCode 200
		{"topcoder.com", "bryfry", "{SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=", 200, false},

		// Case5 Failure, domain not found, StatusCode 404
		{"bryfry.com", "takumi", "{SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=", 404, false},
	}

	r := routerInit()

	for _, test := range tests {

		// setup request and parameters
		data := url.Values{}
		if test.username != "" {
			data.Set("username", test.username)
		}
		if test.password != "" {
			data.Set("password", test.password)
		}
		url := "/api/2/domains/" + test.domain + "/proxyauth"
		req, _ := http.NewRequest("POST", url, bytes.NewBufferString(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// record request response
		rw := httptest.NewRecorder()
		rw.Body = new(bytes.Buffer)

		// make request
		r.ServeHTTP(rw, req)

		// ensure got (g) equals want (w)
		if g, w := rw.Code, test.code; g != w {
			t.Errorf("%s: code = %d, want %d", url, g, w)
		}
		if rw.Code == 200 {
			if g, w := rw.Body.String(), successBody(test.success); g != w {
				t.Errorf("%s: body = %q, want %q", url, g, w)
			}
		}
	}

}
