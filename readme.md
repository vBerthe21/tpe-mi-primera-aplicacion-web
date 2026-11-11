# Nombre del Proyecto: CINE+
# Integrantes: Berthelot Valentin; Rodriguez Giacomasso Santiago.

# Entrega TP5: Mi Primera Aplicación Web.

## Para realizar la automatización y persistencia del proyecto, se debe contar con los siguientes requisitos:
* Go en una versión >= 1.25.0.
* PostgreSQL como motor de base de datos.
* Docker y Docker Compose, necesario para levantar el contenedor con la base de datos.
* Make, para la automatización de tareas necesarias para realizar el test.
* Atlas, utilizada para realizar las migraciones de la base de datos.
* SQLC, usado para la generación de código GO a partir de consultas SQL realizadas en el motor PostgreSQL.
* Templ, utilizado para generar las vistas dinámicas HTML desde el lado del servidor.

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
│       └── db_test.go
├── logic/
│   └── movies.go
├── views/
│   ├── index.templ
│   ├── layout.templ
│   ├── movie_form.templ
│   └── movie_list.templ
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── Makefile
├── readme.md
├── .env.example
├── .gitignore
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


## Ejecución de TP5: SSR con Templ
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