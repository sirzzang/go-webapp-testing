package main

import (
	"os"
	"testing"
)

// all the references of app
var app application

func TestMain(m *testing.M) {

	app.Session = getSession()
	os.Exit(m.Run())
}
