package logic

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	db "proyecto/db/sqlc" // Carpeta donde sqlc guarda el código generado.
	"strconv"
	"strings"
)

type MovieGenreRequest struct {
	GenreID int32 `json:"genre_id"`
}

func ValidateMovieGenre(mg MovieGenreRequest) error {
	if mg.GenreID <= 0 {
		return errors.New("El id de género debe ser mayor a 0")
	}
	return nil
}

// MovieGenresHandler maneja /movies/{id}/genres y /movies/{id}/genres/{genre_id}
func MovieGenresHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	// Caso POST: /movies/{movie_id}/genres (asociar género)
	if len(parts) == 4 {
		movieID, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "ID de película inválido", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPost {
			addMovieGenre(w, r, int32(movieID))
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	// Caso DELETE: /movies/{movie_id}/genres/{genre_id} (desasociar género)
	if len(parts) == 5 {
		movieID, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "ID de película inválido", http.StatusBadRequest)
			return
		}

		genreID, err := strconv.Atoi(parts[4])
		if err != nil {
			http.Error(w, "ID de género inválido", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodDelete {
			removeMovieGenre(w, r, int32(movieID), int32(genreID))
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
		return
	}

	http.Error(w, "Ruta no encontrada", http.StatusNotFound)
}

func addMovieGenre(w http.ResponseWriter, r *http.Request, movieID int32) {
	// 1. Verificar si la película existe
	_, err := Queries.GetMovie(r.Context(), movieID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Película no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener la película", http.StatusInternalServerError)
		return
	}

	// 2. Decodificar el JSON
	var req MovieGenreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if errVal := ValidateMovieGenre(req); errVal != nil {
		http.Error(w, errVal.Error(), http.StatusBadRequest)
		return
	}

	// 3. Verificar si el género existe
	_, err = Queries.GetGenre(r.Context(), req.GenreID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Género no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener el género", http.StatusInternalServerError)
		return
	}

	// 4. Insertar la relación en movie_genre
	params := db.AddMovieGenreParams{
		MovieID: movieID,
		GenreID: req.GenreID,
	}

	result, err := Queries.AddMovieGenre(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al asociar el género a la película (posible duplicado)", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // HTTP 201
	json.NewEncoder(w).Encode(result)
}

func removeMovieGenre(w http.ResponseWriter, r *http.Request, movieID int32, genreID int32) {
	params := db.RemoveMovieGenreParams{
		MovieID: movieID,
		GenreID: genreID,
	}

	if err := Queries.RemoveMovieGenre(r.Context(), params); err != nil {
		http.Error(w, "Error al eliminar la asociación del género", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content
}
