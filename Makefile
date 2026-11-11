-include .env
export

APP_NAME := my-app
DB_URL ?= postgres://postgres:postgres@localhost:5432/peliculas_db?sslmode=disable

.PHONY: all generate migrate apply status build test up run clean down

all: build

generate:
	@sqlc generate
	@templ generate


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

# Levanta la base de datos en Docker
up:
	@docker compose up -d --wait

# Detiene y limpia la infraestructura
down:
	docker compose down -v

# Levanta PostgreSQL (si no está levantado), espera a que esté lista y ejecuta Go
run: generate up apply
	@echo "-> Iniciando servidor Go..."
	@go run main.go

# Ejecución completa para la evaluación del TP
run-app: generate build
	@echo "-> [1/3] Tareas previas: Limpiando entorno..."
	@docker compose down -v
	@echo "-> Levantando BD y esperando..."
	@$(MAKE) up
	@echo "-> Aplicando migraciones..."
	@$(MAKE) apply
	@echo "-> [2/3] Ejecutando la aplicación (Presiona Ctrl+C para detener)..."
	@bash -c 'trap "echo \"\n-> [3/3] Tareas posteriores: Limpiando contenedores...\" && docker compose down -v" EXIT; ./tmp/$(APP_NAME)'
	
# Ejecución de tests automatizando Docker y las migraciones
test: generate build
	@echo "-> Limpiando entorno..."
	@docker compose down -v
	@echo "-> Levantando PostgreSQL..."
	@$(MAKE) up
	@echo "-> Aplicando migraciones pendientes..."
	@atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"
	@echo "-> Corriendo pruebas..."
	@go test -v ./...
	@echo "-> Limpiando entorno..."
	@docker compose down -v


clean:
	@rm -rf tmp