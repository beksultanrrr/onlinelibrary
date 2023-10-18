package data
import (
	 "time"
	 "fmt"
	 "encoding/json"
)


type Book struct {
	ID   int64 `json:"id"`
	CreatedAt time.Time `json:"-"` // Use the - directive
	Author string `json:"author"`
	Title string `json:"title"`
	Year int32 `json:"year,omitempty"`
	Runtime int32 `json:"-"`  // Add the string directive
	Genres []string `json:"genres,omitempty"` 
	Version int32  `json:"version"` 

}

func (m Book) MarshalJSON() ([]byte, error) {
	// Declare a variable to hold the custom runtime string (this will be the empty // string "" by default).
	var runtime string
	// If the value of the Runtime field is not zero, set the runtime variable to be a // string in the format "<runtime> mins".
	if m.Runtime != 0 {
		runtime = fmt.Sprintf("%d mins", m.Runtime) 
	}
	type BookAllias Book
	aux := struct {
		BookAllias
		Runtime string  `json:"runtime,omitempty"`
	}{
	 BookAllias: BookAllias(m),
	 Runtime: runtime,
	}
	
	return json.Marshal(aux) 
}