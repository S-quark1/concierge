package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.ServeFiles("/css/*filepath", http.Dir("C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\css"))
	router.ServeFiles("/img/*filepath", http.Dir("C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\img"))
	router.ServeFiles("/js/*filepath", http.Dir("C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\js"))
	router.ServeFiles("/pages/*filepath", http.Dir("C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\pages"))
	router.ServeFiles("/scss/*filepath", http.Dir("C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\scss"))
	router.ServeFiles("/vendor/*filepath", http.Dir("C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\vendor"))
	// TODO если что нужно тут добавляем

	// LandingPage
	router.HandlerFunc(http.MethodGet, "/", app.LandingPageHandler)

	// Admin
	router.HandlerFunc(http.MethodGet, "/admin", app.AdminPageHandler)

	// B2B

	// B2C

	// Concierge

	// Partner

	//mux := http.NewServeMux()
	//fileServer := http.FileServer(http.Dir("./ui/static(delete)/"))
	//mux.Handle("/static(delete)/", http.StripPrefix("/static(delete)", fileServer))
	//
	//mux.HandleFunc("/", app.homePageHandler)
	//mux.HandleFunc("/auth", app.authorizationPageHandler)

	return router
}
