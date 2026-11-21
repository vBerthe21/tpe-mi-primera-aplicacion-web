package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"proyecto/logic"
	"proyecto/views"
)

// POST /genres
func (h *Handler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulario inválido", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))

	genres, err := h.queries.ListGenre(r.Context())
	if err != nil {
		log.Printf("error al listar géneros: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}

	if err := logic.ValidateGenre(name); err != nil {
		views.GenreList(genres, err.Error()).Render(r.Context(), w)
		return
	}
	for _, g := range genres {
		if strings.EqualFold(g.NameGenre, name) {
			views.GenreList(genres, "el género ya existe").Render(r.Context(), w)
			return
		}
	}

	if _, err := h.queries.CreateGenre(r.Context(), name); err != nil {
		log.Printf("error al crear género: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}

	genres, err = h.queries.ListGenre(r.Context())
	if err != nil {
		log.Printf("error al listar géneros: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}
	views.GenreList(genres, "").Render(r.Context(), w)
	views.GenreCheckboxesOOB(genres).Render(r.Context(), w)
}

// DELETE /genres/{id}
func (h *Handler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.queries.DeleteGenre(r.Context(), int32(id)); err != nil {
		log.Printf("error al eliminar género: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}

	genres, err := h.queries.ListGenre(r.Context())
	if err != nil {
		log.Printf("error al listar géneros: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}
	views.GenreCheckboxesOOB(genres).Render(r.Context(), w)
}
