/*
* @author Colton J. McCurdy
*	GitHub: mccurdyc
* Email:  mccurdyc22@gmail.com
* Date: 2018-01-22
*
*  usage:
*	 go run main.go
*  go run main.go -host localhost -port 8080
*
*  OR
*
*  go build
*  ./hello
*
 */
package main

import (
	"flag"
	"log"

	"github.com/mccurdyc/goblogs/docker-multi-stage-builds/hello/transports/http/service"
)

// global variables
var (
	serverHost = flag.String("host", "localhost", "server host")
	serverPort = flag.Int("port", 8080, "server port")
)

// init gets called before main()
func init() {
	flag.Parse()
}

func main() {
	// pass in the value of serverHost and serverPort
	s := service.NewService(*serverHost, *serverPort)
	log.Printf("started server on %s:%d", *serverHost, *serverPort)

	s.Start()
}
