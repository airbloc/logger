package main

import (
	"github.com/airbloc/logger"
)

var log = logger.New("e-mail")

func main() {
	emailLog := log.WithAttrs(logger.Attrs{
		"from": "foo@bar.com",
	})

	emailLog.Info("Sending an e-mail", logger.Attrs{
		"to": "qux@corge.com",
	})

	emailLog.Warn("Bouncing an e-mail from {from}")
}
