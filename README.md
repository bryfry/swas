#bryfry_swas

## Install Dependencies 
```
go get github.com/codegangsta/cli # MIT Licensed
go get github.com/gorilla/mux     # New BSD Licensed
go get github.com/Sirupsen/logrus # MIT Licensed
```

## Expected Environment
Because bryfry_swas uses an internal package (proxyauth) it is important that go 
has the appropriate access to find this package via $GOPATH.  Please place the bryfry_swas
project directory under the src directory of your $GOPATH and everything should work swimmingly!

## Build
`go build` or `go install`

## Run
`sudo ./bryfry_swas`

## Testing
Two testing methods are provided, external curl test script and internal go test 
* Run `./curl_test.sh` after startup (port 80) to run external tests
* Run `go test` from inside bryfry_swas/proxyauth to run internal package tests

## Documentation
view godoc.html to read the godoc for bryfry_swas/proxyauth

## Help and CLI options
`./proxyauth_server --help`

```
NAME:
   bryfry_swas - Simple Web API Server - Proxy Authentication API endpoint

USAGE:
   bryfry_swas [global options] command [command options] [arguments...]

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

