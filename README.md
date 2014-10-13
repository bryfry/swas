#bryfry_swas

## Checkout
```
go get github.com/bryfry/swas
```

## Install Dependencies 
```
go get github.com/codegangsta/cli # MIT Licensed
go get github.com/gorilla/mux     # New BSD Licensed
go get github.com/Sirupsen/logrus # MIT Licensed
```

## Build
`go build` or `go install`

## Run
`sudo ./swas`

## Testing
Two testing methods are provided, external curl test script and internal go test 
* Run `./curl_test.sh` after startup (port 80) to run external tests
* Run `go test` from inside swas/proxyauth to run internal package tests

## Documentation
view godoc.html to read the godoc for bryfry_swas/proxyauth

## Help and CLI options
`./proxyauth_server --help`

```
NAME:
   swas - Simple Web API Server - Proxy Authentication API endpoint

USAGE:
   swas [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --users, -u './users.json'	Specify users json file
   --port, -p '80'		Specify API Port
   --verbose			Increase verbosity
   --help, -h			show help
   --version, -v		print the version

```
## More Info
See the SPEC.md file for full API specification

