package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) entities(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.crud.Entities())
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	entity := chi.URLParam(r, "entity")
	if !h.requireEntity(w, r, entity, "read") {
		return
	}
	rows, err := h.crud.List(r.Context(), entity)
	if err != nil {
		h.crudError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	entity := chi.URLParam(r, "entity")
	if !h.requireEntity(w, r, entity, "read") {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	row, err := h.crud.Get(r.Context(), entity, id)
	if err != nil {
		h.crudError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, row)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	entity := chi.URLParam(r, "entity")
	if !h.requireEntity(w, r, entity, "create") {
		return
	}
	body, ok := decodeBody(w, r)
	if !ok {
		return
	}
	row, err := h.crud.Create(r.Context(), entity, body)
	if err != nil {
		h.crudError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, row)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	entity := chi.URLParam(r, "entity")
	if !h.requireEntity(w, r, entity, "update") {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	body, ok := decodeBody(w, r)
	if !ok {
		return
	}
	row, err := h.crud.Update(r.Context(), entity, id, body)
	if err != nil {
		h.crudError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, row)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	entity := chi.URLParam(r, "entity")
	if !h.requireEntity(w, r, entity, "delete") {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := h.crud.Delete(r.Context(), entity, id); err != nil {
		h.crudError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func decodeBody(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	defer r.Body.Close()
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return nil, false
	}
	return body, true
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return 0, false
	}
	return id, true
}

func (h *Handler) crudError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		writeError(w, http.StatusNotFound, err)
	case strings.Contains(err.Error(), "unknown entity"):
		writeError(w, http.StatusNotFound, err)
	default:
		writeError(w, http.StatusBadRequest, err)
	}
}
