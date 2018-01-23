/*
* @author Colton J. McCurdy
*	GitHub: mccurdyc
* Email:  mccurdyc22@gmail.com
* Date: 2018-01-22
 */

package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/mccurdyc/goblogs/docker-multi-stage-builds/hello/transports/http/handlers"
)

// Service contains an http Server and the time when it was initially launched.
type Service struct {
	Launched time.Time
	Server   http.Server
}

// NewService intiallizes a new Service with the host where the server should run
// and the port it is running on set.
func NewService(host string, port int) *Service {
	addr := fmt.Sprintf("%s:%d", host, port)

	return &Service{
		Launched: time.Now(),
		Server: http.Server{
			Addr:         addr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

// Start creates a new router, r, where routes will be mapped to handler functions
// and the accepted routes and respective handler functions are defined. Finally,
// the server is started.
func (s *Service) Start() {
	r := mux.NewRouter()

	r.HandleFunc("/hello", handlers.Hello)
	http.Handle("/", r)

	if err := s.Server.ListenAndServe(); err != nil {
		panic(err)
	}
}
