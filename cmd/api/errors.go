package main
import ( "fmt"
"net/http"
)
// The logError() method is a generic helper for logging an error message. Later in the // book we'll upgrade this to use structured logging, and record additional information // about the request including the HTTP method and URL.
func (app *application) logError(r *http.Request, err error) {
		app.logger.Println(err) 
	}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
	app.logError(r, err)
	w.WriteHeader(500) }
}
// The serverErrorResponse() method will be used when our application encounters an // unexpected problem at runtime. It logs the detailed error message, then uses the // errorResponse() helper to send a 500 Internal Server Error status code and JSON // response (containing a generic error message) to the client.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) { 
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message) 
}
// The notFoundResponse() method will be used to send a 404 Not Found status code and // JSON response to the client.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
 }

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message) 
}