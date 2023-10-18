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
	fmt.Fprintln(w, "create a new Book") 
}
// Add a showBookHandler for the "GET /v1/Books/:id" endpoint. For now, we retrieve // the interpolated "id" parameter from the current URL and include it in a placeholder // response.
func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
id, err := app.readIDParam(r)
	if err != nil {
	http.NotFound(w, r)
	return
}	
book := data.Book{
	ID:   id,
	CreatedAt: time.Now(), 
	Title: "Сказка о рыбаке и рыбке ",
	Runtime: 102,
	Genres: []string{"drama","romance","war"},
	Version: 1,
	Author: "Пушкин",
	Year: 1985,


}

err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil) 
if err != nil {
	app.logger.Println(err)
	http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError) }
// Otherwise, interpolate the Book ID in a placeholder response.

}