package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

	routes := app.routes()

	// create a test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

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
	var tests = []struct {
		name         string
		putInSession string
		expectedHTML string
	}{
		{"first visit", "", "<small>From session:"},
		{"second visit", "Hello, World!", "<small>From session: Hello, World!"},
	}

	for _, e := range tests {

		// create a request
		req, _ := http.NewRequest("GET", "/", nil)

		// add context and session information
		req = addContextAndSessionToRequest(req, app)

		_ = app.Session.Destroy(req.Context())

		if e.putInSession != "" {
			app.Session.Put(req.Context(), "test", e.putInSession)
		}

		// response writer
		res := httptest.NewRecorder()

		// handler
		handler := http.HandlerFunc(app.Home)
		handler.ServeHTTP(res, req)

		// check status code
		if res.Code != http.StatusOK {
			t.Errorf("TestAppHome returned wrong status code; expected 200 but got %d", res.Code)
		}

		// check body
		body, _ := io.ReadAll(res.Body)
		if !strings.Contains(string(body), e.expectedHTML) {
			t.Errorf("%s: Did not find %s in response body", e.name, e.expectedHTML)
		}
	}
}

func TestApp_renderWithBadTemplate(t *testing.T) {

	// set templatepath to a location with a bad template
	pathToTemplates = "./testdata/"

	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	res := httptest.NewRecorder()

	err := app.render(res, req, "bad.page.gohtml", &TemplateData{})
	t.Log(err)
	if err == nil {
		t.Error("Expected error from bad template, but did not get one.")
	}

	// restore bad template path
	pathToTemplates = "./../../templates/"

}

func TestApp_renderParseWithBadTemplate(t *testing.T) {

	pathToTemplates = "./testdata/"

	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	res := httptest.NewRecorder()

	err := app.render(res, req, "bad.template.gohtml", &TemplateData{})
	t.Log(err)
	if err == nil {
		t.Error("Expected error when parsing bad template, but did not get one.")
	}
}

func TestApp_login(t *testing.T) {
	tests := []struct {
		name               string
		postedData         url.Values
		expectedStatusCode int
		expectedLoc        string
	}{
		{
			name: "valid login",
			postedData: url.Values{
				"email":    {"admin@example.com"}, // you might have more than one thing(e.g. checkboxes)
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/user/profile",
		},
		{
			name: "missing form data",
			postedData: url.Values{
				"email":    {"admin@example.com"}, // you might have more than one thing(e.g. checkboxes)
				"password": {""},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
		{
			name: "user not found",
			postedData: url.Values{
				"email":    {"hello@world.com"}, // you might have more than one thing(e.g. checkboxes)
				"password": {"password"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
		{
			name: "bad credentials",
			postedData: url.Values{
				"email":    {"admin@example.com"}, // you might have more than one thing(e.g. checkboxes)
				"password": {"secre"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
	}

	for _, e := range tests {

		req, _ := http.NewRequest("POST", "/login", strings.NewReader(e.postedData.Encode()))
		req = addContextAndSessionToRequest(req, app)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Login)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: returned wrong status code; expected: %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		actualLoc, err := rr.Result().Location()
		if err == nil {
			if actualLoc.String() != e.expectedLoc {
				t.Errorf("%s: returned wrong location; expected: %s, but got %s", e.name, e.expectedLoc, actualLoc.String())
			}
		} else {
			t.Errorf("%s: returned no location header", e.name)
		}
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
