package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pintoter/todo-list/internal/entity"
)

const (
	format = "2006-01-02"
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

	date, err := time.Parse(format, input.Date)
	if err != nil {
		newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
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

type getNoteResponse struct {
	Note entity.Note `json:"note"`
}

func (h *Handler) getNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	note, err := h.service.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrNoteNotExits) {
			newErrorResponse(w, r, http.StatusNotFound, entity.ErrNoteNotExits.Error())
		} else {
			newErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(w, r, http.StatusOK, getNoteResponse{Note: note})
}

type getNotesRequest struct {
	Status string `json:"status,omitempty"`
	Date   string `json:"date,omitempty"`
	Limit  int    `json:"limit,omitempty"  binding:"required,min=1"`
	Offset int    `json:"offset,omitempty" binding:"required,min=1"`
}

type getNotesResponse struct {
	Notes []entity.Note `json:"notes"`
}

func (h *Handler) getNotesHandler(w http.ResponseWriter, r *http.Request) {
	var input getNotesRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		newErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var page int = 1
	queryPage := r.URL.Query().Get("page")
	if queryPage != "" {
		page, err = strconv.Atoi(queryPage)
		if err != nil || page < 0 {
			newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidPage.Error())
			return
		}
	}

	if input.Limit <= 0 || input.Offset < 0 {
		newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	notes, err := h.service.GetNotes(context.Background(), input.Limit, (page-1)*input.Offset, input.Status, input.Date)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidDate) {
			newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidDate.Error())
		} else if errors.Is(err, entity.ErrInvalidStatus) {
			newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidStatus.Error())
		} else if errors.Is(err, entity.ErrNoteExists) {
			newErrorResponse(w, r, http.StatusNotFound, entity.ErrNoteExists.Error())
		} else {
			newErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(w, r, http.StatusOK, getNotesResponse{Notes: notes})
}

type updateNoteInput struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

type updateNoteResponse struct {
	Message string `json:"message"`
}

func (h *Handler) updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var input updateNoteInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || (input.Title == "" && input.Description == "" && input.Status == "") {
		newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	err = h.service.UpdateNote(ctx, id, input.Title, input.Description, input.Status)
	if err != nil {
		newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidInput.Error())
	}

	newResponse(w, r, http.StatusAccepted, updateNoteResponse{Message: "Successfully update"})
}

type deleteNoteResponse struct {
	Message string `json:"message"`
}

func (h *Handler) deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err := h.service.DeleteById(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrNoteExists) {
			newErrorResponse(w, r, http.StatusBadRequest, entity.ErrNoteExists.Error())
			return
		} else {
			newErrorResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}

	newResponse(w, r, http.StatusOK, deleteNoteResponse{Message: "Succesfully delete note"})
}

type deleteNotesResponse struct {
	Message string `json:"message"`
}

func (h *Handler) deleteNotesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	err := h.service.DeleteNotes(ctx)
	if err != nil {
		newErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(w, r, http.StatusOK, deleteNotesResponse{Message: "Succesfully delete all notes"})
}
