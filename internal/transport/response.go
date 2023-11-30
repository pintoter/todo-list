package transport

import (
	"encoding/json"
	"log"
	"net/http"
)

func newResponse(w http.ResponseWriter, r *http.Request, code int, data any) {
	log.Printf("[%s] %s - Response: Success", r.Method, r.URL.Path)

	resp, _ := json.MarshalIndent(data, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

type errorResponse struct {
	Err string `json:"error"`
}

func newErrorResponse(w http.ResponseWriter, r *http.Request, code int, err string) {
	log.Printf("[%s] %s - Response: Error", r.Method, r.URL.Path)

	resp, _ := json.MarshalIndent(errorResponse{Err: err}, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
