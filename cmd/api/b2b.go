package main

import (
	"html/template"
	"net/http"
)

func (app *application) B2BClientPageHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\pages\\b2b\\business-page.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.

	// The last parameter to Execute() represents any dynamic data that we
	// want to pass in, which for now we'll leave as nil.
	//err = ts.ExecuteTemplate(w, "base", "smth smth smth smth data")
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
