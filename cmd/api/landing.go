package main

import (
	"context"
	"errors"
	"github.com/concierge/service/internal/data"
	"github.com/go-session/session/v3"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func (app *application) showLandingPageHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/index.html", // the order... matters?
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

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

	email := strings.Split(r.FormValue("emailD"), " ")
	pswd := strings.Split(r.FormValue("passwordD"), " ")

	user, err := app.models.User.GetByEmail(email[0])
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
			app.logError(r, err)
			//http.Redirect(w, r, "http://localhost:8080", http.StatusNotFound)
		}
		return
	}

	match, err := user.Password.Matches(pswd[0])
	if err != nil {
		app.serverErrorResponse(w, r, err)
		http.Redirect(w, r, "http://localhost:8080", http.StatusSeeOther)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		http.Redirect(w, r, "http://localhost:8080", http.StatusSeeOther)
		return
	}

	token, err := app.models.Token.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//d := "Bearer " + token.Plaintext

	store, err := session.Start(context.Background(), w, r)
	store.Set("Bearer", token.Plaintext)

	err = store.Save()
	//r.Header.Add("Vary", "Authorization")
	//r.Header.Add("Authorization", d)

	//d := []string{"Bearer ", token.Plaintext}
	//r = app.contextSetUser(r, user)
	//w.Header()["Authorization"] = d

	if user.UserType == "admin" {
		http.Redirect(w, r, "http://localhost:8080/my-cabinet-admin/", http.StatusSeeOther)
		return
	} else if user.UserType == "cs_employee" {
		http.Redirect(w, r, "http://localhost:8080/my-cabinet", http.StatusSeeOther)
		return
	}
	//http.Redirect(w, r, "http://localhost:8080", http.StatusSeeOther)
}
