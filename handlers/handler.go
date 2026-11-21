package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	db "proyecto/db/sqlc"
	"proyecto/views"
)

type Handler struct {
	queries *db.Queries
}

func New(q *db.Queries) *Handler {
	return &Handler{queries: q}
}

func nullStr(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func (h *Handler) loadGenresByMovie(ctx context.Context) (map[int32][]string, error) {
	rows, err := h.queries.ListMovieGenres(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[int32][]string)
	for _, row := range rows {
		result[row.MovieID] = append(result[row.MovieID], row.NameGenre)
	}
	return result, nil
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	movies, err := h.queries.ListMovie(r.Context())
	if err != nil {
		log.Printf("error al listar películas: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}
	genres, err := h.queries.ListGenre(r.Context())
	if err != nil {
		log.Printf("error al listar géneros: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}
	byMovie, err := h.loadGenresByMovie(r.Context())
	if err != nil {
		log.Printf("error al listar géneros de las películas: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := views.IndexPage("CINE+ | Catálogo de Películas", movies, genres, byMovie).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar: %v", err)
	}
}
