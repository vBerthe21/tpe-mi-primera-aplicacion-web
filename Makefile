APP_NAME := my-app
DB_URL := postgres://postgres:postgres@localhost:5432/peliculas_db?sslmode=disable

.PHONY: all generate migrate apply status build test clean

all: build

generate:
	@sqlc generate
	# @templ generate

# Genera una migración: make migrate name=nombre_de_la_migracion
migrate:
	@test -n "$(name)" || (echo "Uso: make migrate name=nombre" && exit 1)
	atlas migrate diff "$(name)" --dir "file://db/migrations" --to \
	"file://db/schema/schema.sql" --dev-url "docker://postgres/15/dev?search_path=public"

# Aplica las migraciones pendientes
apply:
	atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"

# Muestra el estado de las migraciones
status:
	atlas migrate status --dir "file://db/migrations" --url "$(DB_URL)"

# Construye el binario de la aplicación
build: generate
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .

# Ejecución de tests automatizando Docker y las migraciones de las filminas
test: generate build
	@echo "-> Limpiando entorno..."
	@docker compose down -v
	@echo "-> Levantando PostgreSQL..."
	@docker compose up -d
	@echo "-> Esperando a la base de datos..."
	@sleep 3
	@echo "-> Aplicando migraciones pendientes..."
	@atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"
	@echo "-> Corriendo pruebas..."
	@go test -v ./...
	@echo "-> Limpiando entorno..."
	@docker compose down -v

clean:
	@rm -rf tmp