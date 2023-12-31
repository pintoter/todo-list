package transport

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pintoter/todo-list/internal/service"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Handler struct {
	router  *mux.Router
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	handler := &Handler{
		router:  mux.NewRouter(),
		service: service,
	}

	handler.InitRoutes()

	return handler
}

func (h *Handler) InitRoutes() {
	h.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	auth := h.router.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodPost)
		
	}

	v1 := h.router.PathPrefix("/api/v1").Subrouter()
	{
		v1.HandleFunc("/note", h.createNote).Methods(http.MethodPost)
		v1.HandleFunc("/note/{id:[0-9]+}", h.getNote).Methods(http.MethodGet)
		v1.HandleFunc("/note/{id:[0-9]+}", h.updateNote).Methods(http.MethodPatch)
		v1.HandleFunc("/note/{id:[0-9]+}", h.deleteNote).Methods(http.MethodDelete)
		v1.HandleFunc("/notes", h.getNotes).Methods(http.MethodGet)
		v1.HandleFunc("/notes", h.deleteNotes).Methods(http.MethodDelete)
		v1.HandleFunc("/notes/{page:[0-9]+}", h.getNotesExtended).Methods(http.MethodPost)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Request] [%s] %s - [FROM]: %s", r.Method, r.URL, r.RemoteAddr)
	h.router.ServeHTTP(w, r)
}
