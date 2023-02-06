package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/concierge/service/internal/data"
	"github.com/concierge/service/internal/validator"
	"github.com/go-session/session/v3"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or
			// not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been
				// sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has the type interface{}, so we use
				// fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. In turn, this will log the error using
				// our custom Logger type at the ERROR level and send the client a 500
				// Internal Server Error response.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

//func (app *application) rateLimit(next http.Handler) http.Handler {
//	// Initialize a new rate limiter which allows an average of 2 requests per second,
//	// with a maximum of 4 requests in a single ‘burst’.
//	limiter := rate.NewLimiter(2, 4)
//	// The function we are returning is a closure, which 'closes over' the limiter
//	// variable.
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Call limiter.Allow() to see if the request is permitted, and if it's not,
//		// then we call the rateLimitExceededResponse() helper to return a 429 Too Many
//		// Requests response (we will create this helper in a minute).
//		if !limiter.Allow() {
//			app.rateLimitExceededResponse(w, r)
//			return
//		}
//		next.ServeHTTP(w, r)
//	})
//}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store, err := session.Start(context.Background(), w, r)
		tokenI, ok := store.Get("Bearer")

		if !ok {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		token := fmt.Sprintf("%v", tokenI)

		//Validate the token to make sure it is in a sensible format.
		v := validator.New()
		// If the token isn't valid, use the invalidAuthenticationTokenResponse()
		// helper to send a response, rather than the failedValidationResponse() helper
		// that we'd normally use.
		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		// Retrieve the details of the user associated with the authentication token,
		// again calling the invalidAuthenticationTokenResponse() helper if no
		// matching record was found. IMPORTANT: Notice that we are using
		// ScopeAuthentication as the first parameter here.
		user, err := app.models.User.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		// Call the contextSetUser() helper to add the user information to the request
		// context.
		r = app.contextSetUser(r, user)
		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user == nil {
			//app.authenticationRequiredResponse(w, r)
			http.Redirect(w, r, "http://localhost:8080", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Checks that a user is both authenticated and activated.
func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	// Rather than returning this http.HandlerFunc we assign it to the variable fn.
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		// Check that a user is activated.
		if !user.Activated {
			//app.inactiveAccountResponse(w, r)
			http.Redirect(w, r, "http://localhost:8080", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
	// Wrap fn with the requireAuthenticatedUser() middleware before returning it.
	//print(2)
	return app.requireAuthenticatedUser(fn)
}

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the user from the request context.
		user := app.contextGetUser(r)
		// Get the slice of permissions for the user.
		//permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		//if err != nil {
		if user.UserType != "admin" {
			//app.serverErrorResponse(w, r, err)
			http.Redirect(w, r, "http://localhost:8080", http.StatusSeeOther)
			return
		}
		// Check if the slice includes the required permission. If it doesn't, then
		// return a 403 Forbidden response.
		//if !permissions.Include(code) {
		//	//app.notPermittedResponse(w, r)
		//	http.Redirect(w, r, "http://localhost:8080", http.StatusSeeOther)
		//	return
		//}
		// Otherwise they have the required permission so we call the next handler in
		// the chain.
		next.ServeHTTP(w, r)
	}
	// Wrap this with the requireActivatedUser() middleware before returning it.
	//print(1)
	return app.requireActivatedUser(fn)
}
