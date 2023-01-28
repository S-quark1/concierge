package main

import (
	"github.com/concierge/service/internal/data"
	"html/template"
	"net/http"
)

// get function. only for the ui
func (app *application) AddServicesPageHandler(w http.ResponseWriter, r *http.Request) {
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
func (app *application) AddServicesHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		//"./ui/pages/admin/admin-page.html", // todo нужны html
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = r.ParseForm()
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

	// todo вместо nil возвращаем service and its prices
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
