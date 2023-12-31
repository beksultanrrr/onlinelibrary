package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	// Register the relevant methods, URL patterns and handler functions for our // endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST" // respectively.
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/book", app.createBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/book/:id", app.showBookHandler)

	router.HandlerFunc(http.MethodPut, "/v1/book/:id", app.updateBookHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/book/:id", app.deleteBookHandler)
	// Return the httprouter instance.
	return router
}
