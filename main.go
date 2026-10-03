package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	db "proyecto/db/sqlc" // Import de sqlc.
	"proyecto/logic"      // Import de la lógica de negocio

	_ "github.com/lib/pq"
)

func main() {
	// Cadena de conexión
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=peliculas_db sslmode=disable"

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

	http.HandleFunc("/movies/", logic.MoviesHandler)

	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fileServer)

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %v", err)
	}
}
