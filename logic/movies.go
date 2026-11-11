package logic

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

// MovieInput son los datos de una película tal como llegan del formulario.
type MovieInput struct {
	Name        string
	Author      string
	Duration    int32
	Score       string
	Review      string
	Description string
}

func ValidateMovie(m MovieInput) error {
	switch {
	case strings.TrimSpace(m.Name) == "":
		return errors.New("el nombre de la película es obligatorio")
	case strings.TrimSpace(m.Author) == "":
		return errors.New("el autor de la película es obligatorio")
	case m.Duration <= 0:
		return errors.New("la duración de la película debe ser mayor a 0")
	case utf8.RuneCountInString(m.Review) > 500:
		return errors.New("la reseña no puede superar los 500 caracteres")
	case utf8.RuneCountInString(m.Description) > 500:
		return errors.New("la descripción no puede superar los 500 caracteres")
	}

	if m.Score != "" {
		s, err := strconv.ParseFloat(m.Score, 64)
		if err != nil || !(s >= 0 && s <= 10) {
			return errors.New("el puntaje debe ser un número entre 0 y 10")
		}
	}

	return nil
}
