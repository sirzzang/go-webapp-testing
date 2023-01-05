package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

	var app application

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

	// create an app var of type application
	var app application

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
