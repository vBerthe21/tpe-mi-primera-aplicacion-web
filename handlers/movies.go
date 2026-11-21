package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	db "proyecto/db/sqlc"
	"proyecto/logic"
	"proyecto/views"
)

// Helper privado: lee y valida el formulario tanto al crear como al editar.
func parseAndValidateMovie(r *http.Request) (logic.MovieInput, error) {
	if err := r.ParseForm(); err != nil {
		return logic.MovieInput{}, errors.New("formulario inválido")
	}

	duracion, err := strconv.ParseInt(r.FormValue("duration"), 10, 32)
	if err != nil {
		return logic.MovieInput{}, errors.New("la duración debe ser un número entero")
	}

	in := logic.MovieInput{
		Name:        strings.TrimSpace(r.FormValue("name")),
		Author:      strings.TrimSpace(r.FormValue("author")),
		Duration:    int32(duracion),
		Score:       strings.TrimSpace(r.FormValue("score")),
		Review:      strings.TrimSpace(r.FormValue("review")),
		Description: strings.TrimSpace(r.FormValue("description")),
	}

	if err := logic.ValidateMovie(in); err != nil {
		return logic.MovieInput{}, err
	}
	return in, nil
}

// Helper privado: guarda los checkboxes de géneros tildados para una película.
func (h *Handler) attachGenres(r *http.Request, movieID int32) {
	for _, g := range r.Form["genres"] {
		genreID, err := strconv.Atoi(g)
		if err == nil {
			_, err = h.queries.AddMovieGenre(r.Context(), db.AddMovieGenreParams{
				MovieID: movieID,
				GenreID: int32(genreID),
			})
			if err != nil {
				log.Printf("error al asociar género %d a la película %d: %v", genreID, movieID, err)
			}
		}
	}
}

// Helper privado: vuelve a consultar películas + géneros y renderiza la lista.
func (h *Handler) renderMovieList(w http.ResponseWriter, r *http.Request) {
	movies, err := h.queries.ListMovie(r.Context())
	if err != nil {
		log.Printf("error al listar películas: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}
	byMovie, err := h.loadGenresByMovie(r.Context())
	if err != nil {
		log.Printf("error al listar géneros de las películas: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}
	if err := views.MovieList(movies, byMovie).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar la lista: %v", err)
		http.Error(w, "error interno al renderizar", http.StatusInternalServerError)
	}
}

// POST /movies
func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	in, err := parseAndValidateMovie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMovie, err := h.queries.CreateMovie(r.Context(), db.CreateMovieParams{
		Name:        in.Name,
		Author:      in.Author,
		Duration:    in.Duration,
		Score:       nullStr(in.Score),
		Review:      nullStr(in.Review),
		Description: nullStr(in.Description),
	})
	if err != nil {
		log.Printf("error al crear película: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}

	h.attachGenres(r, newMovie.ID)
	h.renderMovieList(w, r)
}

// DELETE /movies/{id}
func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	if err := h.queries.DeleteMovie(r.Context(), int32(id)); err != nil {
		log.Printf("error al eliminar película: %v", err)
		http.Error(w, "Error interno al eliminar", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GET /movies/{id}/edit
func (h *Handler) EditMovieForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	movie, err := h.queries.GetMovie(r.Context(), int32(id))
	if err == sql.ErrNoRows {
		http.Error(w, "película no encontrada", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("error al obtener película: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}

	genres, err := h.queries.ListGenre(r.Context())
	if err != nil {
		log.Printf("error al listar géneros: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}

	current, err := h.queries.ListGenresByMovie(r.Context(), int32(id))
	if err != nil {
		log.Printf("error al listar géneros de la película: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}
	selected := make(map[int32]bool)
	for _, g := range current {
		selected[g.ID] = true
	}

	views.MovieEditForm(movie, genres, selected).Render(r.Context(), w)
}

// PUT /movies/{id}
func (h *Handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	in, err := parseAndValidateMovie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.queries.UpdateMovie(r.Context(), db.UpdateMovieParams{
		ID:          int32(id),
		Name:        in.Name,
		Author:      in.Author,
		Duration:    in.Duration,
		Score:       nullStr(in.Score),
		Review:      nullStr(in.Review),
		Description: nullStr(in.Description),
	})
	if err != nil {
		log.Printf("error al actualizar película: %v", err)
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
	}

	if err := h.queries.RemoveAllMovieGenres(r.Context(), int32(id)); err != nil {
		log.Printf("error al quitar géneros de la película %d: %v", id, err)
	}
	h.attachGenres(r, int32(id))
	h.renderMovieList(w, r)
}
