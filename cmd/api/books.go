package main

import (
	"errors"
	"fmt"
	"net/http"
	"onlinelibrary.beks.net/internal/data"
	"onlinelibrary.beks.net/internal/validator"
	"time"
)

// Add a createBookHandler for the "POST /v1/Books" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		CreatedAt time.Time     `json:"-"` // Use the - directive
		Author    string        `json:"author"`
		Title     string        `json:"title"`
		Year      int32         `json:"year,omitempty"`
		Readtime  data.Readtime `json:"readtime"` // Add the string directive
		Genres    []string      `json:"genres,omitempty"`
		PageCount int32         `json:"pagecount,omitempty"`
		Rating    float32       `json:"rating,omitempty"`
		Language  []string      `json:"language,omitempty"`
		Version   int32         `json:"version"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book := &data.Book{
		Author:    input.Author,
		Title:     input.Title,
		Year:      input.Year,
		Readtime:  input.Readtime,
		Genres:    input.Genres,
		PageCount: input.PageCount,
		Rating:    input.Rating,
		Languages: input.Language,
		Version:   input.Version,
	}

	// Initialize a new Validator.
	v := validator.New()
	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Call the Insert() method on our movies model, passing in a pointer to the
	// validated movie struct. This will create a record in the database and update the
	// movie struct with the system-generated information.
	err = app.models.Books.Insert(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/book/%d", book.ID))

	// Write a JSON response with a 201 Created status code, the movie data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Add a showBookHandler for the "GET /v1/Books/:id" endpoint. For now, we retrieve // the interpolated "id" parameter from the current URL and include it in a placeholder // response.
func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
	app.notFoundResponse(w,r)
	return
}	
book := data.Book{
	ID:   id,
	CreatedAt: time.Now(), 
	Title: "Skazka ",
	Readtime: 102,
	Genres: []string{"drama","romance","war"},
	Version: 1,
	Author: "Pushkin",
	Year: 1985,
	Language: []string{"Rus","Eng"},
	PageCount: 130,
	Rating: 8.5,


}

err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil) 
if err != nil {
	app.serverErrorResponse(w,r,err)
 }
// Otherwise, interpolate the Book ID in a placeholder response.

}