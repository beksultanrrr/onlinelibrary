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
	Readtime int32 `json:"-"`  // Add the string directive
	Genres []string `json:"genres,omitempty"` 
	PageCount int32 `json:"pagecount,omitempty"`
	Rating float32 `json:"rating,omitempty"`
	Language []string `json:"language,omitempty"`
	Version int32  `json:"version"` 
	
	


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
		Readtime string  `json:"readtime,omitempty"`
	}{
	 BookAllias: BookAllias(m),
	 Readtime: readtime,
	}
	
	return json.Marshal(aux) 
}