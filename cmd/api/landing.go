package main

import (
	"github.com/concierge/service/internal/data"
	"html/template"
	"net/http"
	"strings"
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
		return
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
		//app.serverErrorResponse(w, r, err)
		http.Redirect(w, r, "http://localhost:8080", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "http://localhost:8080", http.StatusAccepted)
}

// post
func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//username := strings.Split(r.FormValue("usernameD"), " ")
	pswd := strings.Split(r.FormValue("passwordD"), " ")

	user, err := app.models.User.GetByUsername("ivan55")
	if err != nil {
		app.notFoundResponse(w, r)
		//http.Redirect(w, r, "http://localhost:8080", http.StatusNotFound)
		return
	}

	//
	//err = user.Password.Set(pswd[0])
	//if err != nil {
	//	return
	//}
	//err = app.models.User.Update(user)
	//if err != nil {
	//	return
	//}

	//print(username)
	print(user.UserType)

	match, err := user.Password.Matches(pswd[0])
	if err != nil {
		app.serverErrorResponse(w, r, err)
		http.Redirect(w, r, "http://localhost:8080", http.StatusInternalServerError)
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		http.Redirect(w, r, "http://localhost:8080", http.StatusUnauthorized)
	}

	_, err = app.models.Token.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if user.UserType == "admin" {
		http.Redirect(w, r, "http://localhost:8080/my-cabinet-admin", http.StatusOK)
	} else if user.UserType == "cs_employee" {
		http.Redirect(w, r, "http://localhost:8080/my-cabinet", http.StatusOK)
	}
	print("huh")
	http.Redirect(w, r, "http://localhost:8080", http.StatusSeeOther)
}
