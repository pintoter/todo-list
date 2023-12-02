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
	dateFormat = "2006-01-02"
)

type createNoteInput struct {
	Title       string `json:"title" binding:"required,min=1,max=80"`
	Description string `json:"description,omitempty"`
	Date        string `json:"date,omitempty" binding:"min=9,max=10"`
	Status      string `json:"status,omitempty"`
}

func (i createNoteInput) Validate() error {
	if i.Date != "" {
		_, err := time.Parse(dateFormat, i.Date)
		if err != nil {
			return err
		}
	}

	if i.Title == "" {
		return entity.ErrInvalidInput
	}
	return nil
}

// @Summary Create note
// @Description create note
// @Tags note
// @Accept json
// @Produce json
// @Param input body createNoteInput true "note info"
// @Success 201 {object} createNoteResponse
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/note [post]
func (h *Handler) createNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background() // change go r.Context()

	var input createNoteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	if err := input.Validate(); err != nil {
		newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	note := entity.Note{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
	}

	if input.Date != "" {
		note.Date, _ = time.Parse(dateFormat, input.Date)
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

	newResponse(w, r, http.StatusCreated, createNoteResponse{ID: id})
}

// @Summary Get note by id
// @Description Get note by id
// @Tags note
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} createNoteResponse
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/note/{id} [get]
func (h *Handler) getNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if id == 0 {
		newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidId.Error())
		return
	}

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

// @Summary Get all notes 
// @Description Get all notes
// @Tags note
// @Accept json
// @Produce json
// @Param id path int false "id"
// @Param input body getNotesRequest true "search params"
// @Success 201 {object} getNotesRequest
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/get_all_extended [post]
func (h *Handler) getNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := h.service.GetNotes(context.Background())
	if err != nil {
		newErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(w, r, http.StatusOK, getNotesResponse{Notes: notes})
}

type getNotesRequest struct {
	Page          int       `json:"-"`
	Status        string    `json:"status,omitempty"`
	Date          string    `json:"date,omitempty"`
	DateFormatted time.Time `json:"-"`
	Limit         int       `json:"limit" binding:"required"`
	Offset        int       `json:"offset" binding:"required"`
}

func (n *getNotesRequest) Set(r *http.Request) error {
	var err error
	n.Page, _ = strconv.Atoi(mux.Vars(r)["page"])
	if n.Page == 0 {
		return entity.ErrInvalidPage
	}

	err = json.NewDecoder(r.Body).Decode(n)
	if err != nil {
		return entity.ErrInvalidInput
	}

	if n.Limit <= 0 || n.Offset < 0 {
		return entity.ErrInvalidInput
	}

	if n.Date != "" {
		n.DateFormatted, err = time.Parse(dateFormat, n.Date)
		if err != nil {
			return entity.ErrInvalidDate
		}
	}
	return nil
}

// @Summary Get notes with filter
// @Description Get notes with filter
// @Tags note
// @Accept json
// @Produce json
// @Param id path int false "id"
// @Param input body getNotesRequest true "search params"
// @Success 201 {object} notesRequest
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/get_all_extended [post]
func (h *Handler) getNotesExtendedHandler(w http.ResponseWriter, r *http.Request) {
	var input getNotesRequest

	if err := input.Set(r); err != nil {
		newErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	notes, err := h.service.GetNotesExtended(context.Background(), input.Limit, (input.Page-1)*input.Offset, input.Status, input.DateFormatted)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidStatus) {
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

// @Summary Update note
// @Description update note by id
// @Tags note
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param input body updateNoteInput true "params for update"
// @Success 201 {object} successResponse
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/note/{id} [patch]
func (h *Handler) updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if id == 0 {
		newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidId.Error())
		return
	}

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

	newResponse(w, r, http.StatusAccepted, putResponse{Message: "Successfully update"})
}

// @Summary Delete note
// @Description Delete note by id
// @Tags note
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} successResponse
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/note/{id} [delete]
func (h *Handler) deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if id == 0 {
		newErrorResponse(w, r, http.StatusBadRequest, entity.ErrInvalidId.Error())
		return
	}

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

	newResponse(w, r, http.StatusOK, putResponse{Message: "Succesfully delete note"})
}
