package main

import (
	"github.com/concierge/service/internal/data"
	"github.com/concierge/service/internal/validator"
	"html/template"
	"net/http"
	"time"
)

func (app *application) showAdminPageHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/pages/admin/admin-page.html",
		"./ui/pages/admin/analytics.html",
		//"C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\pages\\admin\\admin-page.html",
		//"./ui/pages/admin/404.html",
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

// get
func (app *application) showAdminRegisterUsersPageHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		//"./ui/pages/admin/admin-page.html", //todo нужны html
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	regForms, err := app.models.RegForm.GetByDateAsc()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = ts.Execute(w, regForms)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// get function. only for the ui
func (app *application) GetAddServicesPageHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		//"./ui/pages/admin/admin-page.html", // todo нужны html
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

// post function. actual data adding
func (app *application) PostAddServicesHandler(w http.ResponseWriter, r *http.Request) {
	//files := []string{
	//	//"./ui/pages/admin/admin-page.html", // todo нужны html
	//}
	//
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//	return
	//}

	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	var input struct {
		// service
		Name             string `json:"name"`
		Description      string `json:"description"`
		Type             string `json:"type"`
		CreatedByID      int64  `json:"createdBy_id"`
		CompanyProviding int    `json:"company_id"`

		//prices
		Prices   []int    `json:"prices"`
		UserType []string `json:"user_types"`
	}

	input.Name = r.FormValue("nameService")
	input.Description = r.FormValue("descService")
	input.Type = r.FormValue("typeService")
	input.CreatedByID = int64(app.convertInt(r.FormValue("createdByIDService")))
	input.CompanyProviding = app.convertInt(r.FormValue("companyProvidingService"))

	for index, element := range r.Form["pricesService"] {
		input.Prices[index] = app.convertInt(element)
		input.UserType[index] = element
	}
	//input.Prices[0] = app.convertInt(r.Form["pricesService"][0])
	//input.Prices[1] = app.convertInt(r.Form["pricesService"][1])
	//input.Prices[2] = app.convertInt(r.Form["pricesService"][2])
	//
	//input.UserType[0] = r.Form["userTypeService"][0]
	//input.UserType[1] = r.Form["userTypeService"][1]
	//input.UserType[2] = r.Form["userTypeService"][2]

	service := &data.Service{
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,
		CreatedByID: input.CreatedByID,
		CompanyID:   input.CompanyProviding,
	}

	//v := validator.New()
	//if data.ValidateService(v, service); !v.Valid() {
	//	app.failedValidationResponse(w, r, v.Errors)
	//	return
	//}

	// check if company exists
	err = app.models.Company.Exists(service.CompanyID)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Service.Insert(service)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// insert prices
	for index, element := range input.Prices {
		price := &data.Price{
			ServiceID: service.ID,
			Price:     element,
			UserType:  input.UserType[index],
		}
		// index is the index where we are
		// element is the element from someSlice for where we are

		err = app.models.Price.Insert(price)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// todo тут что-то нужно возвращать
	http.Redirect(w, r, "http://localhost:8080/my-cabinet/services", http.StatusOK)
}

func (app *application) PostAddUsersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		UserType    string `json:"user_type"`
		Preferences string `json:"preferences"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		Username:    input.Username,
		Activated:   false,
		UserType:    input.UserType,
		Preferences: input.Preferences,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.User.Insert(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	_, err = app.models.Token.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//app.background(func() {
	//	// As there are now multiple pieces of data that we want to pass to our email
	//	// templates, we create a map to act as a 'holding structure' for the data. This
	//	// contains the plaintext version of the activation token for the user, along
	//	// with their ID.
	//	Data := map[string]interface{}{
	//		"pswd":     input.Password,
	//		"username": user.Username,
	//	}
	//
	//	err = app.mailer.Send(user.Email, "user_welcome.tmpl", Data)
	//	if err != nil {
	//		// Importantly, if there is an error sending the email then we use the
	//		// app.logger.PrintError() helper to manage it, instead of the
	//		// app.serverErrorResponse() helper like before.
	//		app.logger.Print(err)
	//	}
	//})

	// todo тут что-то нужно возвращать
	http.Redirect(w, r, "http://localhost:8080/my-cabinet/services", http.StatusOK)
}
