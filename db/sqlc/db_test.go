package db

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const testDSN = "postgres://postgres:postgres@localhost:5432/peliculas_db?sslmode=disable"

func TestQueries_CRUD_Completo(t *testing.T) {
	conn, err := sql.Open("pgx", testDSN)
	if err != nil {
		t.Fatalf("Error al abrir conexión con la BD: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		t.Fatalf("La base de datos no responde al ping: %v", err)
	}

	testQueries := New(conn)
	ctx := context.Background()

	// =============================================================
	// 1. TESTS DE PELÍCULAS
	// =============================================================

	// A. CreateMovie
	movie, err := testQueries.CreateMovie(ctx, CreateMovieParams{
		Name:        "Matrix",
		Author:      "Lana Wachowski",
		Duration:    136,
		Score:       sql.NullString{String: "8.7", Valid: true},
		Review:      sql.NullString{String: "Clásico del cine", Valid: true},
		Description: sql.NullString{String: "Realidad simulada", Valid: true},
	})
	if err != nil {
		t.Fatalf("CreateMovie falló: %v", err)
	}
	if movie.ID == 0 {
		t.Errorf("Se esperaba un ID autogenerado válido")
	}

	// B. GetMovie
	fetchedMovie, err := testQueries.GetMovie(ctx, movie.ID)
	if err != nil {
		t.Fatalf("GetMovie falló: %v", err)
	}
	if fetchedMovie.Name != movie.Name {
		t.Errorf("GetMovie: se esperaba '%s', se obtuvo '%s'", movie.Name, fetchedMovie.Name)
	}

	// C. UpdateMovie
	err = testQueries.UpdateMovie(ctx, UpdateMovieParams{
		ID:          movie.ID,
		Name:        "Matrix Reloaded",
		Author:      "Hermanas Wachowski",
		Duration:    138,
		Score:       sql.NullString{String: "7.2", Valid: true},
		Review:      sql.NullString{String: "Buena secuela", Valid: true},
		Description: sql.NullString{String: "Saga Matrix", Valid: true},
	})
	if err != nil {
		t.Fatalf("UpdateMovie falló: %v", err)
	}

	// Validar Update
	updatedMovie, err := testQueries.GetMovie(ctx, movie.ID)
	if err != nil || updatedMovie.Name != "Matrix Reloaded" {
		t.Errorf("UpdateMovie no actualizó el nombre correctamente")
	}

	// =============================================================
	// 2. TESTS DE GÉNEROS
	// =============================================================

	// A. CreateGenre
	genre, err := testQueries.CreateGenre(ctx, "Ciencia Ficción")
	if err != nil {
		t.Fatalf("CreateGenre falló: %v", err)
	}

	// B. GetGenre (Usamos fetchedGenre para evitar el error de variable no usada)
	fetchedGenre, err := testQueries.GetGenre(ctx, genre.ID)
	if err != nil {
		t.Fatalf("GetGenre falló: %v", err)
	}
	if fetchedGenre.ID != genre.ID {
		t.Errorf("GetGenre: se esperaba ID %d, pero se obtuvo %d", genre.ID, fetchedGenre.ID)
	}

	// C. UpdateGenre
	err = testQueries.UpdateGenre(ctx, UpdateGenreParams{
		ID:        genre.ID,
		NameGenre: "Sci-Fi",
	})
	if err != nil {
		t.Fatalf("UpdateGenre falló: %v", err)
	}

	// =============================================================
	// 3. TESTS DE RELACIÓN N:M (Película - Género)
	// =============================================================

	// A. Asociar Género a Película (Capturamos los 2 valores devueltos: assoc y err)
	assoc, err := testQueries.AddMovieGenre(ctx, AddMovieGenreParams{
		MovieID: movie.ID,
		GenreID: genre.ID,
	})
	if err != nil {
		t.Fatalf("AddMovieGenre falló: %v", err)
	}
	// Usamos 'assoc' para validar la relación
	if assoc.MovieID != movie.ID || assoc.GenreID != genre.ID {
		t.Errorf("AddMovieGenre: la asociación creada no coincide con la película o el género")
	}

	// B. Desasociar Género de Película
	err = testQueries.RemoveMovieGenre(ctx, RemoveMovieGenreParams{
		MovieID: movie.ID,
		GenreID: genre.ID,
	})
	if err != nil {
		t.Fatalf("RemoveMovieGenre falló: %v", err)
	}

	// =============================================================
	// 4. CLEANUP / ELIMINACIÓN (Delete)
	// =============================================================

	// A. DeleteMovie
	err = testQueries.DeleteMovie(ctx, movie.ID)
	if err != nil {
		t.Fatalf("DeleteMovie falló: %v", err)
	}

	// Validar que ya no existe la película
	_, err = testQueries.GetMovie(ctx, movie.ID)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows al consultar película eliminada")
	}

	// B. DeleteGenre
	err = testQueries.DeleteGenre(ctx, genre.ID)
	if err != nil {
		t.Fatalf("DeleteGenre falló: %v", err)
	}
}
