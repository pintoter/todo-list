package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pintoter/todo-list/internal/entity"
)

type createNoteResponse struct {
	ID int `json:"id"`
}

type getNoteResponse struct {
	Note entity.Note `json:"note"`
}

type getNotesResponse struct {
	Notes []entity.Note `json:"notes"`
}

type putResponse struct {
	Message string `json:"message"`
}

func newResponse(w http.ResponseWriter, r *http.Request, code int, data any) {
	log.Printf("[%s] %s - Response: Success - Status code: [%d]", r.Method, r.URL.Path, code)

	resp, _ := json.MarshalIndent(data, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

type errorResponse struct {
	Err string `json:"error"`
}

func newErrorResponse(w http.ResponseWriter, r *http.Request, code int, err string) {
	log.Printf("[%s] %s - Response: Error - Status code: [%d] %s", r.Method, r.URL.Path, code, err)

	resp, _ := json.MarshalIndent(errorResponse{Err: err}, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
