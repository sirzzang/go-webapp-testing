package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Has(t *testing.T) {

	// test empty form
	form := NewForm(nil)

	has := form.Has("foo")
	if has {
		t.Error("Form has field when it should not.")
	}

	// test form with data
	postedData := url.Values{}
	postedData.Add("foo", "foo bar")
	form = NewForm(postedData)

	has = form.Has("foo")
	if !has {
		t.Error("Form does not have field when it should.")
	}
}

func TestForm_Required(t *testing.T) {

	// test empty form
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := NewForm(r.PostForm)

	form.Required("foo", "bar", "baz")

	if form.Valid() {
		t.Error("Form shows valid when required fields are missing.")
	}

	// test form with data
	postedData := url.Values{}
	postedData.Add("foo", "foo")
	postedData.Add("bar", "bar")
	postedData.Add("baz", "baz")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData

	form = NewForm(r.PostForm)
	form.Required("foo", "bar", "baz")
	if !form.Valid() {
		t.Error("Form does not have required fields when it should.")
	}

}

func TestForm_Check(t *testing.T) {

	form := NewForm(nil)

	form.Check(false, "password", "password is required")
	if form.Valid() {
		t.Error("Valid() returns false, and it should be true when calling Check()")
	}

}

func TestForm_Get(t *testing.T) {

	form := NewForm(nil)

	form.Check(false, "password", "password is required")
	s := form.Errors.Get("password")
	if len(s) == 0 {
		t.Error("Should have an error returned from Get, but do not.")
	}

	s = form.Errors.Get("whatever")
	if len(s) != 0 {
		t.Error("Should not have an error, but got one.")
	}
}
