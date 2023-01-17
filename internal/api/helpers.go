package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/foroozf001/logger-service/internal/data"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		app.HttpReqs.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, r *http.Request, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		app.HttpReqs.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		app.HttpReqs.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method).Inc()
		return err
	}

	app.HttpReqs.WithLabelValues(strconv.Itoa(status), r.Method).Inc()
	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, r *http.Request, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	app.HttpReqs.WithLabelValues(strconv.Itoa(statusCode), r.Method).Inc()

	return app.writeJSON(w, r, statusCode, payload)
}

func (app *Config) logEvent(name, content string) error {
	event := data.LogItem{
		Name: name,
		Data: content,
	}

	return app.Models.LogItem.Insert(event)
}
