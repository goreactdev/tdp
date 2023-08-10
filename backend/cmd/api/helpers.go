package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ton-developer-program/internal/database"
)

func (app *application) newEmailData() map[string]any {
	data := map[string]any{
		"BaseURL": app.config.App.BaseUrl,
	}

	return data
}

func (app *application) backgroundTask(fn func() error) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				app.reportError(fmt.Errorf("%s", err))
			}
		}()

		err := fn()
		if err != nil {
			app.reportError(err)
		}
	}()
}


func basicAuth(handler http.Handler, username, password, realm string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized.\n"))
			return
		}
		handler.ServeHTTP(w, r)
	})
}


func getPagination(r *http.Request) (*database.Pagination, error) {
	startQuery := r.URL.Query().Get("_start")
	endQuery := r.URL.Query().Get("_end")

	// convert to int
	if startQuery == "" {
		return nil, errors.New("_start is required")
	}

	if endQuery == "" {
		return nil, errors.New("_end is required")
	}

	start, err := strconv.Atoi(startQuery)
	if err != nil {
		return  nil, errors.New("_start must be an integer")
	}

	end, err := strconv.Atoi(endQuery)
	if err != nil {
		return nil, errors.New("_end must be an integer")
	}

	return &database.Pagination{
		Start: start,
		End: end,
	}, nil

}