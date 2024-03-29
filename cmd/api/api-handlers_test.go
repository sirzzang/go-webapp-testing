package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	"webapp/pkg/data"

	"github.com/go-chi/chi"
)

func Test_app_authenticate(t *testing.T) {
	var tests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{"valid user", `{"email":"admin@example.com","password":"secret"}`, http.StatusOK},
		{"not json", `"}`, http.StatusUnauthorized},
		{"not json", `I'm not json`, http.StatusUnauthorized},
		{"empty json", `I'm not json`, http.StatusUnauthorized},
		{"empty email", `{"email":""}`, http.StatusUnauthorized},
		{"empty password", `{"email":"admin@exmample.com"}`, http.StatusUnauthorized},
		{"invalid user", `{"email":"admin@someotherdomain.com"}`, http.StatusUnauthorized},
	}

	for _, e := range tests {
		var reader io.Reader
		reader = strings.NewReader(e.requestBody)
		req, _ := http.NewRequest("POST", "/auth", reader)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.authenticate)

		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: expeceted code %v, but got wrong status code %v", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}

func Test_app_refresh(t *testing.T) {
	tests := []struct {
		name               string
		token              string
		expectedStatusCode int
		resetRefreshTime   bool
	}{
		{"valid token", "", http.StatusOK, true},
		{"valid but not expired", "", http.StatusTooEarly, false},
		{"expired token", expiredToken, http.StatusBadRequest, false},
	}

	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}

	// save the old refresh token expiry
	oldRefreshTime := refreshTokenExpiry

	for _, e := range tests {
		var tkn string
		if e.token == "" {
			if e.resetRefreshTime {
				refreshTokenExpiry = time.Second * 1
			}
			tokens, _ := app.generateTokenPairs(&testUser)
			tkn = tokens.RefreshToken
		} else {
			tkn = e.token
		}

		postedData := url.Values{
			"refresh_token": []string{tkn},
		}

		req, _ := http.NewRequest("POST", "/refresh-token", strings.NewReader(postedData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.refresh)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: expeceted code %v, but got wrong status code %v", e.name, e.expectedStatusCode, rr.Code)
		}

		// recover the old refresh token expiry
		refreshTokenExpiry = oldRefreshTime

	}
}

func Test_app_UserHandlers(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		json               string
		paramId            string
		handler            http.HandlerFunc
		expectedStatusCode int
	}{
		{"allUsers", "GET", "", "", app.allUsers, http.StatusOK},
		{"deleteUser", "DELETE", "", "1", app.deleteUser, http.StatusNoContent},
		// {"allUsers", "GET", "", "", app.allUsers, http.StatusOK},
		// {"allUsers", "GET", "", "", app.allUsers, http.StatusOK},

	}

	for _, e := range tests {
		var req *http.Request
		if e.json == "" {
			req, _ = http.NewRequest(e.method, "/users", nil)
		} else {
			req, _ = http.NewRequest(e.method, "/users", strings.NewReader(e.json))
		}

		if e.paramId != "" {
			log.Printf("paramId: %s", e.paramId)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("userID", e.paramId)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(e.handler)
		handler.ServeHTTP(rr, req)

		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: expeceted code %v, but got wrong status code %v", e.name, e.expectedStatusCode, rr.Code)
		}

	}
}
