# Nombre del Proyecto: CINE+
# Integrantes: Berthelot Valentin; Rodriguez Giacomasso Santiago.

# Entrega TP1: Mi Primera Aplicación Web.
* Para ejecutar el proyecto: abrir una terminal en la raíz y ejecutar "go run .", luego, abrir el navegador de preferencia e ingresar a http://localhost:8080

# Entrega TP2: Mi Primera Aplicación Web.
## Para realizar la automatización y persistencia del proyecto, se debe contar con los siguientes requisitos:
* Go en una versión >= 1.25.0.
* PostgreSQL como motor de base de datos.
* Docker y Docker Compose, necesario para levantar el contenedor con la base de datos.
* Make, para la automatización de tareas necesarias para realizar el test.
* Atlas, utilizada para realizar las migraciones de la base de datos.
* SQLC, usado para la generación de código GO a partir de consultas SQL realizadas en el motor PostgreSQL.

## Estructura del proyecto:
```text
├── db/
│   ├── migrations/
│   │   └── atlas.sum
│   ├── queries/
│   │   └── queries.sql
│   ├── schema/
│   │   └── schema.sql
│   └── sqlc/
│       └── db_test.go   # Único archivo rastreado (tests)
├── logic/
│   └── movies.go
├── static/
│   └── index.html
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── Makefile
├── readme.md
└── sqlc.yaml
```

## Composición de entidades.
#### Movies.
La tabla que contiene la información referente a las películas. Contiene id (PK), nombre, autor, duración, y, como opcionales: puntaje, review y descripción.
#### Genre.
Tabla que contiene los géneros de las películas. Contiene id (PK) y nombre.
#### MovieGenre.
Tabla de la relación que surge entre las anteriores. Contiene las claves de las otras tablas: id de la pelicula y id del genero.

## Persistencia:
Las modificaciones al esquema se realizan en db/schema/schema.sql y se migran a la base de datos PostgreSQL usando Atlas.
Las consultas a la base de datos se escriben en SQL en db/queries/queries.sql.
SQLC lee las consultas y genera el código en Go dentro de db/sqlc/

## Ejecución de parte 2:
* Para ejecutar esta segunda parte, debe hacerse un clonado del repositorio, obteniendose de esa forma el link, posteriormente en el IDE de preferencia se
abre una terminal, se ejecuta el comando **git clone** junto con el link del repositorio. Por último, desde la ruta principal del proyecto debe ejecutarse **make test**.