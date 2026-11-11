package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	db "proyecto/db/sqlc" // Import de sqlc.
	"proyecto/logic"
	"proyecto/views"

	_ "github.com/lib/pq"

	"strconv"
	"strings"
)

func getEnv(key, fallback string) string { // Función para obtener variables de entorno.
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// nullStr convierte una cadena vacía en NULL para la base.
func nullStr(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func main() {
	// Cadena de conexión obtenida del entorno.
	connStr := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/peliculas_db?sslmode=disable")

	// 1. Abrir conexión con la base de datos
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error al abrir conexión: %v", err)
	}
	defer dbConn.Close()

	// 2. Comprobar que responde
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a PostgreSQL: %v", err)
	}
	fmt.Println("¡Conexión a PostgreSQL exitosa!")

	mux := http.NewServeMux()
	queries := db.New(dbConn)

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		// 1. obtener las películas con queries.ListMovie(...)
		movies, err := queries.ListMovie(r.Context())
		if err != nil {
			log.Printf("error al listar películas: %v", err)
			http.Error(w, "error interno del servidor", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = views.IndexPage("CINE+ | Catálogo de Películas", movies).Render(r.Context(), w)
		if err != nil {
			log.Printf("error al renderizar: %v", err)
		}
	})

	mux.HandleFunc("POST /movies", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "formulario inválido", http.StatusBadRequest)
			return
		}

		duracion, err := strconv.ParseInt(r.FormValue("duration"), 10, 32)
		if err != nil {
			http.Error(w, "la duración debe ser un número entero", http.StatusBadRequest)
			return
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = queries.CreateMovie(r.Context(), db.CreateMovieParams{
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

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	port := getEnv("PORT", "8080")
	if port[0] != ':' {
		port = ":" + port
	}

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
