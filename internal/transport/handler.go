package transport

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pintoter/todo-list/internal/service"
)

/*
link:
https://github.com/eliben/code-for-blog/blob/master/2021/go-rest-servers/gorilla/gorilla.go

 add to repository:
email, loginTime := "human@example.com", time.Now()
result, err := db.Exec("INSERT INTO UserAccount VALUES ($1, $2)", email, loginTime)
if err != nil {
  panic(err)
}

parse time:
const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)
date := "1999-12-31"
t, _ := time.Parse(layoutISO, date)
fmt.Println(t)                  // 1999-12-31 00:00:00 +0000 UTC
fmt.Println(t.Format(layoutUS)) // December 31, 1999
*/

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
	v1 := h.router.PathPrefix("/api/v1").Subrouter()
	{
		v1.HandleFunc("/note", h.createNoteHandler).Methods("POST")
		v1.HandleFunc("/note", h.getNotesHandler).Methods("GET")
		v1.HandleFunc("/note/{id:[0-9]+}/update", h.updateNoteHandler).Methods("PATCH")
		// h.router.HandleFunc("/note", h.deleteNotesHandler).Methods("DELETE")
		// h.router.HandleFunc("/note/{id:[0-9]+}", h.getNoteHandler).Methods("GET")
		// h.router.HandleFunc("/note/{id:[0-9]+}", h.deleteNoteHandler).Methods("DELETE")
		// h.router.HandleFunc("/note/{id:[0-9]+}/update_info", h.updateNoteInfoHandler).Methods("PATCH")
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s", r.Method, r.URL)

	h.router.ServeHTTP(w, r)
}
