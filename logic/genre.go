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

type GenreRequest struct {
	NameGenre string `json:"name_genre"`
}

func GenresHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/") // Si existe una barra al final la elimina.
	parts := strings.Split(path, "/")           // Para dividir el path en partes usando / como separador.

	if len(parts) == 2 && parts[1] == "genres" { // Caso /genres.
		switch r.Method {
		case http.MethodGet:
			getGenres(w, r)
		case http.MethodPost:
			createGenre(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed) // HTTP 405.
		}
		return
	}

	if len(parts) == 3 && parts[1] == "genres" { // Caso /genres/{id}.
		id, err := strconv.Atoi(parts[2]) // Convierte el id de string a int.
		if err != nil {                   // Si no se puede convertir a int, devuelve un error.
			http.Error(w, "ID de género inválido", http.StatusBadRequest) // HTTP 400.
			return
		}
		switch r.Method {
		case http.MethodGet:
			getGenreByID(w, r, int32(id))
		case http.MethodPut:
			updateGenre(w, r, int32(id))
		case http.MethodDelete:
			deleteGenre(w, r, int32(id))
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed) // HTTP 405.
		}
	}
}

func ValidateGenre(g GenreRequest) error {
	if strings.TrimSpace(g.NameGenre) == "" {
		return errors.New("El nombre del género es obligatorio")
	}
	return nil
}

func getGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := Queries.ListGenre(r.Context()) // Llama a ListGenre de sqlc para obtener todos los géneros.
	if err != nil {                               // Si hay un error al obtener los géneros, devuelve un error 500.
		http.Error(w, "Error al obtener los géneros", http.StatusInternalServerError) // HTTP 500.
		return
	}
	w.Header().Set("Content-Type", "application/json") // Establece el tipo de contenido a JSON.
	json.NewEncoder(w).Encode(genres)                  // Codifica los géneros a JSON y los envía en la respuesta.
}

func createGenre(w http.ResponseWriter, r *http.Request) {
	var g GenreRequest
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil { // Decodifica el JSON del cuerpo de la petición.
		http.Error(w, "JSON inválido", http.StatusBadRequest) // HTTP 400.
		return
	}
	errG := ValidateGenre(g) // Valida los datos del género.
	if errG != nil {         // Si hay un error de validación, devuelve un error 400.
		http.Error(w, errG.Error(), http.StatusBadRequest) // HTTP 400.
		return
	}

	newGenre, err := Queries.CreateGenre(r.Context(), g.NameGenre) // Llama a CreateGenre de sqlc.
	if err != nil {                                                // Si hay un error al crear el género, devuelve un error 500.
		http.Error(w, "Error al crear el género", http.StatusInternalServerError) // HTTP 500.
		return
	}

	w.Header().Set("Content-Type", "application/json") // Establece el tipo de contenido a JSON.
	w.WriteHeader(http.StatusCreated)                  // HTTP 201.
	json.NewEncoder(w).Encode(newGenre)                // Codifica el nuevo género a JSON y lo envía en la respuesta.
}

func getGenreByID(w http.ResponseWriter, r *http.Request, id int32) {
	genre, err := Queries.GetGenre(r.Context(), id) // Llama a GetGenre de sqlc para obtener el género por ID.
	if err != nil {                                 // Si hay un error al obtener el género, devuelve un error.
		if err == sql.ErrNoRows { // Si no se encuentra el género, devuelve un error 404.
			http.Error(w, "Género no encontrado", http.StatusNotFound) // HTTP 404.
			return
		}
		http.Error(w, "Error al obtener el género", http.StatusInternalServerError) // HTTP 500.
		return
	}
	w.Header().Set("Content-Type", "application/json") // Establece el tipo de contenido a JSON.
	json.NewEncoder(w).Encode(genre)                   // Codifica el género a JSON y lo envía en la respuesta.
}

func updateGenre(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := Queries.GetGenre(r.Context(), id) // Verifica si el género existe.
	if err != nil {
		if err == sql.ErrNoRows { // Si no se encuentra el género, devuelve un error 404.
			http.Error(w, "Género no encontrado", http.StatusNotFound) // HTTP 404.
			return
		}
		http.Error(w, "Error al obtener el género", http.StatusInternalServerError) // HTTP 500.
		return
	}

	var g GenreRequest
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil { // Decodifica el JSON del cuerpo de la petición.
		http.Error(w, "JSON inválido", http.StatusBadRequest) // HTTP 400.
		return
	}
	errG := ValidateGenre(g)
	if errG != nil { // Si hay un error de validación, devuelve un error 400.
		http.Error(w, errG.Error(), http.StatusBadRequest) // HTTP 400.
		return
	}

	params := db.UpdateGenreParams{
		ID:        id,
		NameGenre: g.NameGenre,
	}

	if err := Queries.UpdateGenre(r.Context(), params); err != nil { // Llama a UpdateGenre de sqlc para actualizar el género.
		http.Error(w, "Error al actualizar el género", http.StatusInternalServerError) // HTTP 500.
		return
	}

	updatedGenre, err := Queries.GetGenre(r.Context(), id) // Obtiene el género actualizado.
	if err != nil {                                        // Si hay un error al obtener el género, devuelve un error 500.
		http.Error(w, "Error al obtener el género actualizado", http.StatusInternalServerError) // HTTP 500.
		return
	}

	w.Header().Set("Content-Type", "application/json") // Establece el tipo de contenido a JSON.
	json.NewEncoder(w).Encode(updatedGenre)            // Codifica el género actualizado a JSON y lo envía en la respuesta.
}

func deleteGenre(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := Queries.GetGenre(r.Context(), id) // Verifica si el género existe.
	if err != nil {
		if err == sql.ErrNoRows { // Si no se encuentra el género, devuelve un error 404.
			http.Error(w, "Género no encontrado", http.StatusNotFound) // HTTP 404.
			return
		}
		http.Error(w, "Error al obtener el género", http.StatusInternalServerError) // HTTP 500.
		return
	}

	if err := Queries.DeleteGenre(r.Context(), id); err != nil { // Llama a DeleteGenre de sqlc para eliminar el género.
		http.Error(w, "Error al eliminar el género", http.StatusInternalServerError) // HTTP 500.
		return
	}
	w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content.
}
