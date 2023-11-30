package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/pintoter/todo-list/internal/entity"
)

const (
	layout = "2006-01-02"
)

type createNoteInput struct {
	Title       string `json:"title" binding:"required,min=1,max=80"`
	Description string `json:"description,omitempty"`
	Date        string `json:"date,omitempty" binding:"required,min=10,max=80"`
	Status      string `json:"status,omitempty"`
}

type createNoteResponse struct {
	ID int `json:"id"`
}

func (h *Handler) createNoteHandler(w http.ResponseWriter, r *http.Request) {
	var input createNoteInput
	ctx := context.Background()

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	date, err := time.Parse(layout, input.Date)
	if err != nil {
		newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidInput.Error())
	}

	note := entity.Note{
		Title:       input.Title,
		Description: input.Description,
		Date:        date,
		Status:      input.Status,
	}

	id, err := h.service.CreateNote(ctx, note)
	if err != nil {
		if errors.Is(err, entity.ErrNoteExists) {
			newErrorResponse(w, r, http.StatusConflict, err.Error())
		} else {
			newErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(w, r, http.StatusCreated, id)
}

// type getNoteResponse struct {
// 	Note entity.Note `json:"note"`
// }

// func (h *Handler) getNoteHandler(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	id := params["id"]

// 	noteId, _ := strconv.Atoi(id)

// 	note, err := h.service.GetById(context.Background(), noteId)
// 	if err != nil {
// 		if errors.Is(err, entity.ErrNoteNotFound) {
// 			http.Error(w, entity.ErrArticleNotFound.Error(), http.StatusNotFound)
// 		} else {
// 			http.Error(w, http.StatusInternalServerError, http.StatusNotFound)
// 		}
// 		return
// 	}

// 	resp, _ := json.Marshal(createNoteResponse{
// 		Note: note,
// 	})

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(resp)
// }

// type getNoteResponse struct {
// 	Note []entity.Note `json:"note"`
// }

// func (h *Handler) getNoteHandler(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	id := params["id"]

// 	noteId, _ := strconv.Atoi(id)

// 	note, err := h.service.GetNote(context.Background(), noteId)
// 	if err != nil {
// 		if errors.Is(err, entity.ErrNoteNotFound) {
// 			http.Error(w, entity.ErrArticleNotFound.Error(), http.StatusNotFound)
// 		} else {
// 			http.Error(w, http.StatusInternalServerError, http.StatusNotFound)
// 		}
// 		return
// 	}

// 	resp, _ := json.Marshal(getNoteResponse{
// 		note: note,
// 	})

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(resp)
// }

// type getNotesResponse struct {
// 	Notes []entity.Note `json:"notes"`
// }

// func (h *Handler) getNotesHandler(w http.ResponseWriter, r *http.Request) {
// 	var limit, offest int = 3, 0
// 	var curStatus = "not done" // ??
// 	var t time.Time

// 	page := r.URL.Query().Get("page")
// 	if page != "" {
// 		limit, err = strconv.Atoi(page)
// 		if err != nil || limit <= 0 {
// 			http.Error(w, entity.ErrInvalidId.Error(), http.StatusBadRequest)
// 			return
// 		}
// 	}

// 	count := r.URL.Query().Get("count")
// 	if count != "" {
// 		offset, err = strconv.Atoi(count)
// 		if err != nil || offset < 0 {
// 			http.Error(w, entity.ErrInvalidId.Error(), http.StatusBadRequest)
// 			return
// 		}
// 	}

// 	queryStatus := r.URL.Query().Get("status")
// 	if queryStatus != "" && (queryStatus != statusDone || queryStatus != statusNotDone) {
// 		http.Error(w, entity.ErrInvalidStatus.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	queryDate := r.URL.Query().Get("date")
// 	if queryDate != "" {
// 		t, err = time.Parse(layoutISO, queryDate)
// 		if err != nil {
// 			http.Error(w, entity.ErrInvalidDate.Error(), http.StatusBadRequest)
// 			return
// 		}
// 	}

// 	if queryStatus == "" && queryOrder != "" {
// 		http.Error(w, entity.ErrInvalidFilter.Error(), http.StatusBadRequest)
// 		return
// 	} else {
// 		curOrder = queryOrder
// 	}

// 	var notes []entity.Note

// 	switch {
// 	case queryStatus != "" && queryOrder != "":
// 		note, err = h.service.GetByDateAndStatus(context.Background(), "ASC", "DONE", limit, offset) // asc или desc
// 	case queryStatus != "":
// 		note, err = h.service.GetByDateAndStatus(context.Background(), "ASC", limit, offset) // asc или desc
// 	default:
// 		note, err = h.service.GetAll(context.Background(), limit, offset)
// 	}

// 	if err != nil {
// 		if errors.Is(err, entity.ErrNoteNotFound) {
// 			http.Error(w, entity.ErrArticleNotFound.Error(), http.StatusNotFound)
// 		} else {
// 			http.Error(w, http.StatusInternalServerError, http.StatusNotFound)
// 		}
// 		return
// 	}

// 	resp, _ := json.Marshal(getNotesResponse{
// 		notes: notes,
// 	})

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(resp)
// }

// type updateNoteStatusInput struct {
// 	Status string `json:"status"`
// }

// func (h *Handler) updateNoteStatusHandler(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	id := params["id"]

// 	noteId, _ := strconv.Atoi(id)

// 	var input createNoteInput

// 	err := json.NewDecoder(r.Body).Decode(&input)
// 	if err != nil {
// 		http.Error(w, entity.ErrInvalidId.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	err := h.service.UpdateStatus(ctx, noteId, input.Status)
// 	if err != nil {
// 		if errors.Is(err, entity.ErrNoteNotFound) {
// 			http.Error(w, entity.ErrArticleNotFound.Error(), http.StatusNotFound)
// 		} else {
// 			http.Error(w, http.StatusInternalServerError, http.StatusNotFound)
// 		}
// 		return
// 	}

// 	resp, _ := json.Marshal(getNoteResponse{
// 		note: note,
// 	})

// 	// ПОДУМАТЬ ЧТО ОТДАВАТЬ
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusNoContent)
// 	w.Write(resp)
// }

// type updateNoteInfoInput struct {
// 	Title       string `json:"title,omitempty"`
// 	Description string `json:"description,omitempty"`
// }

// func (h *Handler) updateNoteInfoHandler(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	id := params["id"]

// 	noteId, _ := strconv.Atoi(id)

// 	var input updateNoteInfoInput

// 	err := json.NewDecoder(r.Body).Decode(&input)
// 	if err != nil {
// 		http.Error(w, entity.ErrInvalidId.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	err := h.service.UpdateInfo(ctx, noteId, input.Title, input.Description)
// 	if err != nil {
// 		if errors.Is(err, entity.ErrNoteNotFound) {
// 			http.Error(w, entity.ErrNoteNotFound.Error(), http.StatusNotFound)
// 		} else {
// 			http.Error(w, http.StatusInternalServerError, http.StatusNotFound)
// 		}
// 		return
// 	}

// 	resp, _ := json.Marshal(getNoteResponse{
// 		note: note,
// 	})

// 	// ПОДУМАТЬ ЧТО ОТДАВАТЬ
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(resp)
// }

// func (h *Handler) deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	id := params["id"]

// 	noteId, _ := strconv.Atoi(id)

// 	err := h.service.DeleteById(ctx, noteId)
// 	if err != nil {
// 		if errors.Is(err, entity.ErrNoteNotFound) {
// 			http.Error(w, entity.ErrNoteNotFound.Error(), http.StatusBadRequest)
// 		} else {
// 			http.Error(w, http.StatusInternalServerError, http.StatusNotFound)
// 		}
// 		return
// 	}

// 	// ПОДУМАТЬ ЧТО ОТДАВАТЬ
// 	w.WriteHeader(http.StatusNoContent)
// 	w.Write(nil) // ?
// }

// func (h *Handler) deleteNotesHandler(w http.ResponseWriter, r *http.Request) {
// 	err := h.service.DeleteNotes(ctx)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// ПОДУМАТЬ ЧТО ОТДАВАТЬ
// 	w.WriteHeader(http.StatusNoContent)
// 	w.Write(nil) // ?
// }
