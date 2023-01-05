package main

import (
	"html/template"
	"net/http"
)

// Home handler
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "this is the home page")
	_ = app.render(w, r, "home.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP   string
	Data map[string]any
}

// renderer
func (app *application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	// parse the template from disk
	parsedTemplate, err := template.ParseFiles("./templates/" + t) // 현재 디렉토리에 없는 거 아닌가?
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return err
	}

	// execute the template, passing it data, if any
	err = parsedTemplate.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
