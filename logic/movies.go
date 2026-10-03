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

var Queries *db.Queries

type MovieRequest struct {
	Name        string  `json:"name"`
	Author      string  `json:"author"`
	Duration    int32   `json:"duration"`
	Score       *string `json:"score"`
	Review      *string `json:"review"`
	Description *string `json:"description"`
}

func ValidateMovie(m MovieRequest) error {
	switch {
	case (m.Name == ""):
		return errors.New("name es obligatorio")
	case (m.Author == ""):
		return errors.New("author es obligatorio")
	case (m.Duration <= 0):
		return errors.New("duration debe ser mayor a 0")
	}
	return nil
}

func MoviesHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/") //Si existe una barra al final la elimina.
	parts := strings.Split(path, "/")           //Para dividir el path en partes usando / como separador.

	if len(parts) == 2 && parts[1] == "movies" { //Caso /movies.
		switch r.Method {
		case http.MethodGet:
			getMovies(w, r)
		case http.MethodPost:
			createMovie(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed) //HTTP 405.
		}
		return
	}

	if len(parts) == 3 && parts[1] == "movies" { //Caso /movies/{id}.
		id, err := strconv.Atoi(parts[2]) //Convierte el id de string a int.
		if err != nil {                   //Si no se puede convertir a int, devuelve un error.
			http.Error(w, "ID de película inválido", http.StatusBadRequest) //HTTP 400.
			return
		}
		switch r.Method {
		case http.MethodGet:
			getMovieByID(w, r, int32(id))
		case http.MethodPut:
			updateMovie(w, r, int32(id))
		case http.MethodDelete:
			deleteMovie(w, r, int32(id))
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed) //HTTP 405.
		}
	}
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := Queries.ListMovie(r.Context()) //Llama a ListMovie de sqlc para obtener todas las películas.
	if err != nil {                               //Si hay un error al obtener las películas, devuelve un error 500.
		http.Error(w, "Error al obtener las películas", http.StatusInternalServerError) //HTTP 500.
		return
	}
	w.Header().Set("Content-Type", "application/json") //Establece el tipo de contenido a JSON.
	json.NewEncoder(w).Encode(movies)                  //Codifica las películas a JSON y las envía en la respuesta.
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var mov MovieRequest
	if err := json.NewDecoder(r.Body).Decode(&mov); err != nil { //Decodifica el JSON del cuerpo de la petición.
		http.Error(w, "JSON inválido", http.StatusBadRequest) //HTTP 400.
		return
	}
	errMov := ValidateMovie(mov) //Valida los datos de la película.
	if errMov != nil {           //Si hay un error de validación, devuelve un error 400.
		http.Error(w, errMov.Error(), http.StatusBadRequest) //HTTP 400.
		return
	}

	params := db.CreateMovieParams{
		Name:     mov.Name,
		Author:   mov.Author,
		Duration: mov.Duration,
	}
	if mov.Score != nil {
		params.Score = sql.NullString{String: *mov.Score, Valid: true} //Si Score no es nulo, lo asigna a params.
	}
	if mov.Review != nil {
		params.Review = sql.NullString{String: *mov.Review, Valid: true} //Si Review no es nulo, lo asigna a params.
	}
	if mov.Description != nil {
		params.Description = sql.NullString{String: *mov.Description, Valid: true} //Si Description no es nulo, lo asigna a params.
	}

	newMovie, err := Queries.CreateMovie(r.Context(), params) //Llama a CreateMovie de sqlc para crear la película.
	if err != nil {                                           //Si hay un error al crear la película, devuelve un error 500.
		http.Error(w, "Error al crear la película", http.StatusInternalServerError) //HTTP 500.
		return
	}

	w.Header().Set("Content-Type", "application/json") //Establece el tipo de contenido a JSON.
	w.WriteHeader(http.StatusCreated)                  //HTTP 201.
	json.NewEncoder(w).Encode(newMovie)                //Codifica la nueva película a JSON y la envía en la respuesta.
}

func getMovieByID(w http.ResponseWriter, r *http.Request, id int32) {
	movie, err := Queries.GetMovie(r.Context(), id) //Llama a GetMovie de sqlc para obtener la película por ID.
	if err != nil {                                 //Si hay un error al obtener la película, devuelve un error 500.
		if err == sql.ErrNoRows { //Si no se encuentra la película, devuelve un error 404.
			http.Error(w, "Película no encontrada", http.StatusNotFound) //HTTP 404.
			return
		}
		http.Error(w, "Error al obtener la película", http.StatusInternalServerError) //HTTP 500.
		return
	}
	w.Header().Set("Content-Type", "application/json") //Establece el tipo de contenido a JSON.
	json.NewEncoder(w).Encode(movie)                   //Codifica la película a JSON y la envía en la respuesta.
}

func updateMovie(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := Queries.GetMovie(r.Context(), id) //Verifica si la película existe.
	if err != nil {
		if err == sql.ErrNoRows { //Si no se encuentra la película, devuelve un error 404.
			http.Error(w, "Película no encontrada", http.StatusNotFound) //HTTP 404.
			return
		}
		http.Error(w, "Error al obtener la película", http.StatusInternalServerError) //HTTP 500.
		return
	}
	var mov MovieRequest
	if err := json.NewDecoder(r.Body).Decode(&mov); err != nil { //Decodifica el JSON del cuerpo de la petición.
		http.Error(w, "JSON inválido", http.StatusBadRequest) //HTTP 400.
		return
	}
	errMov := ValidateMovie(mov)
	if errMov != nil { //Si hay un error de validación, devuelve un error 400.
		http.Error(w, errMov.Error(), http.StatusBadRequest) //HTTP 400.
		return
	}

	params := db.UpdateMovieParams{
		ID:       id,
		Name:     mov.Name,
		Author:   mov.Author,
		Duration: mov.Duration,
	}
	if mov.Score != nil {
		params.Score = sql.NullString{String: *mov.Score, Valid: true} //Si Score no es nulo, lo asigna a params.
	}
	if mov.Review != nil {
		params.Review = sql.NullString{String: *mov.Review, Valid: true} //Si Review no es nulo, lo asigna a params.
	}
	if mov.Description != nil {
		params.Description = sql.NullString{String: *mov.Description, Valid: true} //Si Description no es nulo, lo asigna a params.
	}

	if err := Queries.UpdateMovie(r.Context(), params); err != nil { //Llama a UpdateMovie de sqlc para actualizar la película.
		http.Error(w, "Error al actualizar la película", http.StatusInternalServerError) //HTTP 500.
		return
	}

	updatedMovie, err := Queries.GetMovie(r.Context(), id) //Obtiene la película actualizada.
	if err != nil {                                        //Si hay un error al obtener la película, devuelve un error 500.
		http.Error(w, "Error al obtener la película actualizada", http.StatusInternalServerError) //HTTP 500.
		return
	}

	w.Header().Set("Content-Type", "application/json") //Establece el tipo de contenido a JSON.
	json.NewEncoder(w).Encode(updatedMovie)            //Codifica la película actualizada a JSON y la envía en la respuesta.
}

func deleteMovie(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := Queries.GetMovie(r.Context(), id) //Verifica si la película existe.
	if err != nil {
		if err == sql.ErrNoRows { //Si no se encuentra la película, devuelve un error 404.
			http.Error(w, "Película no encontrada", http.StatusNotFound) //HTTP 404.
			return
		}
		http.Error(w, "Error al obtener la película", http.StatusInternalServerError) //HTTP 500.
		return
	}

	if err := Queries.DeleteMovie(r.Context(), id); err != nil { //Llama a DeleteMovie de sqlc para eliminar la película.
		http.Error(w, "Error al eliminar la película", http.StatusInternalServerError) //HTTP 500.
		return
	}
	w.WriteHeader(http.StatusNoContent) //HTTP 204 No Content.
}
