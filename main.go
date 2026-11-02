package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	db "proyecto/db/sqlc" // Import de sqlc.
	"proyecto/logic"      // Import de la lógica de negocio

	_ "github.com/lib/pq"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	// Cadena de conexión obtenida del entorno o fallback local
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

	logic.Queries = db.New(dbConn)

	http.HandleFunc("/genres", logic.GenresHandler)
	http.HandleFunc("/genres/", logic.GenresHandler)

	http.HandleFunc("/movies", logic.MoviesHandler)
	http.HandleFunc("/movies/", logic.MoviesHandler)

	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fileServer)

	port := getEnv("PORT", "8080")
	if port[0] != ':' {
		port = ":" + port
	}

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %v", err)
	}
}
