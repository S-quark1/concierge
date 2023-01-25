package main

import (
	"html/template"
	"net/http"
)

func (app *application) AdminPageHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\pages\\admin\\admin-page.html",
		// //"./ui/pages/admin/blank.html",
		// //"./ui/pages/admin/404.html", // TODO и остальное... или еще рано?
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
