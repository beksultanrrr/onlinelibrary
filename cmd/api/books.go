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
		app.notFoundResponse(w, r)
		return
	}
	// Call the Get() method to fetch the data for a specific movie. We also need to
	// use the errors.Is() function to check if it returns a data.ErrRecordNotFound
	// error, in which case we send a 404 Not Found response to the client.
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Fetch the existing movie record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Author    string        `json:"author"`
		Title     string        `json:"title"`
		Year      int32         `json:"year,omitempty"`
		Readtime  data.Readtime `json:"readtime"` // Add the string directive
		Genres    []string      `json:"genres,omitempty"`
		PageCount int32         `json:"pagecount,omitempty"`
		Rating    float32       `json:"rating,omitempty"`
		Languages []string      `json:"language,omitempty"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book.Author = input.Author
	book.Title = input.Title
	book.Year = input.Year
	book.Readtime = input.Readtime
	book.Genres = input.Genres
	book.PageCount = input.PageCount
	book.Rating = input.Rating
	book.Languages = input.Languages

	// Validate the updated movie record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	if data.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Pass the updated movie record to our new Update() method.
	err = app.models.Books.Update(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Write the updated movie record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the movie from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = app.models.Books.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "book successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
