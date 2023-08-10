package main

import (
	"net/http"

	"github.com/ton-developer-program/internal/response"
	"github.com/tonkeeper/tongo/liteapi"
)

var networks = map[string]*liteapi.Client{}


func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}
	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
	}
}

func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, "protected")
}
