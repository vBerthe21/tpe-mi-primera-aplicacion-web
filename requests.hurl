# ==========================================
# 1. PRUEBAS PARA LA ENTIDAD GÉNERO (/genres)
# ==========================================

# A) Crear Género 1
POST http://localhost:8080/genres
Content-Type: application/json
{
  "name_genre": "Acción"
}
HTTP 201
[Captures]
genre_id_1: jsonpath "$.id"

# B) Crear Género 2
POST http://localhost:8080/genres
Content-Type: application/json
{
  "name_genre": "Ciencia Ficción"
}
HTTP 201
[Captures]
genre_id_2: jsonpath "$.id"

# C) Listar Géneros
GET http://localhost:8080/genres
HTTP 200

# D) Obtener Género por ID
GET http://localhost:8080/genres/{{genre_id_1}}
HTTP 200

# E) Actualizar Género
PUT http://localhost:8080/genres/{{genre_id_1}}
Content-Type: application/json
{
  "name_genre": "Acción y Aventura"
}
HTTP 200

# ==========================================
# 2. PRUEBAS PARA LA ENTIDAD PELÍCULA (/movies)
# ==========================================

# A) Crear Película
POST http://localhost:8080/movies
Content-Type: application/json
{
  "name": "Inception",
  "author": "Christopher Nolan",
  "duration": 148,
  "score": "9.0",
  "review": "Excelente película de ciencia ficción",
  "description": "Un ladrón que roba secretos corporativos a través del uso de la tecnología de intercambio de sueños."
}
HTTP 201
[Captures]
movie_id: jsonpath "$.id"

# B) Listar Películas
GET http://localhost:8080/movies
HTTP 200

# C) Obtener Película por ID
GET http://localhost:8080/movies/{{movie_id}}
HTTP 200

# D) Actualizar Película
PUT http://localhost:8080/movies/{{movie_id}}
Content-Type: application/json
{
  "name": "Inception (El Origen)",
  "author": "Christopher Nolan",
  "duration": 148,
  "score": "9.5",
  "review": "Obra maestra del cine moderno",
  "description": "Un ladrón que roba secretos corporativos a través del uso de la tecnología de intercambio de sueños."
}
HTTP 200

# ==========================================
# 3. PRUEBAS PARA LA RELACIÓN MOVIE_GENRE
# ==========================================

# A) Asociar Género 1 a Película
POST http://localhost:8080/movies/{{movie_id}}/genres
Content-Type: application/json
{
  "genre_id": {{genre_id_1}}
}
HTTP 201

# B) Asociar Género 2 a Película
POST http://localhost:8080/movies/{{movie_id}}/genres
Content-Type: application/json
{
  "genre_id": {{genre_id_2}}
}
HTTP 201

# C) Desasociar Género 1 de la Película
DELETE http://localhost:8080/movies/{{movie_id}}/genres/{{genre_id_1}}
HTTP 204

# ==========================================
# 4. LIMPIEZA / ELIMINACIÓN DE ENTIDADES
# ==========================================

# A) Eliminar Película
DELETE http://localhost:8080/movies/{{movie_id}}
HTTP 204

# B) Eliminar Géneros creados
DELETE http://localhost:8080/genres/{{genre_id_1}}
HTTP 204

DELETE http://localhost:8080/genres/{{genre_id_2}}
HTTP 204