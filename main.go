package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bryfry/swas/proxyauth"
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
)

// Define the command line arguments and help info
func cliInit() *cli.App {
	app := cli.NewApp()
	app.Name = "swas"
	app.Usage = "Simple Web API Server - Proxy Authentication API endpoint"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "users, u",
			Value: "./users.json",
			Usage: "Specify users json file",
		},
		cli.IntFlag{
			Name:  "port, p",
			Value: 80,
			Usage: "Specify API Port",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Increase verbosity",
		},
	}
	return app
}

// Open and handle errors associated with creating a new Proxy data model
func proxyInit(file string) *proxyauth.Proxy {
	p, err := proxyauth.NewProxy(file)
	if err != nil {
		log.WithFields(log.Fields{
			"file": file,
			"err":  err,
		}).Fatal("Proxy init failed")
	}
	return p
}

// Define the api handles, expected methods, and content type
func apiInit(p *proxyauth.Proxy) *mux.Router {
	r := mux.NewRouter()

	api := r.
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded").
		Subrouter()

	api.Handle("/api/2/domains/{domain}/proxyauth", p.Authenticate())
	return r
}

// The httpInterceptor pattern is used to intercept all http requests.  This
// enables logging on all requests before they reach the mux router
func httpInterceptor(router http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()
		router.ServeHTTP(w, r)
		elapsedTime := time.Now().Sub(startTime)

		log.WithFields(log.Fields{
			"time": elapsedTime,
		}).Debug(r.URL.Path)

	})
}

func main() {

	app := cliInit()

	app.Action = func(c *cli.Context) {

		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		log.WithFields(log.Fields{"users": c.String("users")}).Info("Proxy Auth Initalizing Users...")
		proxy := proxyInit(c.String("users"))

		log.WithFields(log.Fields{"port": c.Int("port")}).Info("Proxy Auth API Server Starting...")
		router := apiInit(proxy)
		address := fmt.Sprintf(": %d", c.Int("port"))

		// using default router with interceptor pattern (uses api mux)
		http.Handle("/", httpInterceptor(router))
		log.Fatal(http.ListenAndServe(address, nil))
	}

	app.Run(os.Args)
}
