package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	db "proyecto/db/sqlc"
	"proyecto/handlers"

	_ "github.com/lib/pq"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
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

	// 3. Inyectar dependencias (queries -> handlers)
	queries := db.New(dbConn)
	h := handlers.New(queries)

	// 4. Registrar rutas
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", h.Index)

	// Películas
	mux.HandleFunc("POST /movies", h.CreateMovie)
	mux.HandleFunc("DELETE /movies/{id}", h.DeleteMovie)
	mux.HandleFunc("GET /movies/{id}/edit", h.EditMovieForm)
	mux.HandleFunc("PUT /movies/{id}", h.UpdateMovie)

	// Géneros
	mux.HandleFunc("POST /genres", h.CreateGenre)
	mux.HandleFunc("DELETE /genres/{id}", h.DeleteGenre)

	// 5. Levantar el servidor
	port := getEnv("PORT", "8080")
	if port[0] != ':' {
		port = ":" + port
	}

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
