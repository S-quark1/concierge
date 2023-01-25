package main

import (
	"github.com/concierge/service/internal/data"
	"html/template"
	"net/http"
	"time"
)

func (app *application) LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/index.html", // the order... matters?
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

func (app *application) RegFormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	email := r.FormValue("emailReg")
	companyName := r.FormValue("companyNameReg")
	phone := r.FormValue("phoneReg")

	regForm := &data.RegForm{
		CompanyName: companyName,
		Email:       email,
		PhoneNumber: phone,
		CreatedAt:   time.Time{},
	}

	err = app.models.RegForm.Insert(regForm)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	print("а дальше что?")
}
