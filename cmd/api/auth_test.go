package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"webapp/pkg/data"
)

func TestApp_getTokenFromHeaderAndVerify(t *testing.T) {
	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}

	tokens, _ := app.generateTokenPairs(&testUser)

	tests := []struct {
		name      string
		token     string
		isError   bool
		setHeader bool
		issuer    string
	}{
		{"valid", fmt.Sprintf("Bearer %s", tokens.Token), false, true, app.Domain},
		{"valid expired", fmt.Sprintf("Bearer %s", expiredToken), true, true, app.Domain},
		{"no header", "", true, false, app.Domain},
		{"invalid token", fmt.Sprintf("Bearer %s111", tokens.Token), true, true, app.Domain},
		{"no bearer in token", fmt.Sprintf("Bear %s", tokens.Token), true, true, app.Domain},
		{"three header parts", fmt.Sprintf("Bearer %s 11", tokens.Token), true, true, app.Domain},
		// make sure the next test is the last one to run
		{"wrong issuer", fmt.Sprintf("Bearer %s", tokens.Token), true, true, "abcd"},
	}

	for _, e := range tests {
		if e.issuer != app.Domain {
			// generate new token with wrong issuer
			app.Domain = e.issuer
			tokens, _ = app.generateTokenPairs(&testUser)
		}

		req, _ := http.NewRequest("GET", "/", nil)
		if e.setHeader {
			req.Header.Set("Authorization", e.token)
		}

		rr := httptest.NewRecorder()

		_, _, err := app.getTokenFromHeaderAndVerify(rr, req)
		if !e.isError && err != nil {
			t.Errorf("%s: did not expect error, but got an error - %s", e.name, err.Error())
		}
		if e.isError && err == nil {
			t.Errorf("%s: expected error, but got no error", e.name)
		}

		// recover domain
		app.Domain = "example.com"

	}
}
