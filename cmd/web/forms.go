package main

import (
	"net/url"
	"strings"
)

// convenience type to have functions tied to our map
type errors map[string][]string

func (e errors) Get(field string) string {
	errorSlice := e[field]
	if len(errorSlice) == 0 {
		return ""
	}
	return errorSlice[0]
}

// Add adds an error message for a given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Form is a type used to instantiate form validation
type Form struct {
	Data   url.Values
	Errors errors
}

// NewForm instantiates a form struct
func NewForm(data url.Values) *Form {
	return &Form{
		Data:   data,
		Errors: map[string][]string{},
	}
}

// Has validates if a form has a given field
func (f *Form) Has(field string) bool {
	x := f.Data.Get(field)
	if x == "" {
		return false
	}
	return true
}

// Required valiedates for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Data.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) Check(ok bool, key, message string) {
	if !ok {
		f.Errors.Add(key, message)
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
