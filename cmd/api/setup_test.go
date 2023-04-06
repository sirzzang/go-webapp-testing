package main

import (
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"
)

var app application
var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXVkIjoiZXhhbXBsZS5jb20iLCJleHAiOjE2Nzk1NTg4MjQsIm5hbWUiOiJKb2huIERvZSIsInN1YiI6IjEifQ.CJ8t5glWSTIdgGbIYE_bUytbGAMlrUgyQJbmCnKESUE"

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBRepo{}
	app.Domain = "example.com"
	app.JWTSecret = "eraser-secret"
	os.Exit(m.Run())
}
