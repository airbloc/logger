package main

import (
	"errors"
	"github.com/airbloc/logger"
	"net/http"
)

func main() {
	log := logger.New("api")
	defer log.Recover(logger.Attrs{"currentTask": "myJob"})

	err := errors.New("some serious error")
	log.Wtf(err)

	err = http.ErrBodyNotAllowed
	log.Wtf("failed to post message to server", err)

	letsPanic()
}

func letsPanic() {
	var srv *http.Server
	_ = srv.ListenAndServe()
}
