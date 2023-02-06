package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.ServeFiles("/css/*filepath", http.Dir("./ui/css/"))
	router.ServeFiles("/img/*filepath", http.Dir("./ui/img/"))
	router.ServeFiles("/js/*filepath", http.Dir("./ui/js/"))
	router.ServeFiles("/pages/*filepath", http.Dir("./ui/pages/"))
	router.ServeFiles("/scss/*filepath", http.Dir("./ui/scss/"))
	router.ServeFiles("/vendor/*filepath", http.Dir("./ui/vendor/"))

	//router.ServeFiles("/css/*filepath", http.Dir("C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\css"))
	//router.ServeFiles("/img/*filepath", http.Dir("C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\img"))
	//router.ServeFiles("/js/*filepath", http.Dir("C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\js"))
	//router.ServeFiles("/pages/*filepath", http.Dir("C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\pages"))
	//router.ServeFiles("/scss/*filepath", http.Dir("C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\scss"))
	//router.ServeFiles("/vendor/*filepath", http.Dir("C:\\Users\\mapol\\IdeaProjects\\concierge\\ui\\vendor"))

	// LandingPage
	router.HandlerFunc(http.MethodGet, "/", app.showLandingPageHandler)
	router.HandlerFunc(http.MethodPost, "/regForm", app.RegFormHandler)
	router.HandlerFunc(http.MethodPost, "/login", app.LoginHandler)

	// Admin
	router.HandlerFunc(http.MethodGet, "/my-cabinet-admin/services", app.requirePermission("is_admin", app.GetAddServicesPageHandler))
	router.HandlerFunc(http.MethodPost, "/my-cabinet-admin/services", app.requirePermission("is_admin", app.PostAddServicesHandler))
	router.HandlerFunc(http.MethodGet, "/my-cabinet-admin/", app.requirePermission("is_admin", app.showAdminPageHandler))
	router.HandlerFunc(http.MethodGet, "/my-cabinet-admin/analytics", app.requirePermission("is_admin", app.showAdminRegisterUsersPageHandler))

	// B2B
	router.HandlerFunc(http.MethodGet, "/my-cabinet-b-client", app.requirePermission("is_B2B", app.B2BClientPageHandler))

	// B2C
	router.HandlerFunc(http.MethodGet, "/my-cabinet", app.requirePermission("is_B2С", app.B2CClientPageHandler))

	// Concierge
	//router.HandlerFunc(http.MethodGet, "/my-cabinet", app.CSPageHandler)

	// Partner

	//mux := http.NewServeMux()
	//fileServer := http.FileServer(http.Dir("./ui/static(delete)/"))
	//mux.Handle("/static(delete)/", http.StripPrefix("/static(delete)", fileServer))
	//
	//mux.HandleFunc("/", app.homePageHandler)
	//mux.HandleFunc("/auth", app.authorizationPageHandler)

	router.HandlerFunc(http.MethodPost, "/debug", app.PostAddUsersHandler)
	router.HandlerFunc(http.MethodPost, "/debug/token", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.authenticate(router))
}
