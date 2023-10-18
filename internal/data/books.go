package data
import (
	 "time"
	 "fmt"
)


type Book struct {
	ID   int64 
	CreatedAt time.Time // Use the - directive
	Author string 
	Title string 
	Year int32 
	Runtime int32
	Genres []string 
	Version int32 

}

func (m Book) MarshalJSON() ([]byte, error) {
	// Declare a variable to hold the custom runtime string (this will be the empty // string "" by default).
	var runtime string
	// If the value of the Runtime field is not zero, set the runtime variable to be a // string in the format "<runtime> mins".
	if m.Runtime != 0 {
		runtime = fmt.Sprintf("%d mins", m.Runtime) 
	}
	// Create an anonymous struct to hold the data for JSON encoding. This has exactly // the same fields, types and tags as our Movie struct, except that the Runtime
	// field here is a string, instead of an int32. Also notice that we don't include // a CreatedAt field at all (there's no point including one, because we don't want // it to appear in the JSON output).
	aux := struct {
		ID   int64 
	CreatedAt time.Time `json:"-"` // Use the - directive
	Author string `json:"author"`
	Title string `json:"title"`
	Year int32 `json:"year,omitempty"`
	Runtime string `json:"runtime,omitempty"` // Add the string directive
	Genres []string `json:"genres,omitempty"` 
	Version int32  `json:"version"` 
	}{
	ID:  m.ID,
	CreatedAt: m.CreatedAt,
	Author: m.Author,
	Title: m.Title,
	Year: m.Year,
	Runtime: m.Runtime,// Add the string directive
	Genres: m.Genres,
	Version: m.Version,
	}
	// Set the values for the anonymous struct.

	// Encode the anonymous struct to JSON, and return it.
	return json.Marshal(aux) 
}