package main
import ( 
	"net/http"
)
// Declare a handler which writes a plain-text response with information about the // application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// js := `{"status": "available", "environment": %q, "version": %q}` 
	// js = fmt.Sprintf(js, app.config.env, version)
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(js))

	env := envelope{
		"status" : "availabe",
		"system_info" : map[string]string{
		"enviroment": app.config.env,
		"version": version,
		},
	}
	err := app.writeJSON(w, http.StatusOK,env, nil)


	if  err  != nil {
		app.logger.Println(err)
		http.Error(w,"The server encountered a problem and could not process your request", http.StatusInternalServerError)
	
	}
	
}