# Nombre del Proyecto: CINE+
# Integrantes: Berthelot Valentin; Rodriguez Giacomasso Santiago.

# Entrega TP6: Mi Primera Aplicación Web.

## Para realizar la automatización y persistencia del proyecto, se debe contar con los siguientes requisitos:
* Go en una versión >= 1.25.0.
* PostgreSQL como motor de base de datos.
* Docker y Docker Compose, necesario para levantar el contenedor con la base de datos.
* Make, para la automatización de tareas necesarias para realizar el test.
* Atlas, utilizada para realizar las migraciones de la base de datos.
* SQLC, usado para la generación de código GO a partir de consultas SQL realizadas en el motor PostgreSQL.
* Conexión a internet para cargar Pico.css y HTMX desde un CDN (se incluyen en `views/layout.templ`).

## Herramientas gestionadas automáticamente por Go:

* Templ, utilizado para generar las vistas dinámicas HTML desde el lado del servidor. Se encuentra declarado en go.mod y se ejecuta mediante go tool templ.

## Estructura del proyecto:
```text
├── db/
│   ├── migrations/
│   │   ├── 20260910145033_init.sql
│   │   └── atlas.sum
│   ├── queries/
│   │   └── queries.sql
│   ├── schema/
│   │   └── schema.sql
│   └── sqlc/
│       └── db_test.go
├── handlers/
│   ├── genres.go
│   ├── handler.go
│   └── movies.go
├── logic/
│   └── movies.go
├── views/
│   ├── genre_form.templ
│   ├── index.templ
│   ├── layout.templ
│   ├── movie_edit.templ
│   ├── movie_form.templ
│   └── movie_list.templ
├── .env.example
├── .gitignore
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── Makefile
├── README.md
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

## Arquitectura del proyecto:

main.go: punto de entrada de la aplicación. Lee las variables de entorno, abre la conexión con PostgreSQL, inyecta *db.Queries en handlers y registra las rutas HTTP.

handlers/: controladores HTTP separados por responsabilidad (handler.go para la carga inicial y helpers comunes, movies.go para el CRUD de películas y genres.go para la gestión de géneros). Devuelven fragmentos HTML para las peticiones HTMX.

logic/: validaciones puras de reglas de negocio (ValidateMovie y ValidateGenre), desacopladas de HTTP y de la base de datos.

views/: componentes .templ renderizados desde el servidor (layout base, formularios, listados, edición inline de películas y fragmentos con hx-swap-oob para sincronizar múltiples secciones del DOM).

db/: esquema, migraciones de Atlas, consultas SQL y código tipado autogenerado por SQLC.

## Funcionalidades con HTMX:

Películas: creación (POST /movies), edición inline (GET /movies/{id}/edit y PUT /movies/{id}) y borrado con confirmación (DELETE /movies/{id}) sin recargar la página.

Géneros: creación (POST /genres) con validación de duplicados y borrado (DELETE /genres/{id}). Mediante Out-of-Band Swaps (hx-swap-oob="true"), al crear o eliminar un género se actualizan simultáneamente la lista de géneros, los checkboxes del formulario de películas y las tarjetas afectadas en el catálogo.

## Ejecución de TP6: Interfaces Reactivas con HTMX
Todo el ciclo de vida del entorno (compilación, generación de código, levantamiento de bases de datos, migraciones y limpieza) está orquestado dentro del archivo `Makefile` para facilitar su evaluación.

Para probar el proyecto completo:

1. Clona este repositorio y abre una terminal en la carpeta raíz.
2. Ejecuta el comando automatizado:
   ```bash
   make run-app
   ```

**¿Qué hace este comando automáticamente?**
* **Tareas Previas:** Genera los archivos `.go` a partir de las plantillas `templ` y `sqlc`, compila el binario, limpia cualquier contenedor residual, levanta PostgreSQL y aplica las migraciones estructurales con Atlas.
* **Ejecución:** Inicia el servidor Go. Podrás acceder a la aplicación abriendo [http://localhost:8080](http://localhost:8080) en tu navegador web.
* **Tareas Posteriores:** Al presionar `Ctrl + C` en la terminal para apagar el servidor, el script interceptará la señal y destruirá de manera segura los contenedores y los volúmenes persistentes asociados, dejando tu sistema completamente limpio.

Nota: Si se desea correr únicamente la batería de pruebas unitarias sobre la base de datos de manera automatizada, se puede ejecutar make test