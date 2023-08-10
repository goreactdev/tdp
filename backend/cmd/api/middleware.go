package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ton-developer-program/internal/database"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}


func (app *application) NoDirListingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/") {
					http.NotFound(w, r)
					return
			}
			h.ServeHTTP(w, r)
	})
}



func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, database.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}


		headerParts := strings.Split(authorizationHeader, " ")

		if len(headerParts) == 2 && headerParts[0] == "Bearer" {
			tokenString := headerParts[1]


			if tokenString == "" {
				app.invalidAuthenticationToken(w, r)
				return
			}


			user, err := app.sqlModels.Users.GetForToken(database.ScopeAuthentication, tokenString)
			if err != nil {
				switch {
				case errors.Is(err, sql.ErrNoRows):
					app.notPermittedResponse(w, r)
				default:
					app.serverError(w, r, err) 
					app.logger.Error(err,nil)
				}
				return
			}

			// get user permissions
			permissions, err := app.sqlModels.Users.GetAllPermissions(user.ID)
			if err != nil {
				app.serverError(w, r, err) 
				app.logger.Error(err,nil)
				return
			}

			user.Permissions = permissions


			r = app.contextSetUser(r, user)

			next.ServeHTTP(w, r)

		}
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			app.notPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}


func (app *application) requirePermission(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user == nil {
			app.authenticationRequired(w, r)
			return
		}

		// Split the URL path into segments.
		pathSegments := strings.Split(r.URL.Path, "/")
		if len(pathSegments) < 3 {
			http.Error(w, "Invalid route", http.StatusNotFound)
			return
		}

		// if path segments contain "admin" then check if user is admin
		if !contains(pathSegments, "admin") {
			next.ServeHTTP(w, r)
			return
		}
		
		// The second segment should be the resource (like "rewards").
		resource := pathSegments[3]

		// Map HTTP methods to permission actions.
		actionMap := map[string]string{
			"GET":    "read",
			"POST":   "create",
			"PUT":    "edit",
			"PATCH":  "edit",
			"DELETE": "delete",
		}
		// Get the action for this HTTP method.
		action, exists := actionMap[r.Method]
		if !exists {
			app.methodNotAllowed(w, r)
			return
		}

		// Construct the required permission.
		requiredPermission := "permissions:" + resource + "-" + action

		// Check if the user has this permission.
		hasPermission := false
		
		for _, permission := range user.Permissions {

			if permission.Name == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			app.notPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)

	}
	return http.HandlerFunc(fn)
}

func contains(s []string, e string) bool {
	for _, v := range s {
			if v == e {
					return true
			}
	}
	return false
}