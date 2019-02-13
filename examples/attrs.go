package main

import (
	"errors"
	"time"

	"github.com/airbloc/logger"
)

var log = logger.New("e-mail")

func main() {
	log.Info("Sending an e-mail", logger.Attrs{
		"from": "foo@bar.com",
		"to":   "qux@corge.com",
	})

	log.Info("Sending e-mail to {address}", logger.Attrs{
		"from":    "foo@bar.com",
		"address": "qux@corge.com",
	})

	err := errors.New("Too busy")

	// if the last argument (except Attrs) is `error`, it'll be displayed automatically
	log.Error("Failed to send e-mail", err, logger.Attrs{
		"from": "foo@bar.com",
		"to":   "qux@corge.com",
	})

	timer := log.Timer()
	time.Sleep(time.Millisecond * 500)
	timer.End("Created a new {} {} images", 300, "bike")
}
