package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
	"webapp/pkg/data"
)

func Test_application_handlers(t *testing.T) {
	var theTests = []struct {
		name                    string
		url                     string
		expectedStatusCode      int
		expectedURL             string
		expectedFirstStatusCode int
	}{
		{"home", "/", http.StatusOK, "/", http.StatusOK},
		{"404", "/foo", http.StatusNotFound, "/foo", http.StatusNotFound},
		{"profile", "/user/profile", http.StatusOK, "/", http.StatusTemporaryRedirect},
	}

	routes := app.routes()

	// create a test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url) // built-in server root url + "/"
		t.Log(ts.URL + e.url)                        // http://127.0.0.1:51120/home
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("%s: expected status %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

		if resp.Request.URL.Path != e.expectedURL {
			t.Errorf("%s: expected final url of %s, but got %s", e.name, e.expectedURL, resp.Request.URL.Path)
		}

		resp2, _ := client.Get(ts.URL + e.url)
		if resp2.StatusCode != e.expectedFirstStatusCode {
			t.Errorf("%s: expected first returned status code to be %d, but got %d", e.name, e.expectedFirstStatusCode, resp2.StatusCode)
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

func TestApp_uploadFiles(t *testing.T) {
	// set up pipes
	pr, pw := io.Pipe()

	// create a new writer, of type *io.Writer
	writer := multipart.NewWriter(pw)

	// create a waitgroup, and add 1
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// simulate uploading a file using a goroutine and our writer
	go simulatePNGUpload("./testdata/img.png", writer, t, wg)

	// read from the pipe which receives data
	request := httptest.NewRequest("POST", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	// call app.UploadFiles
	uploadedFiles, err := app.uploadFiles(request, "./testdata/uploads/")
	if err != nil {
		t.Error(err)
	}

	// perform tests
	if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].OriginalFileName)); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", err.Error())
	}

	// clean up
	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].OriginalFileName))

	wg.Wait()
}

func TestApp_UploadProfilePic(t *testing.T) {
	uploadPath = "./testdata/uploads"
	fileName := "img.png"
	filePath := fmt.Sprintf("./testdata/%s", fileName)

	// specify a field name for the form
	fieldName := "file"

	// create bytes buffer to act as the request body
	body := new(bytes.Buffer)

	// create a multi writer
	mw := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}

	w, err := mw.CreateFormFile(fieldName, filePath)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(w, file); err != nil {
		t.Fatal(err)
	}

	mw.Close()

	// mock request
	req := httptest.NewRequest("POST", "/upload", body)
	req = addContextAndSessionToRequest(req, app)
	app.Session.Put(req.Context(), "user", data.User{ID: 1})
	req.Header.Add("Content-Type", mw.FormDataContentType())

	// record response
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.UploadProfilePic)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("wrong status code")
	}

	_ = os.Remove(fmt.Sprintf("%s/%s", uploadPath, fileName))
}

func simulatePNGUpload(fileToUpload string, writer *multipart.Writer, t *testing.T, wg *sync.WaitGroup) {
	defer writer.Close()
	defer wg.Done()

	// create the form data filled 'file' with value being filename
	part, err := writer.CreateFormFile("file", path.Base(fileToUpload))
	if err != nil {
		t.Error(err)
	}

	// open the actual file
	f, err := os.Open(fileToUpload)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	// decode the image
	img, _, err := image.Decode(f)
	if err != nil {
		t.Error("error decoding image:", err)
	}

	// write the png to our io.Writer
	err = png.Encode(part, img)
	if err != nil {
		t.Error(err)
	}
}
