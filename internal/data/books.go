package data

import (
	 "time"
	 "fmt"
	 "encoding/json"
	 "onlinelibrary.beks.net/internal/validator" 
)

type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // Use the - directive
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Readtime  Readtime  `json:"-"` // Add the string directive
	Genres    []string  `json:"genres,omitempty"`
	PageCount int32     `json:"pagecount,omitempty"`
	Rating    float32   `json:"rating,omitempty"`
	Languages []string  `json:"language,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(book.Title != "", "title", "must be provided")
	v.Check(len(book.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(book.Year != 0, "year", "must be provided")
	v.Check(book.Year >= 1888, "year", "must be greater than 1888")
	v.Check(book.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(book.Readtime != 0, "readtime", "must be provided")
	v.Check(book.Readtime > 0, "readtime", "must be a positive integer")
	v.Check(book.Genres != nil, "genres", "must be provided")
	v.Check(len(book.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(validator.Unique(book.Genres), "genres", "must not contain duplicate values")
	v.Check(book.Languages != nil, "language", "must be provided")
	v.Check(len(book.Languages) >= 1, "language", "must contain at least 1 language")
	v.Check(book.Rating <= 10, "rating", "must be less than 10")
	v.Check(book.Author != "", "author", "must be provided")
	v.Check(len(book.Author) <= 500, "author", "must not be more than 500 bytes long")
	v.Check(book.PageCount > 0, "pagecount", "must be less than 0")
}

func (m Book) MarshalJSON() ([]byte, error) {
	// Declare a variable to hold the custom readtime string (this will be the empty // string "" by default).
	var readtime string
	// If the value of the readtime field is not zero, set the readtime variable to be a // string in the format "<readtime> mins".
	if m.Readtime != 0 {
		readtime = fmt.Sprintf("%d mins", m.Readtime)
	}
	type BookAllias Book
	aux := struct {
		BookAllias
		Readtime string `json:"readtime,omitempty"`
	}{
		BookAllias: BookAllias(m),
		Readtime:   readtime,
	}

	return json.Marshal(aux)
}