package main

import (
	"net/http"
)

// JSONPayload is the type for json posted to this api
type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// WriteLog accepts json post requests and writes to mongo
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload

	_ = app.readJSON(w, r, &requestPayload)

	err := app.logEvent(requestPayload.Name, requestPayload.Data)
	if err != nil {
		_ = app.errorJSON(w, r, err, http.StatusBadRequest)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged event " + requestPayload.Name,
	}

	_ = app.writeJSON(w, r, http.StatusAccepted, resp)
}
