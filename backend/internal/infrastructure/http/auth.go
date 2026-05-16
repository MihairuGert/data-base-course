package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"autopark/internal/application"
)

type authContextKey struct{}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	session, err := h.auth.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	user, ok := currentUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user":        user,
		"permissions": user.Permissions,
	})
}

func (h *Handler) authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || strings.TrimSpace(token) == "" {
			writeError(w, http.StatusUnauthorized, errors.New("missing bearer token"))
			return
		}
		user, err := h.auth.Parse(strings.TrimSpace(token))
		if err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return
		}
		ctx := context.WithValue(r.Context(), authContextKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func currentUser(r *http.Request) (application.AuthUser, bool) {
	user, ok := r.Context().Value(authContextKey{}).(application.AuthUser)
	return user, ok
}

func (h *Handler) requireEntity(w http.ResponseWriter, r *http.Request, entity, action string) bool {
	user, ok := currentUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return false
	}
	if !application.CanEntity(user.Role, entity, action) {
		writeError(w, http.StatusForbidden, errors.New("forbidden"))
		return false
	}
	return true
}

func (h *Handler) requireReport(w http.ResponseWriter, r *http.Request, report string) bool {
	user, ok := currentUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return false
	}
	if !application.CanReport(user.Role, report) {
		writeError(w, http.StatusForbidden, errors.New("forbidden"))
		return false
	}
	return true
}
