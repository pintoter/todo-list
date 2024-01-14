package transport

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pintoter/todo-list/internal/config"
	"github.com/pintoter/todo-list/internal/service"
	"github.com/pintoter/todo-list/pkg/auth"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Config interface {
	GetMode() string
}

type Handler struct {
	router       *mux.Router
	service      *service.Service
	tokenManager auth.TokenManager
}

func NewHandler(service *service.Service, cfg Config) *Handler {
	handler := &Handler{
		router:  mux.NewRouter(),
		service: service,
	}

	if cfg.GetMode() != config.Production {
		handler.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		)).Methods(http.MethodGet)
	}

	handler.InitRoutes()

	return handler
}

func (h *Handler) InitRoutes() {
	auth := h.router.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodPost)
		auth.HandleFunc("/refresh", h.refresh).Methods(http.MethodPost)
	}

	v1 := h.router.PathPrefix("/api/v1").Subrouter()
	{
		v1.Use(h.authMiddleware)
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
