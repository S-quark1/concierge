package main

import (
	"html/template"
	"net/http"
)

func (app *application) AddServicesPageHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		//"./ui/pages/admin/admin-page.html", // todo нужны html
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		// service
		Name              string `json:"name"`
		Description       string `json:"description"`
		ServiceType       string `json:"type"`
		CreatedByEmployee int64  `json:"createdBy_id"`
		CompanyProviding  int32  `json:"company_id"`

		//prices
		Prices   []int32  `json:"prices"`
		UserType []string `json:"user_types"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//service := &data.Service{
	//	Name:              input.Name,
	//	Description:       input.Description,
	//	ServiceType:       input.ServiceType,
	//	CreatedByEmployee: input.CreatedByEmployee,
	//	CompanyProviding:  input.CompanyProviding,
	//}

	// todo
	//v := validator.New()
	//if data.ValidateService(v, service); !v.Valid() {
	//	app.failedValidationResponse(w, r, v.Errors)
	//	return
	//}

	// check if company exists
	//err = app.models.Company.Exists(service.CompanyProviding)
	//if err != nil {
	//	app.notFoundResponse(w, r)
	//	return
	//}
	//
	//err = app.models.Service.Insert(service)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//	return
	//}
	//
	//// insert prices
	//for index, element := range input.Prices {
	//	price := &data.Price{
	//		ServiceID: service.ID,
	//		Price:     element,
	//		UserType:  input.UserType[index],
	//	}
	//	// index is the index where we are
	//	// element is the element from someSlice for where we are
	//
	//	err = app.models.Price.Insert(price)
	//	if err != nil {
	//		app.serverErrorResponse(w, r, err)
	//		return
	//	}
	//}

	// todo вместо nil возвращаем service and its prices
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
