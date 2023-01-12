package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.ServeFiles("/img/*filepath", http.Dir("./ui/img/"))
	router.ServeFiles("/js/*filepath", http.Dir("./ui/js/"))
	router.ServeFiles("/pages/*filepath", http.Dir("./ui/pages/"))
	router.ServeFiles("/scss/*filepath", http.Dir("./ui/scss/"))
	router.ServeFiles("/vendor/*filepath", http.Dir("./ui/vendor/"))
	// TODO если что нужно тут добавляем
	// LandingPage
	router.HandlerFunc(http.MethodGet, "/", app.LandingPageHandler)

	// Admin

	// B2B

	// B2C

	// Concierge

	// Partner

	//mux := http.NewServeMux()
	//fileServer := http.FileServer(http.Dir("./ui/static/"))
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	//
	//mux.HandleFunc("/", app.homePageHandler)
	//mux.HandleFunc("/auth", app.authorizationPageHandler)

	return router
}
