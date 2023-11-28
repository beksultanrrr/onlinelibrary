package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"onlinelibrary.beks.net/internal/validator"
	"time"
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

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type BookModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m BookModel) Insert(book *Book) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
		INSERT INTO books (author, title, year, readtime, genres, pagecount, rating, language)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, version`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []interface{}{book.Author, book.Title, book.Year, book.Readtime, pq.Array(book.Genres), book.PageCount, book.Rating, pq.Array(book.Languages)}
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m BookModel) Get(id int64) (*Book, error) {
	// The PostgreSQL bigserial type that we're using for the movie ID starts
	// auto-incrementing at 1 by default, so we know that no movies will have ID values
	// less than that. To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the movie data.
	query := `
		SELECT id, created_at, author, title, year, readtime, genres, pagecount, rating, language, version
		FROM books
		WHERE id = $1`
	var book Book

	// Use the context.WithTimeout() function to create a context.Context which carries a
	// 3-second timeout deadline. Note that we're using the empty context.Background()
	// as the 'parent' context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns
	defer cancel()

	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Author,
		&book.Title,
		&book.Year,
		&book.Readtime,
		pq.Array(&book.Genres),
		&book.PageCount,
		&book.Rating,
		pq.Array(&book.Languages),
		&book.Version,
	)

	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &book, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m BookModel) Update(book *Book) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
		UPDATE books
		SET author = $1, title = $2, year = $3, readtime = $4, genres = $5, pagecount = $6, rating = $7, language = $8, version = version + 1
		WHERE id = $9 AND version = $10
		RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		book.Author,
		book.Title,
		book.Year,
		book.Readtime,
		pq.Array(book.Genres),
		book.PageCount,
		book.Rating,
		pq.Array(book.Languages),
		book.ID,
		book.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the SQL query. If no matching row could be found, we know the movie
	// version has changed (or the record has been deleted) and we return our custom
	// ErrEditConflict error.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&book.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m BookModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM books
		WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m BookModel) GetAll(title string, genres []string, filters Filters) ([]*Book, Metadata, error) {
	// Construct the SQL query to retrieve all movie records.
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, author, title, year, readtime, genres, pagecount, rating, language, version
		FROM books
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, pq.Array(genres), filters.limit(), filters.offset()}
	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	totalRecords := 0
	// Initialize an empty slice to hold the movie data.
	books := []*Book{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var book Book
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords,
			&book.ID,
			&book.CreatedAt,
			&book.Author,
			&book.Title,
			&book.Year,
			&book.Readtime,
			pq.Array(book.Genres),
			&book.PageCount,
			&book.Rating,
			pq.Array(book.Languages),
			&book.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the Movie struct to the slice.
		books = append(books, &book)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of movies.
	return books, metadata, nil
}
