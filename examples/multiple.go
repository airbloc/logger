package main

import (
	"errors"
	"time"

	"github.com/airbloc/logger"
)

var app = logger.New("app")
var images = logger.New("images")
var socket = logger.New("websocket")
var users = logger.New("users")
var db = logger.New("database")

func main() {
	go app.Info("Starting at {}", 9088)

	go db.Info("Connecting to mysql://azer@localhost:9900/foobar")
	go images.Info("Requesting an image at foo/bar.jpg")

	timer := images.Timer()
	time.Sleep(time.Millisecond * 250)
	timer.End("Fetched foo/bar.jpg")

	db.Error("Fatal connection error.")

	users.Info("{} just logged from {}", "John", "Istanbul")

	socket.Info("Connecting...")

	err := errors.New("unable to connect")
	socket.Error("error", err)

	time.Sleep(time.Millisecond * 250)

	app.Error("Failed to start, shutting down...")
}
