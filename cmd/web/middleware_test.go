package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"webapp/pkg/data"
)

func Test_application_addIPToContext(t *testing.T) {
	tests := []struct {
		name        string
		headerName  string
		headerValue string
		addr        string
		emptyAddr   bool // should be ignored
	}{
		{"default", "", "", "", false},
		{"empty", "", "", "", true},
		{"forwarded", "X-Forwarded-For", "192.188.159.1", "", false},
		{"invalid port", "", "", "hello:world", false},
	}

	// create a dummy handler we'll use to check the context
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// make sure that the value exists in context
		val := r.Context().Value(contextUserKey)
		if val == nil {
			t.Errorf("%s not present in context", contextUserKey)
		}

		// make sure we got a string back
		ip, ok := val.(string)
		if !ok {
			t.Errorf("%v not string", ip)
		}
		t.Log(ip)
	})

	for _, e := range tests {

		// create the handler to test
		handlerToTest := app.addIPToContext(nextHandler)

		// mock request with test case values
		req := httptest.NewRequest("GET", "http://testing", nil)
		if e.emptyAddr {
			req.RemoteAddr = ""
		}
		if len(e.headerName) > 0 {
			req.Header.Add(e.headerName, e.headerValue)
		}
		if len(e.addr) > 0 {
			req.RemoteAddr = e.addr
		}

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}
}

func Test_application_ipFromContext(t *testing.T) {

	// get a context
	var ctx = context.Background()

	// put something in the context
	ctx = context.WithValue(ctx, contextUserKey, "192.168.159.21")

	// call the function
	ip := app.ipFromContext(ctx)

	// perform the test
	if !strings.EqualFold("192.168.159.21", ip) {
		t.Error("Wrong value returned from context")
	}

}

func TestApp_auth(t *testing.T) {

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

	tests := []struct {
		name            string
		isAuthenticated bool
		expectedCode    int
	}{
		{"logged in", true, http.StatusOK},
		{"not logged in", false, http.StatusTemporaryRedirect},
	}

	for _, e := range tests {

		handlerToTest := app.auth(nextHandler)

		req := httptest.NewRequest("GET", "http://testing", nil)
		req = addContextAndSessionToRequest(req, app)
		if e.isAuthenticated {
			app.Session.Put(req.Context(), "user", data.User{ID: 1})
		}

		rr := httptest.NewRecorder()
		handlerToTest.ServeHTTP(rr, req)

		if e.isAuthenticated && rr.Code != e.expectedCode {
			t.Errorf("%s: expected status code %d, but got %d", e.name, e.expectedCode, rr.Code)
		}

		if !e.isAuthenticated && rr.Code != e.expectedCode {
			t.Errorf("%s: expected status code %d, but got %d", e.name, e.expectedCode, rr.Code)
		}

	}

}
