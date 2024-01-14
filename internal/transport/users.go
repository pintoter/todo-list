package transport

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pintoter/todo-list/internal/entity"
)

// @Summary Sign Up
// @Description Sign Up
// @Tags auth
// @Accept json
// @Produce json
// @Param input body signUpInput true "input"
// @Success 200 {object} successCUDResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var input signUpInput
	if err := input.Set(r); err != nil {
		renderJSON(w, r, http.StatusBadRequest, errorResponse{err.Error()})
		return
	}

	_, err := h.service.SignUp(r.Context(), input.Email, input.Login, input.Password)
	if err != nil {
		renderJSON(w, r, http.StatusInternalServerError, errorResponse{
			Err: err.Error(),
		})
		return
	}

	renderJSON(w, r, http.StatusCreated, successCUDResponse{Message: "user successfully registered"})
}

// @Summary Sign In
// @Description Sign In
// @Tags auth
// @Accept json
// @Produce json
// @Param input body signInInput true "input"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var input signInInput
	if err := input.Set(r); err != nil {
		renderJSON(w, r, http.StatusBadRequest, errorResponse{err.Error()})
		return
	}

	tokens, err := h.service.SignIn(r.Context(), input.Login, input.Password)
	if err != nil {
		renderJSON(w, r, http.StatusInternalServerError, errorResponse{Err: err.Error()})
		return
	}

	r.Header.Set("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", tokens.RefreshToken))
	renderJSON(w, r, http.StatusOK, tokenResponse{Token: tokens.AccessToken})
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("refresh-token")

	if refreshToken == "" {
		renderJSON(w, r, http.StatusBadRequest, errorResponse{Err: "refresh-token not found"})
		return
	}

	tokens, err := h.service.RefreshTokens(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, entity.ErrSessionDoesntExist) {
			renderJSON(w, r, http.StatusBadRequest, errorResponse{Err: err.Error()})
		} else {
			renderJSON(w, r, http.StatusInternalServerError, errorResponse{Err: err.Error()})
		}
		return
	}

	r.Header.Set("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", tokens.RefreshToken))
	renderJSON(w, r, http.StatusOK, tokenResponse{Token: tokens.AccessToken})
}
