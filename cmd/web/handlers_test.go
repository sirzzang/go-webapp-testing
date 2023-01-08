package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTests = []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/foo", http.StatusNotFound},
	}

	// var app application
	routes := app.routes()

	// create a test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// fix template directory path
	pathToTemplates = "./../../templates/"

	// range through test data
	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url) // built-in server root url + "/"
		t.Log(ts.URL + e.url)                        // http://127.0.0.1:51120/home
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s: expected status %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestAppHome(t *testing.T) {
	// create a request
	req, _ := http.NewRequest("GET", "/", nil)

	// add context and session information
	req = addContextAndSessionToRequest(req, app)

	// response writer
	res := httptest.NewRecorder()

	// handler
	handler := http.HandlerFunc(app.Home)
	handler.ServeHTTP(res, req)

	// check status code
	if res.Code != http.StatusOK {
		t.Errorf("TestAppHome returned wrong status code. Expected 200 but got %d", res.Code)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
	return ctx
}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {

	// add context to request
	req = req.WithContext(getCtx(req))

	// add information to session
	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))

	return req.WithContext(ctx)

}
