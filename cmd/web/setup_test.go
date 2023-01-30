package main

import (
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"
)

// all the references of app
var app application

func TestMain(m *testing.M) {

	pathToTemplates = "./../../templates/"
	app.Session = getSession()

	app.DB = &dbrepo.TestDBRepo{}

	os.Exit(m.Run())

}
