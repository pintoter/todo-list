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

	v1 := h.router.PathPrefix("/api/v1").Subrouter()
	{
		v1.HandleFunc("/note", h.createNoteHandler).Methods(http.MethodPost)
		v1.HandleFunc("/note/{id:[0-9]+}", h.updateNoteHandler).Methods(http.MethodPatch)
		v1.HandleFunc("/note/{id:[0-9]+}", h.getNoteHandler).Methods(http.MethodGet)
		v1.HandleFunc("/note/{id:[0-9]+}", h.deleteNoteHandler).Methods(http.MethodDelete)
		v1.HandleFunc("/notes", h.getNotesHandler).Methods(http.MethodGet)
		v1.HandleFunc("/notes", h.deleteNotesHandler).Methods(http.MethodDelete)
		v1.HandleFunc("/notes/{page:[0-9]+}", h.getNotesExtendedHandler).Methods(http.MethodPost)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s", r.Method, r.URL)
	h.router.ServeHTTP(w, r)
}
