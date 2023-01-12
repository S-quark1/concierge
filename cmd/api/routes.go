package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	//styles := http.FileServer(http.Dir("./ui"))
	//http.Handle("/ui/", http.StripPrefix("/ui/", styles))
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
