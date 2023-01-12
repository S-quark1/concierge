package main

import (
	"html/template"
	"net/http"
)

// Add a createMovieHandler for the "POST /v1/movies" endpoint.
// return a JSON response.
func (app *application) homePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFoundResponse(w, r)
		return
	}

	files := []string{
		"./ui/html(delete)(delete)/base.html(delete)", // the order... matters?
		"./ui/html(delete)(delete)/partials/nav.html(delete)",
		"./ui/html(delete)(delete)/pages/home.html(delete)",
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
	err = ts.ExecuteTemplate(w, "base", "smth smth smth smth data")
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
