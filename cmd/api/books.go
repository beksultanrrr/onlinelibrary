package main
import ( 
	"fmt"
	"net/http"
	"time"
	"onlinelibrary.beks.net/internal/data"
)
// Add a createBookHandler for the "POST /v1/Books" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	
	var input struct {
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
	err := app.readJSON(w,r,&input)
	if err != nil {
		app.errorResponse(w,r,http.StatusBadRequest,err.Error())
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
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